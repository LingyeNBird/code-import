package cnb

import (
	"ccrctl/pkg/config"
	"ccrctl/pkg/http_client"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/util"
	"ccrctl/pkg/vcs"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"
)

//const (
//	RootOrganization              = "liuhuantest4"
//	ListSubOrganizationEndPoint   = "/" + RootOrganization + "/-/sub-groups?page=1&page_size=1000"
//	CreateSubOrganizationEndPoint = "/groups"
//	CreateRepoEndPoint            = "/" + RootOrganization + "/%s/-/repos"
//	ListRepoEndPoint              = CreateRepoEndPoint + "?page=1&page_size=1000"
//)

const (
	PageSize = "100"
)

var (
	RootOrganization                        = config.Cfg.GetString("cnb.root_organization")
	listRootSubOrganizationEndPoint         = "/" + RootOrganization + "/-/sub-groups?page=%s&page_size=" + PageSize
	listSubOrganizationEndPoint             = "/" + RootOrganization + "/" + "%s" + "/-/sub-groups?page=%s&page_size=" + PageSize
	createSubOrganizationEndPoint           = "/groups"
	createRepoEndPoint                      = "/" + RootOrganization + "/%s/-/repos"
	createRepoUnderRootOrganizationEndPoint = "/" + RootOrganization + "/-/repos"
	listRepoEndPoint                        = createRepoEndPoint + "?page=%s&page_size=" + PageSize
	getRepoInfoEndPoint                     = "/" + RootOrganization + "/%s"
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

type repos struct {
	Name            string    `json:"name"`
	Freeze          time.Time `json:"freeze"`
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
	StarCount          int         `json:"star_count"`
	ForkCount          int         `json:"fork_count"`
	MarkCount          int         `json:"mark_count"`
	LastUpdatedAt      time.Time   `json:"last_updated_at"`
	Language           string      `json:"language"`
	Path               string      `json:"path"`
	ForkedFrom         string      `json:"forked_from"`
	Tags               interface{} `json:"tags"`
	LastUpdateUsername string      `json:"last_update_username"`
	LastUpdateNickname string      `json:"last_update_nickname"`
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

func CreateSubOrganizationIfNotExists(url, token string, depotList []vcs.VCS) (err error) {
	defer logger.Logger.Debugw(util.GetFunctionName(), "url", url, "token", token, "depotList", depotList)
	subGroups, err := GetSubGroupsByRootGroup(url, token)
	if err != nil {
		return err
	}
	for _, depot := range depotList {
		subGroupName := depot.GetSubGroupName()
		parts := strings.Split(subGroupName, "/")
		tmpPath := ""
		if len(parts) > 1 {
			for i := 0; i < len(parts); i++ {
				if i == 0 {
					tmpPath = parts[i]
				} else {
					tmpPath = path.Join(tmpPath, parts[i])
				}
				err := CreateSubOrganization(url, token, tmpPath)
				if err != nil {
					return err
				}
			}
		} else {
			_, exists := subGroups[subGroupName]
			if !exists {
				err := CreateSubOrganization(url, token, subGroupName)
				if err != nil {
					return err
				}
				subGroups[subGroupName] = true
			}
		}

	}
	logger.Logger.Infof("创建子组织完成")
	return nil
}

func CreateSubOrganization(url, token, subGroupName string) (err error) {
	c := http_client.NewClient(url)
	groupPath := path.Join(RootOrganization, subGroupName)
	logger.Logger.Infof("开始创建子组织%s", groupPath)
	body := &CreateOrganization{
		Path: groupPath,
	}
	_, err = c.Request("POST", createSubOrganizationEndPoint, token, body)

	if err != nil {
		if strings.Contains(err.Error(), "组织名称已经被占用") {
			logger.Logger.Infof("子组织%s已存在", groupPath)
			return nil
		} else {
			return fmt.Errorf("创建子组织%s失败: %v", groupPath, err)
		}
	}
	logger.Logger.Infof("创建子组织%s成功", groupPath)
	return nil
}

func CreateRepo(url, token, group, repo string, organizationMappingLevel int, private bool) (err error) {
	c := http_client.NewClient(url)
	var endpoint, visibility string
	switch organizationMappingLevel {
	case 1:
		endpoint = fmt.Sprintf(createRepoEndPoint, group)
	case 2:
		endpoint = createRepoUnderRootOrganizationEndPoint
	}

	if private {
		visibility = "private"
	} else {
		visibility = "public"
	}
	body := &CreateRepoBody{
		Name:       repo,
		Visibility: visibility,
	}
	_, err = c.Request("POST", endpoint, token, body)
	if err != nil {
		return err
	}
	return nil
}

func GetCnbRepoPath(subgroupName, repoName string, organizationMappingLevel int) (repoPath string) {
	switch organizationMappingLevel {
	case 1:
		repoPath = RootOrganization + "/" + subgroupName + "/" + repoName
	case 2:
		repoPath = RootOrganization + "/" + repoName
	}
	return repoPath
}

func HasRepo(url, token, group, repo, organizationMappingLevel string) (has bool, err error) {
	Data, err := GetReposByRepoPath(url, token, group)
	if err != nil {
		return false, err
	}
	_, ok := Data[repo]
	return ok, nil
}

func HasRepoV2(url, token, group, repo string, organizationMappingLevel int) (has bool, err error) {
	c := http_client.NewClient(url)
	var repoPath string
	switch organizationMappingLevel {
	case 1:
		repoPath = group + "/" + repo
	case 2:
		repoPath = repo
	}
	endpoint := fmt.Sprintf(getRepoInfoEndPoint, repoPath)
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

func GetReposByRepoPathFetchPage(url, token, subGroupsName string, page int) (repos []repos, totalRow, pageSize int, err error) {
	c := http_client.NewClient(url)
	endpoint := fmt.Sprintf(listRepoEndPoint, subGroupsName, strconv.Itoa(page))
	resp, header, err := c.RequestV2("GET", endpoint, token, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	err = c.Unmarshal(resp, &repos)
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
	return repos, totalRow, pageSize, nil
}

func GetReposByRepoPath(url, token, group string) (Data map[string]repos, err error) {
	Data = make(map[string]repos)
	page := 1
	for {
		apiRepos, totalRow, pageSize, err := GetReposByRepoPathFetchPage(url, token, group, page)
		if err != nil {
			return nil, err
		}
		for _, v := range apiRepos {
			Data[v.Name] = v
		}
		if page*pageSize >= totalRow {
			break
		}
		page++
	}
	return Data, nil
}

func GetSubGroupsByGroupFetchPage(url, token string, page int) (subGroups []subGroups, totalRow, pageSize int, err error) {
	c := http_client.NewClient(url)
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
	c := http_client.NewClient(url)
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
	c := http_client.NewClient(url)
	endpoint := "/" + RootOrganization
	_, _, respStatusCode, err := c.RequestV3("GET", endpoint, token, nil)
	if err != nil {
		return fmt.Errorf("判断根组织是否存在失败%s", err)
	}
	if respStatusCode == 404 {
		// 创建根组织
		logger.Logger.Infof("根组织不存在:%s", RootOrganization)
		err = CreateRootOrganization(url, token)
		if err != nil {
			return err
		}
		return nil
	}
	if respStatusCode == 200 {
		logger.Logger.Infof("根组织%s已存在", RootOrganization)
		return nil
	}
	return fmt.Errorf("判断根组织是否存在错误的状态码:%d", respStatusCode)
}

func RootOrganizationExists(url, token string) (exists bool, err error) {
	defer logger.Logger.Debugw(util.GetFunctionName(), "url", url, "token", token)
	c := http_client.NewClient(url)
	endpoint := "/" + RootOrganization
	_, _, respStatusCode, err := c.RequestV3("GET", endpoint, token, nil)
	if err != nil {
		return false, fmt.Errorf("判断根组织是否存在失败%s", err)
	}
	if respStatusCode == 404 {
		return false, nil
	}
	if respStatusCode == 200 {
		return true, nil
	}
	return false, fmt.Errorf("判断根组织是否存在错误的状态码:%d", respStatusCode)
}

func CreateRootOrganization(url, token string) (err error) {
	logger.Logger.Infof("开始创建根组织%s", RootOrganization)
	c := http_client.NewClient(url)
	path := RootOrganization
	body := &CreateOrganization{
		Path: path,
	}
	_, err = c.Request("POST", createSubOrganizationEndPoint, token, body)
	if err != nil {
		return fmt.Errorf("创建根组织失败%s", err)
	}
	logger.Logger.Infof("创建根组织%s成功", RootOrganization)
	return nil
}

func GetPushUrl(organizationMappingLevel int, cnbURL, userName, token, projectName, repoName string) string {
	var pushURL string
	parts := strings.Split(cnbURL, "://")
	switch organizationMappingLevel {
	case 1:
		pushURL = parts[0] + "://" + userName + ":" + token + "@" + parts[1] + "/" + RootOrganization + "/" + path.Join(projectName, repoName) + ".git"
	case 2:
		pushURL = parts[0] + "://" + userName + ":" + token + "@" + parts[1] + "/" + RootOrganization + "/" + repoName + ".git"
	}
	return pushURL
}
