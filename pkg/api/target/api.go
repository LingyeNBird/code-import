package target

import (
	"bytes"
	"ccrctl/pkg/config"
	"ccrctl/pkg/http_client"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/util"
	"ccrctl/pkg/vcs"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	PageSize             = "100"
	ReleaseAssetMaxSize  = 1024 * 1024 * 1024 * 50
	DescAssetMaxSize     = 1024 * 1024 * 5
	GroupDescLimitSize   = 200
	GroupRemarkLimitSize = 50
	RepoDescLimitSize    = 350
)

var (
	RootOrganizationName                    = config.Cfg.GetString("cnb.root_organization")
	listRootSubOrganizationEndPoint         = "/" + RootOrganizationName + "/-/sub-groups?page=%s&page_size=" + PageSize
	listSubOrganizationEndPoint             = "/" + RootOrganizationName + "/" + "%s" + "/-/sub-groups?page=%s&page_size=" + PageSize
	createSubOrganizationEndPoint           = "/groups"
	createRepoEndPoint                      = "/" + RootOrganizationName + "/%s/-/Repos"
	createRepoUnderRootOrganizationEndPoint = "/" + RootOrganizationName + "/-/Repos"
	listRepoEndPoint                        = createRepoEndPoint + "?page=%s&page_size=" + PageSize
	getRepoInfoEndPoint                     = "/" + RootOrganizationName + "/%s"
	c                                       = http_client.NewClientV2()
	UploadImgEndPoint                       = "/%s/-/upload/imgs"
	UploadFileEndPoint                      = "/%s/-/upload/files"
	CnbURL                                  = config.Cfg.GetString("cnb.url")
	CnbApiURL                               = config.ConvertToApiURL(CnbURL)
	CnbToken                                = config.Cfg.GetString("cnb.token")
	sourcePlatformURL                       = config.Cfg.GetString("source.url")
)

type users struct {
	Address          string `json:"address"`
	AppreciateStatus int    `json:"appreciate_status"`
	Bio              string `json:"bio"`
	Company          string `json:"company"`
	CreatedAt        string `json:"created_at"`
	Editable         int    `json:"editable"`
	Email            string `json:"email"`
	FollowCount      int    `json:"follow_count"`
	FollowRepoCount  int    `json:"follow_repo_count"`
	FollowerCount    int    `json:"follower_count"`
	Freeze           string `json:"freeze"`
	Gender           int    `json:"gender"`
	GroupCount       int    `json:"group_count"`
	LastLoginAt      string `json:"last_login_at"`
	LastLoginIp      string `json:"last_login_ip"`
	Location         string `json:"location"`
	Mail             string `json:"mail"`
	Nickname         string `json:"nickname"`
	RepoCount        int    `json:"repo_count"`
	RewardAmount     int    `json:"reward_amount"`
	RewardCount      int    `json:"reward_count"`
	Site             string `json:"site"`
	StarsCount       int    `json:"stars_count"`
	Type             int    `json:"type"`
	UpdatedNameAt    string `json:"updated_name_at"`
	UpdatedNickAt    string `json:"updated_nick_at"`
	Username         string `json:"username"`
	Verified         int    `json:"verified"`
	WechatMp         string `json:"wechat_mp"`
	WechatMpQrcode   string `json:"wechat_mp_qrcode"`
}

type subGroups struct {
	Name             string    `json:"name"`
	Remark           string    `json:"remark"`
	Description      string    `json:"description"`
	Site             string    `json:"site"`
	Email            string    `json:"email"`
	Freeze           bool      `json:"freeze"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	FollowCount      int       `json:"follow_count"`
	MemberCount      int       `json:"member_count"`
	AllMemberCount   int       `json:"all_member_count"`
	SubGroupCount    int       `json:"sub_group_count"`
	SubRepoCount     int       `json:"sub_repo_count"`
	AllSubGroupCount int       `json:"all_sub_group_count"`
	AllSubRepoCount  int       `json:"all_sub_repo_count"`
	HasSubGroup      bool      `json:"has_sub_group"`
	Path             string    `json:"path"`
}

type Repos struct {
	Id              int64     `json:"id"`
	Name            string    `json:"name"`
	Freeze          bool      `json:"freeze"`
	Status          int       `json:"status"`
	VisibilityLevel string    `json:"visibility_level"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Description     string    `json:"description"`
	Site            string    `json:"site"`
	Topics          string    `json:"topics"`
	License         string    `json:"license"`
	DisplayModule   struct {
		Activity     bool `json:"activity"`
		Contributors bool `json:"contributors"`
		Release      bool `json:"release"`
	} `json:"display_module"`
	StarCount            int           `json:"star_count"`
	ForkCount            int           `json:"fork_count"`
	MarkCount            int           `json:"mark_count"`
	LastUpdatedAt        interface{}   `json:"last_updated_at"`
	Language             string        `json:"language"`
	WebUrl               string        `json:"web_url"`
	Path                 string        `json:"path"`
	Tags                 interface{}   `json:"tags"`
	OpenIssueCount       int           `json:"open_issue_count"`
	OpenPullRequestCount int           `json:"open_pull_request_count"`
	Languages            []interface{} `json:"languages"`
	LastUpdateUsername   string        `json:"last_update_username"`
	LastUpdateNickname   string        `json:"last_update_nickname"`
}

type errResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (r *Repos) GetSubGroupName() string {
	parts := strings.Split(r.Path, "/")
	return parts[1]
}

func (r *Repos) GetRepoPath() string {
	return r.Path
}

type CreateOrganization struct {
	Description string `json:"description"`
	Path        string `json:"path"`
	Remark      string `json:"remark"`
}

type CreateRepoBody struct {
	Description string `json:"description"`
	License     string `json:"license"`
	Name        string `json:"name"`
	Visibility  string `json:"visibility"`
}

// CreateSubOrganizationIfNotExists 创建子组织，如果不存在则创建（简化优化版本）
func CreateSubOrganizationIfNotExists(url, token string, depotList []vcs.VCS) (err error) {
	defer logger.Logger.Debugw(util.GetFunctionName(), "url", url, "token", token, "depotList", depotList)

	// 1. 收集所有需要创建的子组织路径并去重
	uniqueSubGroups := collectUniqueSubGroups(depotList)
	logger.Logger.Infof("预处理完成，发现 %d 个唯一子组织", len(uniqueSubGroups))

	// 2. 一次性获取现有子组织列表
	existingSubGroups, err := GetSubGroupsByRootGroup(url, token)
	if err != nil {
		return fmt.Errorf("获取子组织列表失败: %v", err)
	}

	// 3. 过滤出需要创建的子组织
	toCreate := filterSubGroupsToCreate(uniqueSubGroups, existingSubGroups)
	if len(toCreate) == 0 {
		logger.Logger.Infof("所有子组织都已存在，无需创建")
		return nil
	}

	logger.Logger.Infof("需要创建 %d 个新的子组织", len(toCreate))

	// 4. 按层级深度顺序创建子组织
	return createSubGroupsSequentially(url, token, toCreate)
}

// collectUniqueSubGroups 收集所有需要创建的子组织路径并去重
func collectUniqueSubGroups(depotList []vcs.VCS) map[string]*vcs.SubGroup {
	uniqueSubGroups := make(map[string]*vcs.SubGroup)

	for _, depot := range depotList {
		subGroup := depot.GetSubGroup()
		subGroupName := subGroup.Name

		// 如果子组织名称为空，跳过
		if subGroupName == "" {
			continue
		}

		// 处理多层级路径，确保所有父级路径都被包含
		parts := strings.Split(subGroupName, "/")
		tmpPath := ""

		for i, part := range parts {
			if tmpPath == "" {
				tmpPath = part
			} else {
				tmpPath = path.Join(tmpPath, part)
			}

			// 保留原有组织信息，不过度加工
			if _, exists := uniqueSubGroups[tmpPath]; !exists {
				if i == len(parts)-1 {
					// 最深层使用原始SubGroup信息
					uniqueSubGroups[tmpPath] = subGroup
				} else {
					// 父级使用原始信息但调整Name字段
					parentSubGroup := &vcs.SubGroup{
						Name:   tmpPath,
						Desc:   subGroup.Desc,   // 保留原有描述
						Remark: subGroup.Remark, // 保留原有备注
					}
					uniqueSubGroups[tmpPath] = parentSubGroup
				}
			}
		}
	}

	return uniqueSubGroups
}

// filterSubGroupsToCreate 过滤出需要创建的子组织
func filterSubGroupsToCreate(uniqueSubGroups map[string]*vcs.SubGroup, existingSubGroups map[string]bool) map[string]*vcs.SubGroup {
	toCreate := make(map[string]*vcs.SubGroup)

	for subGroupPath, subGroup := range uniqueSubGroups {
		if !existingSubGroups[subGroupPath] {
			toCreate[subGroupPath] = subGroup
		}
	}

	return toCreate
}

// createSubGroupsSequentially 按层级深度顺序创建子组织
func createSubGroupsSequentially(url, token string, toCreate map[string]*vcs.SubGroup) error {
	// 按路径深度排序
	paths := make([]string, 0, len(toCreate))
	for path := range toCreate {
		paths = append(paths, path)
	}

	// 简单的深度排序：按斜杠数量排序
	sort.Slice(paths, func(i, j int) bool {
		depthI := strings.Count(paths[i], "/")
		depthJ := strings.Count(paths[j], "/")
		if depthI != depthJ {
			return depthI < depthJ
		}
		return paths[i] < paths[j]
	})

	// 顺序创建每个子组织
	createdCount := 0
	for _, subGroupPath := range paths {
		subGroup := toCreate[subGroupPath]

		err := CreateSubOrganization(url, token, subGroupPath, *subGroup)
		if err != nil {
			// 如果是组织已存在的错误，继续处理下一个
			if strings.Contains(err.Error(), "已存在") {
				logger.Logger.Infof("子组织 %s 已存在，跳过", subGroupPath)
				continue
			}
			return fmt.Errorf("创建子组织 %s 失败: %v", subGroupPath, err)
		}

		createdCount++
		logger.Logger.Infof("成功创建子组织 %s (%d/%d)", subGroupPath, createdCount, len(paths))
	}

	logger.Logger.Infof("子组织创建完成，共创建 %d 个", createdCount)
	return nil
}

func CreateSubOrganization(url, token, subGroupName string, subGroup vcs.SubGroup) (err error) {
	subGroupName = normalizeGroupName(subGroupName)
	groupPath := path.Join(RootOrganizationName, subGroupName)
	logger.Logger.Infof("开始创建子组织%s", groupPath)
	if len(subGroup.Desc) > GroupDescLimitSize {
		subGroup.Desc = subGroup.Desc[:GroupDescLimitSize]
	}
	if len(subGroup.Remark) > GroupRemarkLimitSize {
		subGroup.Remark = subGroup.Remark[:GroupRemarkLimitSize]
	}
	body := &CreateOrganization{
		Path:        groupPath,
		Remark:      subGroup.Remark,
		Description: subGroup.Desc,
	}
	resp, _, statusCode, err := c.RequestV3("POST", createSubOrganizationEndPoint, token, body)

	if err != nil {
		return fmt.Errorf("创建子组织%s失败: %v", groupPath, err)
	}

	if statusCode == 409 {
		var apiErr *errResp
		if unmarshalErr := json.Unmarshal(resp, &apiErr); unmarshalErr != nil {
			return fmt.Errorf("%s: 解析错误响应失败: %v, 原始响应: %s", groupPath, unmarshalErr, string(resp))
		}
		// 10009 仓库占用了组织名称, 10010 组织已存在
		if apiErr.ErrCode == 10009 {
			return fmt.Errorf("%s仓库与要创建的子组织冲突，请先重命名或删除该仓库后再次运行迁移任务", groupPath)
		}
		logger.Logger.Infof("子组织%s已存在", groupPath)
		return nil
	}
	if statusCode == 201 {
		logger.Logger.Infof("创建子组织%s成功", groupPath)
		return nil
	}
	return fmt.Errorf("创建子组织%s失败: %s", groupPath, string(resp))
}

func CreateRepo(url, token, group, repoName, repoDesc string, private bool) (err error) {
	var visibility string
	endpoint := group + "/-/repos"
	if private {
		visibility = "private"
	} else {
		visibility = "public"
	}
	// 修剪描述的前后空格
	repoDesc = strings.TrimSpace(repoDesc)
	if len(repoDesc) > RepoDescLimitSize {
		repoDesc = repoDesc[:RepoDescLimitSize]
	}
	body := &CreateRepoBody{
		Name:        repoName,
		Visibility:  visibility,
		Description: repoDesc,
	}
	_, err = c.Request("POST", endpoint, token, body)
	if err != nil {
		return err
	}
	return nil
}

func GetCnbRepoPathAndGroup(subgroupName, repoName string, organizationMappingLevel int) (repoPath, repoGroup string) {
	switch organizationMappingLevel {
	case 1:
		// 处理当 subgroupName 为空时的情况
		if subgroupName == "" {
			repoPath = "/" + RootOrganizationName + "/" + repoName
			repoGroup = "/" + RootOrganizationName
		} else {
			repoPath = "/" + RootOrganizationName + "/" + subgroupName + "/" + repoName
			repoGroup = "/" + RootOrganizationName + "/" + subgroupName
		}
	case 2:
		repoPath = "/" + RootOrganizationName + "/" + repoName
		repoGroup = "/" + RootOrganizationName
	}
	return repoPath, repoGroup
}

//func HasRepo(url, token, group, repo, organizationMappingLevel string) (has bool, err error) {
//	Data, err := GetReposByGroup(group)
//	if err != nil {
//		return false, err
//	}
//	_, ok := Data[repo]
//	return ok, nil
//}

func HasRepoV2(url, token, repoPath string) (has bool, err error) {
	endpoint := repoPath
	_, _, respStatusCode, err := c.RequestV3("GET", endpoint, token, nil)
	if err != nil {
		return false, fmt.Errorf("判断仓库是否存在失败: %v", err)
	}
	if respStatusCode == 200 {
		return true, nil
	}
	if respStatusCode == 404 {
		return false, nil
	}
	return false, fmt.Errorf("判断仓库是否存在失败: 未知的状态码: %d", respStatusCode)
}

func GetSubGroupsByGroupFetchPage(url, token string, page int) (subGroups []subGroups, totalRow, pageSize int, err error) {
	endpoint := fmt.Sprintf(listRootSubOrganizationEndPoint, strconv.Itoa(page))
	resp, header, err := c.RequestV2("GET", endpoint, token, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	err = c.Unmarshal(resp, &subGroups)
	if err != nil {
		return nil, 0, 0, err
	}
	totalRow, err = strconv.Atoi(header.Get("x-cnb-total"))
	if err != nil {
		return nil, 0, 0, err
	}
	pageSize, err = strconv.Atoi(header.Get("x-cnb-page-size"))
	if err != nil {
		return nil, 0, 0, err
	}
	return subGroups, totalRow, pageSize, nil
}

func GetSubGroupsFetchPage(url, token, subGroupPath string, page int) (subGroups []subGroups, totalRow, pageSize int, err error) {
	endpoint := fmt.Sprintf(listSubOrganizationEndPoint, subGroupPath, strconv.Itoa(page))
	resp, header, err := c.RequestV2("GET", endpoint, token, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	err = c.Unmarshal(resp, &subGroups)
	if err != nil {
		return nil, 0, 0, err
	}
	totalRow, err = strconv.Atoi(header.Get("x-cnb-total"))
	if err != nil {
		return nil, 0, 0, err
	}
	pageSize, err = strconv.Atoi(header.Get("x-cnb-page-size"))
	if err != nil {
		return nil, 0, 0, err
	}
	return subGroups, totalRow, pageSize, nil
}

func GetSubGroups(url, token, subGroupPath string) (Data map[string]bool, err error) {
	Data = make(map[string]bool)
	page := 1
	for {
		apiSubGroups, totalRow, pageSize, err := GetSubGroupsFetchPage(url, token, subGroupPath, page)
		if err != nil {
			return nil, err
		}
		for _, v := range apiSubGroups {
			Data[v.Name] = true
		}
		if page*pageSize >= totalRow {
			break
		}
		page++
	}
	return Data, nil
}

func GetSubGroupsByRootGroup(url, token string) (Data map[string]bool, err error) {
	Data = make(map[string]bool)
	page := 1
	for {
		apiSubGroups, totalRow, pageSize, err := GetSubGroupsByGroupFetchPage(url, token, page)
		if err != nil {
			return nil, err
		}
		for _, v := range apiSubGroups {
			Data[v.Name] = true
		}
		if page*pageSize >= totalRow {
			break
		}
		page++
	}
	return Data, nil
}

func CreateRootOrganizationIfNotExists(url, token string) (err error) {
	defer logger.Logger.Debugw(util.GetFunctionName(), "url", url, "token", token)
	endpoint := "/" + RootOrganizationName
	_, _, respStatusCode, err := c.RequestV3("GET", endpoint, token, nil)
	if err != nil {
		return fmt.Errorf("判断根组织是否存在失败%s", err)
	}
	if respStatusCode == 404 {
		// 创建根组织
		logger.Logger.Infof("根组织不存在:%s", RootOrganizationName)
		err = CreateRootOrganization(url, token)
		if err != nil {
			return err
		}
		return nil
	}
	if respStatusCode == 200 {
		logger.Logger.Infof("根组织%s已存在", RootOrganizationName)
		return nil
	}
	return fmt.Errorf("判断根组织是否存在错误的状态码:%d", respStatusCode)
}

func RootOrganizationExists(url, token string) (exists bool, err error) {
	defer logger.Logger.Debugw(util.GetFunctionName(), "url", url, "token", token)
	endpoint := "/" + RootOrganizationName
	body, _, respStatusCode, err := c.RequestV3("GET", endpoint, token, nil)
	if err != nil {
		return false, err
	}
	if respStatusCode == 404 {
		return false, nil
	}
	if respStatusCode == 200 {
		return true, nil
	}
	return false, fmt.Errorf("判断根组织是否存在错误的状态码:%d, 错误详情:%s", respStatusCode, string(body))
}

func CreateRootOrganization(url, token string) (err error) {
	logger.Logger.Infof("开始创建根组织%s", RootOrganizationName)
	path := RootOrganizationName
	body := &CreateOrganization{
		Path: path,
	}
	_, err = c.Request("POST", createSubOrganizationEndPoint, token, body)
	if err != nil {
		return fmt.Errorf("创建根组织失败%s", err)
	}
	logger.Logger.Infof("创建根组织%s成功", RootOrganizationName)
	return nil
}

func GetPushUrl(organizationMappingLevel int, cnbURL, userName, token, projectName, repoName string) string {
	var pushURL string
	u, _ := url.Parse(cnbURL)
	switch organizationMappingLevel {
	case 1:
		pushURL = fmt.Sprintf("%s://%s", u.Scheme, path.Join(fmt.Sprintf("%s:%s@%s", userName, token, u.Host), RootOrganizationName, projectName, repoName))
	case 2:
		pushURL = fmt.Sprintf("%s://%s", u.Scheme, path.Join(fmt.Sprintf("%s:%s@%s", userName, token, u.Host), RootOrganizationName, repoName))
	}
	return pushURL
}

type CreateReleaseReq struct {
	Body            string `json:"body"`
	Draft           bool   `json:"draft"`
	MakeLatest      string `json:"make_latest"`
	Name            string `json:"name"`
	Prerelease      bool   `json:"prerelease"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
}

type CreateReleaseRes struct {
	Id           string      `json:"id"`
	TagName      string      `json:"tag_name"`
	TagCommitish string      `json:"tag_commitish"`
	Name         string      `json:"name"`
	Body         string      `json:"body"`
	Draft        bool        `json:"draft"`
	Prerelease   bool        `json:"prerelease"`
	IsLatest     bool        `json:"is_latest"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	PublishedAt  time.Time   `json:"published_at"`
	Author       interface{} `json:"author"`
	Assets       interface{} `json:"assets"`
}

// CreateRelease 在指定仓库中创建一个新的发布版本
// 参数:
//   - repoPath: 仓库路径
//   - name: 发布版本名称
//   - desc: 发布版本描述
//   - tagName: 标签名称
//   - projectID: 项目ID
//   - preRelease: 是否为预发布版本
//   - vcs: VCS接口实现
//
// 返回值:
//   - releaseID: 发布版本的ID
//   - exist: 是否已存在相同发布版本
//   - err: 错误信息
func CreateRelease(repoPath, projectID string, release vcs.Releases, vcs vcs.VCS) (releaseID string, exist bool, err error) {
	// 记录函数调用信息
	defer logger.Logger.Debugw(util.GetFunctionName(),
		"repoPath", repoPath,
		"name", release.Name,
		"body", release.Body,
		"tagName", release.TagName)

	// 构建API端点
	endpoint := fmt.Sprintf("/%s/%s/-/releases", RootOrganizationName, repoPath)

	// 获取发布版本中的附件信息
	attachments, err := vcs.GetReleaseAttachments(release.Body, repoPath, projectID)
	if err != nil {
		logger.Logger.Errorf("获取发布版本附件失败 [%s:%s]: %v", repoPath, release.Name, err)
		return "", false, fmt.Errorf("获取发布版本附件失败: %w", err)
	}

	// 处理附件并更新描述内容
	newDesc := release.Body
	for _, attachment := range attachments {
		// 上传附件到CNB平台
		path, err := UploadReleaseDescImgAndAttachments(attachment)
		if err != nil {
			logger.Logger.Errorf("上传发布版本附件失败 [%s:%s]: %v", repoPath, release.Name, err)
			return "", false, fmt.Errorf("上传发布版本附件失败: %w", err)
		}
		// 替换描述中的附件链接
		newDesc = strings.ReplaceAll(newDesc, attachment.Url, path)
	}

	// 构建创建发布版本的请求体
	body := &CreateReleaseReq{
		Body:       newDesc,
		Name:       release.Name,
		TagName:    release.TagName,
		Prerelease: release.Prerelease,
		Draft:      release.Draft,
	}

	// 发送创建发布版本的请求
	res, _, statusCode, err := c.RequestV4(http.MethodPost, endpoint, body)
	if err != nil {
		if statusCode == http.StatusConflict {
			logger.Logger.Warnf("发布版本已存在 [%s:%s]", repoPath, release.Name)
			return "", true, nil
		}
		logger.Logger.Errorf("创建发布版本失败 [%s:%s]: %v", repoPath, release.Name, err)
		return "", false, fmt.Errorf("创建发布版本失败: %v", err)
	}

	// 解析响应数据
	var data CreateReleaseRes
	if err := json.Unmarshal(res, &data); err != nil {
		logger.Logger.Errorf("解析发布版本响应失败 [%s:%s]: %v", repoPath, release.Name, err)
		return "", false, fmt.Errorf("解析发布版本响应失败: %w", err)
	}

	return data.Id, false, nil
}

type UploadImgOrFileRes struct {
	Assets struct {
		Path        string `json:"path"`
		ContentType string `json:"content_type"`
		Name        string `json:"name"`
		Size        int    `json:"size"`
	} `json:"assets"`
	UploadUrl string        `json:"upload_url"`
	Form      UploadImgForm `json:"form"`
	Token     string        `json:"token"`
}

type UploadImgForm struct {
	ContentType    string `json:"Content-Type"`
	Bucket         string `json:"bucket"`
	Key            string `json:"key"`
	Policy         string `json:"policy"`
	XAmzAlgorithm  string `json:"x-amz-algorithm"`
	XAmzCredential string `json:"x-amz-credential"`
	XAmzDate       string `json:"x-amz-date"`
	XAmzSignature  string `json:"x-amz-signature"`
}

type UploadImgsReq struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

func UploadReleaseDescImgAndAttachments(attachment vcs.Attachment) (path string, err error) {
	res, err := GetCosUploadUrlAndForm(attachment)
	if err != nil {
		logger.Logger.Errorf("Get cos  upload form error: %v", err)
		return "", err
	}
	err = UploadFileToCos(res.UploadUrl, res.Form, attachment)
	if err != nil {
		logger.Logger.Errorf("Upload file to cos error: %v", err)
		return "", err
	}
	return res.Assets.Path, nil
}

func UploadReleaseAsset(repoPath, releaseID, assetName string, data []byte) (err error) {
	if len(data) > ReleaseAssetMaxSize {
		logger.Logger.Warnf("%s附件大小超过5GiB，跳过上传", assetName)
		return nil
	}
	uploadURL, err := GetReleaseAssetUploadUrl(repoPath, releaseID, assetName, len(data))
	if err != nil {
		logger.Logger.Errorf("Get upload url error: %v", err)
		return err
	}
	err = c.UploadData(uploadURL.UploadUrl, data)
	if err != nil {
		logger.Logger.Errorf("Upload data error: %v", err)
		return err
	}
	err = ConfirmUpload(uploadURL.VerifyUrl)
	if err != nil {
		return err
	}
	return nil
}

func ConfirmUpload(verifyUrl string) (err error) {
	_, _, _, err = c.RequestWithURL(http.MethodPost, verifyUrl, nil)
	if err != nil {
		logger.Logger.Errorf("Confirm  upload error: %v", err)
		return err
	}
	return nil
}

type UploadUrl struct {
	UploadUrl    string `json:"upload_url"`
	ExpiresInSec int    `json:"expires_in_sec"`
	VerifyUrl    string `json:"verify_url"`
}

type GetReleaseUploadUrlReq struct {
	AssetName string `json:"asset_name"`
	Size      int    `json:"size"`
}

func GetReleaseAssetUploadUrl(repoPath, releaseID, assetName string, size int) (uploadURL UploadUrl, err error) {
	reqPath := fmt.Sprintf("/%s/%s/-/releases/%s/asset-upload-url", RootOrganizationName, repoPath, releaseID)
	body := &GetReleaseUploadUrlReq{
		AssetName: assetName,
		Size:      size,
	}
	res, _, _, err := c.RequestV4(http.MethodPost, reqPath, body)
	if err != nil {
		logger.Logger.Errorf("Get upload url error: %v", err)
		return uploadURL, err
	}
	err = json.Unmarshal(res, &uploadURL)
	if err != nil {
		logger.Logger.Errorf("Unmarshal upload url error: %v", err)
		return uploadURL, err
	}
	return uploadURL, nil
}

func UploadFileToCos(reqUrl string, form UploadImgForm, attachment vcs.Attachment) (err error) {
	// 创建一个缓冲区以写入我们的表单数据
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// 添加其他字段
	w.WriteField("Content-Type", form.ContentType)
	w.WriteField("bucket", form.Bucket)
	w.WriteField("key", form.Key)
	w.WriteField("policy", form.Policy)
	w.WriteField("x-amz-algorithm", form.XAmzAlgorithm)
	w.WriteField("x-amz-credential", form.XAmzCredential)
	w.WriteField("x-amz-date", form.XAmzDate)
	w.WriteField("x-amz-signature", form.XAmzSignature)

	//添加 file 字段
	//if err := http_client.AddFormFile(w, "file", "c436cc0b-7aeb-4029-9ff1-4fa4cdf1f3d1.png", data); err != nil {
	//	return err
	//}

	writer, err := w.CreateFormFile("file", attachment.Name)

	io.Copy(writer, bytes.NewReader(attachment.Data))

	if err != nil {
		fmt.Println(err)
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	_, err = c.SendUploadRequest(reqUrl, w.FormDataContentType(), &b)
	if err != nil {
		return err
	}
	return nil
}

func PlatformConfirmUpload(repoPath, token string) (err error) {
	reqPath := fmt.Sprintf("%s?token=%s", repoPath, token)
	_, _, _, err = c.RequestV4(http.MethodPut, reqPath, nil)
	if err != nil {
		logger.Logger.Errorf("Confirm  upload error: %v", err)
		return err
	}
	return nil
}

func GetCosUploadUrlAndForm(attachment vcs.Attachment) (form UploadImgOrFileRes, err error) {
	var reqPath string
	repoPath := RootOrganizationName + "/" + attachment.RepoPath
	if attachment.Type == "img" {
		reqPath = fmt.Sprintf(UploadImgEndPoint, repoPath)
	} else {
		reqPath = fmt.Sprintf(UploadFileEndPoint, repoPath)
	}
	// 如果上传图片，文件名必须是图片格式后缀
	body := &UploadImgsReq{
		Name: attachment.Name,
		Size: attachment.Size,
	}
	res, _, _, err := c.RequestV4(http.MethodPost, reqPath, body)
	if err != nil {
		logger.Logger.Errorf("Get upload url and form error: %v", err)
		return form, err
	}
	err = json.Unmarshal(res, &form)
	if err != nil {
		logger.Logger.Errorf("Unmarshal upload url error: %v", err)
		return form, err
	}
	return form, nil
}

// normalizeGroupName 规范化子组织名称，使其符合命名规则
// 规则：只能以字母或数字开头和结尾，长度1-50个字符
// 中间可以包含点(.)、下划线(_)和连字符(-)，后缀不能以.git结尾
func normalizeGroupName(input string) string {
	if input == "" {
		return input
	}

	// 去除开头和结尾的非字母数字字符
	runes := []rune(input)

	// 去除开头非字母数字字符
	start := 0
	for start < len(runes) && !isAlphanumeric(runes[start]) {
		start++
	}

	// 去除结尾非字母数字字符
	end := len(runes) - 1
	for end >= start && !isAlphanumeric(runes[end]) {
		end--
	}

	if start > end {
		// 如果没有字母数字字符，返回空字符串
		return ""
	}

	// 截取有效部分
	result := string(runes[start : end+1])

	// 检查长度，超过50则截取
	if len(result) > 50 {
		result = result[:50]
	}

	// 确保不以.git结尾和.svn 结尾
	result = strings.TrimSuffix(result, ".git")
	result = strings.TrimSuffix(result, ".svn")

	return result
}

// isAlphanumeric 检查字符是否为字母或数字
func isAlphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}
