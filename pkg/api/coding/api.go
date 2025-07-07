package coding

import (
	"ccrctl/pkg/config"
	"ccrctl/pkg/http_client"
	"ccrctl/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

const (
	endpoint    = "/open-api"
	PageSize    = 100
	SvnVcsType  = "svn"
	UserName    = "coding"
	ProjectType = "project"
	RepoType    = "repo"
	Team        = "team"
)

var (
	SourceToken = config.Cfg.GetString("source.token")
	SourceURL   = config.Cfg.GetString("source.url")
	Projects    = config.Cfg.GetStringSlice("source.project")
	Repos       = config.Cfg.GetStringSlice("source.repo")
	c           = http_client.NewClient(SourceURL)
	c2          = http_client.NewCodingClient()
)

type UserInfo struct {
	Response struct {
		RequestId string `json:"RequestId"`
		User      struct {
			Id              int    `json:"Id"`
			Status          int    `json:"Status"`
			Email           string `json:"Email"`
			GlobalKey       string `json:"GlobalKey"`
			Avatar          string `json:"Avatar"`
			Name            string `json:"Name"`
			NamePinYin      string `json:"NamePinYin"`
			Phone           string `json:"Phone"`
			PhoneValidation int    `json:"PhoneValidation"`
			EmailValidation int    `json:"EmailValidation"`
			PhoneRegionCode string `json:"PhoneRegionCode"`
			TeamId          int    `json:"TeamId"`
		} `json:"User"`
	} `json:"Response"`
}

type DescribeProjectByNameResponse struct {
	Response struct {
		Project   Project `json:"Project"`
		RequestId string  `json:"RequestId"`
	} `json:"Response"`
}

type Project struct {
	Name        string `json:"Name"`
	Id          int    `json:"Id"`
	Type        int    `json:"Type"`
	DisplayName string `json:"DisplayName"`
	Icon        string `json:"Icon"`
	Description string `json:"Description"`
	CreatedAt   int64  `json:"CreatedAt"`
	MaxMember   int    `json:"MaxMember"`
	TeamId      int    `json:"TeamId"`
	UserOwnerId int    `json:"UserOwnerId"`
	IsDemo      bool   `json:"IsDemo"`
	Archived    bool   `json:"Archived"`
	StartDate   int    `json:"StartDate"`
	UpdatedAt   int64  `json:"UpdatedAt"`
	TeamOwnerId int    `json:"TeamOwnerId"`
	EndDate     int    `json:"EndDate"`
	Status      int    `json:"Status"`
}

type RepoList struct {
	ResponseData struct {
		RequestId string `json:"RequestId"`
		DepotData struct {
			Depots []struct {
				Name     string `json:"Name"`
				HttpsUrl string `json:"HttpsUrl"`
			} `json:"Depots"`
		} `json:"DepotData"`
	} `json:"Response"`
}

type ProjectList struct {
	Data struct {
		PageNumber  int `json:"PageNumber"`
		PageSize    int `json:"PageSize"`
		TotalCount  int `json:"TotalCount"`
		ProjectList []struct {
			Id          int    `json:"Id"`
			CreatedAt   int64  `json:"CreatedAt"`
			UpdatedAt   int64  `json:"UpdatedAt"`
			Status      int    `json:"Status"`
			Type        int    `json:"Type"`
			MaxMember   int    `json:"MaxMember"`
			Name        string `json:"Name"`
			DisplayName string `json:"DisplayName"`
			Description string `json:"Description"`
			Icon        string `json:"Icon"`
			TeamOwnerId int    `json:"TeamOwnerId"`
			UserOwnerId int    `json:"UserOwnerId"`
			StartDate   int    `json:"StartDate"`
			EndDate     int    `json:"EndDate"`
			TeamId      int    `json:"TeamId"`
			IsDemo      bool   `json:"IsDemo"`
			Archived    bool   `json:"Archived"`
		} `json:"ProjectList"`
	} `json:"Data"`
}

type DescribeCodingCurrentUser struct {
	Action string `json:"Action"`
}

type DescribeProjectDepotInfoListRequest struct {
	Action     string
	ProjectId  int64
	PageNumber int
	PageSize   int
}

type DescribeTeamDepotInfoListRequest struct {
	Action     string
	PageNumber int
	PageSize   int
}

type DescribeProjectByNameRequest struct {
	Action      string
	ProjectName string
}

type ResponseError struct {
	Response struct {
		Error struct {
			Message string `json:"Message"`
			Code    string `json:"Code"`
		} `json:"Error"`
		RequestId string `json:"RequestId"`
	} `json:"Response"`
}

type Page struct {
	PageNumber int `json:"PageNumber"`
	PageSize   int `json:"PageSize"`
	TotalPage  int `json:"TotalPage"`
	TotalRow   int `json:"TotalRow"`
}

type Depots struct {
	Id            int    `json:"Id"`
	Name          string `json:"Name"`
	HttpsUrl      string `json:"HttpsUrl"`
	ProjectId     int    `json:"ProjectId"`
	SshUrl        string `json:"SshUrl"`
	WebUrl        string `json:"WebUrl"`
	VcsType       string `json:"VcsType"`
	ProjectName   string `json:"ProjectName"`
	Description   string `json:"Description"`
	CreatedAt     int64  `json:"CreatedAt"`
	LastPushAt    int    `json:"LastPushAt"`
	GroupId       int    `json:"GroupId"`
	GroupName     string `json:"GroupName"`
	GroupType     string `json:"GroupType"`
	DefaultBranch string `json:"DefaultBranch"`
	RepoType      string `json:"RepoType"`
	IsShared      bool   `json:"IsShared"`
}

type DepotData struct {
	Depots []Depots `json:"Depots"`
	Page   Page     `json:"Page"`
}

type DepotInfo struct {
	Response struct {
		DepotData struct {
			Depots []Depots `json:"Depots"`
			Page   struct {
				PageNumber int `json:"PageNumber"`
				PageSize   int `json:"PageSize"`
				TotalPage  int `json:"TotalPage"`
				TotalRow   int `json:"TotalRow"`
			} `json:"Page"`
		} `json:"DepotData"`
		RequestId string `json:"RequestId"`
	} `json:"Response"`
}

type DescribeGitDepotRequest struct {
	Action    string `json:"Action"`
	DepotPath string `json:"DepotPath"`
}

type DescribeGitDepotResponse struct {
	Response struct {
		Depot     Depots `json:"Depot"`
		RequestId string `json:"RequestId"`
	} `json:"Response"`
}

func GetDepotByRepoPath(url, token, repoPath string) (depot Depots, err error) {

	body := &DescribeGitDepotRequest{
		Action:    "DescribeGitDepot",
		DepotPath: repoPath,
	}
	resp, err := c.Request("POST", endpoint, token, body)
	if err != nil {
		return depot, err
	}
	err = checkResponse(resp)
	if err != nil {
		return depot, err
	}
	var describeGitDepotResponse DescribeGitDepotResponse
	err = c.Unmarshal(resp, &describeGitDepotResponse)
	if err != nil {
		return depot, err
	}
	return describeGitDepotResponse.Response.Depot, nil
}

func GetCurrentUserName(url, token string) (userName string, err error) {

	body := &DescribeCodingCurrentUser{
		Action: "DescribeCodingCurrentUser",
	}
	resp, err := c.Request("POST", endpoint, token, body)
	if err != nil {
		return "", err
	}
	var userInfo UserInfo
	err = c.Unmarshal(resp, &userInfo)
	if err != nil {
		return "", err
	}
	return userInfo.Response.User.Name, nil
}

func GetProjectByName(url, token, projectName string) (project Project, err error) {

	body := &DescribeProjectByNameRequest{
		Action:      "DescribeProjectByName",
		ProjectName: projectName,
	}
	resp, err := c.Request("POST", endpoint, token, body)
	if err != nil {
		return project, err
	}
	err = checkResponse(resp)
	if err != nil {
		return project, err
	}
	var projectInfo DescribeProjectByNameResponse
	err = c.Unmarshal(resp, &projectInfo)
	if err != nil {
		return project, err
	}
	logger.Logger.Debugf("%s项目ID: %d", projectName, projectInfo.Response.Project.Id)
	return projectInfo.Response.Project, nil
}

func GetProjectIdsByNames(projects []string, url, token string) (projectIds []int, err error) {
	for _, projectName := range projects {
		project, err := GetProjectByName(url, token, projectName)
		if err != nil {
			return nil, err
		}
		projectIds = append(projectIds, project.Id)
	}
	return projectIds, nil
}

func GetRepoByProjectIdFetchPage(url, token string, projectId, pageNumber int) (DepotInfo, error) {

	body := &DescribeProjectDepotInfoListRequest{
		Action:     "DescribeProjectDepotInfoList",
		ProjectId:  int64(projectId),
		PageNumber: pageNumber,
		PageSize:   PageSize,
	}
	resp, err := c.Request("POST", endpoint, token, body)
	if err != nil {
		return DepotInfo{}, err
	}
	var depotData DepotInfo
	err = c.Unmarshal(resp, &depotData)
	if err != nil {
		return DepotInfo{}, err
	}
	return depotData, nil
}

func GetTeamRepoFetchPage(url, token string, pageNumber int) (DepotInfo, error) {

	body := &DescribeTeamDepotInfoListRequest{
		Action:     "DescribeTeamDepotInfoList",
		PageNumber: pageNumber,
		PageSize:   PageSize,
	}
	resp, err := c.Request("POST", endpoint, token, body)
	if err != nil {
		return DepotInfo{}, err
	}
	err = checkResponse(resp)
	if err != nil {
		return DepotInfo{}, err
	}
	var depotData DepotInfo
	err = c.Unmarshal(resp, &depotData)
	if err != nil {
		return DepotInfo{}, err
	}
	return depotData, nil
}

func GetRepoByProjectId(url, token string, projectId int) ([]Depots, error) {
	Data := []Depots{}
	page := 1
	for {
		apiResp, err := GetRepoByProjectIdFetchPage(url, token, projectId, page)
		if err != nil {
			return nil, err
		}
		for _, v := range apiResp.Response.DepotData.Depots {
			Data = append(Data, v)
		}
		if apiResp.Response.DepotData.Page.PageNumber >= apiResp.Response.DepotData.Page.TotalPage {
			break
		}
		page++
	}
	logger.Logger.Debugw("获取项目仓库列表成功", "projectId", projectId, "depots", Data)
	return Data, nil
}

func GetDepotListByTeam(url, token string) ([]Depots, error) {
	Data := []Depots{}
	page := 1
	for {
		apiResp, err := GetTeamRepoFetchPage(url, token, page)
		if err != nil {
			return nil, err
		}
		for _, v := range apiResp.Response.DepotData.Depots {
			Data = append(Data, v)
		}
		if apiResp.Response.DepotData.Page.PageNumber >= apiResp.Response.DepotData.Page.TotalPage {
			break
		}
		page++
	}
	return Data, nil
}

func GetReposByProjectIds(url, token string, projectIds []int) (depots []Depots, err error) {
	for _, projectId := range projectIds {
		repos, err := GetRepoByProjectId(url, token, projectId)
		if err != nil {
			return nil, err
		}
		depots = append(depots, repos...)
	}
	return depots, nil
}

func checkResponse(data []byte) error {
	var responseError ResponseError
	err := parseJSON(data, &responseError)
	if err != nil {
		fmt.Println("解析JSON失败:", err)
		return err
	}
	if responseError.Response.Error.Code != "" {
		return fmt.Errorf(responseError.Response.Error.Message)
	}
	return nil
}

func parseJSON(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

func (d Depots) GetRepoPath() string {
	return fmt.Sprintf("%s/%s", d.ProjectName, d.Name)
}

func (d Depots) GetMirrorRepoDirPath() string {
	return fmt.Sprintf("%s%s", d.GetRepoPath(), ".git")
}

func GetDepotListByProjectNames(url, token string, projectNames []string) (depots []Depots, err error) {
	projectIds, err := GetProjectIdsByNames(projectNames, url, token)
	if err != nil {
		return nil, fmt.Errorf("%s 通过项目名获取项目ID获取失败: %s", projectNames, err)
	}
	repos, err := GetReposByProjectIds(url, token, projectIds)
	if err != nil {
		return nil, fmt.Errorf("%d 通过项目ID获取代码仓库列表失败: %s", projectIds, err)
	}
	return repos, nil
}

func GetCloneUrl(repoHttpsURL, userName, token string) string {
	parts := strings.Split(repoHttpsURL, "://")
	cloneURL := parts[0] + "://" + userName + ":" + token + "@" + parts[1]
	return cloneURL
}

func GetDepotListByRepoPath(url, token string, repos []string) (depotList []Depots, err error) {
	for _, repoPath := range repos {
		depot, err := GetDepotByRepoPath(url, token, repoPath)
		if err != nil {
			return nil, fmt.Errorf("查询%s仓库信息失败: %s", repoPath, err)
		}
		depotList = append(depotList, depot)
	}
	return depotList, nil
}

func GetDepotList(migrateType string) ([]Depots, error) {
	logger.Logger.Infof("获取仓库列表中...")
	var depotList []Depots
	var err error
	switch migrateType {
	case ProjectType:
		depotList, err = GetDepotListByProjectNames(SourceURL, SourceToken, Projects)
	case RepoType:
		depotList, err = GetDepotListByRepoPath(SourceURL, SourceToken, Repos)
	case Team:
		depotList, err = GetDepotListByTeam(SourceURL, SourceToken)
	default:
		return nil, fmt.Errorf("未知的迁移类型: %s", migrateType)
	}
	if err != nil {
		return nil, err
	}
	logger.Logger.Debugw("仓库列表清单", "depotList", depotList)
	return depotList, err
}

type GetReleasesListResp struct {
	Response struct {
		ReleasePageList struct {
			Releases   []Releases `json:"Releases"`
			TotalCount int        `json:"TotalCount"`
		} `json:"ReleasePageList"`
		RequestId string `json:"RequestId"`
	} `json:"Response"`
}

type Releases struct {
	Body               string   `json:"Body"`
	CommitSha          string   `json:"CommitSha"`
	CreatedAt          int64    `json:"CreatedAt"`
	CreatorId          int      `json:"CreatorId"`
	DepotId            int      `json:"DepotId"`
	Html               string   `json:"Html"`
	Id                 int      `json:"Id"`
	Pre                bool     `json:"Pre"`
	ProjectId          int      `json:"ProjectId"`
	ReleaseId          int      `json:"ReleaseId"`
	TagName            string   `json:"TagName"`
	TargetCommitBranch string   `json:"TargetCommitBranch"`
	Title              string   `json:"Title"`
	UpdatedAt          int64    `json:"UpdatedAt"`
	Iid                int      `json:"iid"`
	ImageDownloadUrl   []string `json:"ImageDownloadUrl"`
	ReleaseAttachment  []struct {
		AttachmentName        string `json:"AttachmentName"`
		AttachmentDownloadUrl string `json:"AttachmentDownloadUrl"`
		AttachmentSize        int    `json:"AttachmentSize"`
	} `json:"ReleaseAttachment"`
}

type GetReleasesListReq struct {
	Action          string `json:"Action"`
	DepotId         int    `json:"DepotId"`
	DepotPath       string `json:"DepotPath"`
	FromDate        string `json:"FromDate"`
	PageNumber      int    `json:"PageNumber"`
	PageSize        int    `json:"PageSize"`
	Status          string `json:"Status"`
	TagName         string `json:"TagName"`
	ToDate          string `json:"ToDate"`
	ShowResourceUrl bool   `json:"ShowResourceUrl"`
}

func GetReleasesFetchPage(repoID, page int) (GetReleasesListResp, error) {
	var getReleasesListResp GetReleasesListResp
	body := GetReleasesListReq{
		Action:          "DescribeGitReleases",
		DepotId:         repoID,
		PageNumber:      page,
		PageSize:        100,
		ShowResourceUrl: true,
	}
	resp, _, _, err := c2.RequestV4(http.MethodPost, "", body)
	if err != nil {
		return getReleasesListResp, err
	}
	logger.Logger.Debugw("获取仓库Release列表", "resp", string(resp))
	err = checkResponse(resp)
	if err != nil {
		return getReleasesListResp, err
	}
	err = json.Unmarshal(resp, &getReleasesListResp)
	if err != nil {
		return getReleasesListResp, err
	}
	return getReleasesListResp, nil
}

func GetReleasesList(repoID int) ([]Releases, error) {
	var releases []Releases
	page := 1
	for {
		resp, err := GetReleasesFetchPage(repoID, page)
		if err != nil {
			return releases, err
		}
		releases = append(releases, resp.Response.ReleasePageList.Releases...)
		if len(releases) < resp.Response.ReleasePageList.TotalCount {
			page++
		} else {
			break
		}
	}
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].Iid < releases[j].Iid
	})
	return releases, nil
}
