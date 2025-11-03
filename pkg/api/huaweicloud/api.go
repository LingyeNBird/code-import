package huaweicloud

import (
	"ccrctl/pkg/config"
	"ccrctl/pkg/logger"
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	codehub "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/codehub/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/codehub/v3/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/codehub/v3/region"
	projectman "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/projectman/v4"
	v4model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/projectman/v4/model"
	v4region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/projectman/v4/region"
)

var (
	client        *codehub.CodeHubClient
	projectClient *projectman.ProjectManClient
)

// InitClient 初始化华为云CodeArts客户端
func InitClient() error {
	// 从配置获取认证信息
	ak := config.Cfg.GetString("source.ak")
	sk := config.Cfg.GetString("source.sk")
	regionName := config.Cfg.GetString("source.region")

	// 检查必需的配置
	if ak == "" || sk == "" {
		return fmt.Errorf("华为云AK/SK未配置")
	}

	// 创建认证信息（使用AK/SK）
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).Build()
	// 创建CodeHub客户端
	regionValue, err := region.SafeValueOf(regionName)
	if err != nil {
		return fmt.Errorf("无效的区域名称: %w", err)
	}

	client = codehub.NewCodeHubClient(
		codehub.CodeHubClientBuilder().
			WithRegion(regionValue).
			WithCredential(auth).
			Build())

	return nil
}

func initProjectClient() error {
	// 从配置获取认证信息
	ak := config.Cfg.GetString("source.ak")
	sk := config.Cfg.GetString("source.sk")
	regionName := config.Cfg.GetString("source.region")

	// 检查必需的配置
	if ak == "" || sk == "" {
		return fmt.Errorf("华为云AK/SK未配置")
	}

	// 创建认证信息（使用AK/SK）
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()
	// 创建Project客户端
	regionValue, err := v4region.SafeValueOf(regionName)
	if err != nil {
		return fmt.Errorf("无效的区域名称: %w", err)
	}

	projectClient = projectman.NewProjectManClient(
		projectman.ProjectManClientBuilder().
			WithRegion(regionValue).
			WithCredential(auth).
			Build())
	return nil
}

// GetRepositories 获取华为云CodeArts仓库列表（支持分页）
func GetRepositories() ([]model.RepoInfoV2, error) {
	if client == nil {
		// 尝试初始化客户端
		if err := InitClient(); err != nil {
			return nil, fmt.Errorf("华为云客户端未初始化: %w", err)
		}
	}

	var allRepositories []model.RepoInfoV2
	pageIndex := int32(1)  // 华为云仓库API使用pageIndex，从1开始
	pageSize := int32(100) // 每页获取100个仓库
	maxPages := 100        // 最大页数限制，防止无限循环

	for page := 0; page < maxPages; page++ {
		// 创建请求
		request := &model.ListUserAllRepositoriesRequest{
			PageIndex: &pageIndex,
			PageSize:  &pageSize,
		}

		logger.Logger.Debugf("正在获取华为云CodeArts仓库列表，第%d页，pageIndex: %d, pageSize: %d", page+1, pageIndex, pageSize)

		// 调用API获取仓库列表
		response, err := client.ListUserAllRepositories(request)
		if err != nil {
			return nil, fmt.Errorf("获取华为云CodeArts仓库列表失败: %w", err)
		}

		// 检查响应是否有效
		if response.Result == nil || response.Result.Repositories == nil {
			logger.Logger.Debugf("仓库列表为空，结束分页")
			break
		}

		// 添加当前页的仓库到结果中
		currentPageRepos := *response.Result.Repositories
		allRepositories = append(allRepositories, currentPageRepos...)

		logger.Logger.Debugf("当前页获取到 %d 个仓库，总计已获取 %d 个仓库", len(currentPageRepos), len(allRepositories))

		// 检查是否还有更多数据
		// 如果当前页返回的仓库数量小于pageSize，说明已经是最后一页
		if len(currentPageRepos) < int(pageSize) {
			logger.Logger.Debugf("已获取所有仓库，结束分页")
			break
		}

		// 如果有total字段，可以用来判断是否还有更多数据
		if response.Result.Total != nil {
			totalCount := *response.Result.Total
			if len(allRepositories) >= int(totalCount) {
				logger.Logger.Debugf("已获取所有仓库 (总数: %d)，结束分页", totalCount)
				break
			}
		}

		// 更新pageIndex到下一页
		pageIndex++
	}

	logger.Logger.Debugf("成功获取华为云CodeArts仓库列表，共 %d 个仓库", len(allRepositories))
	return allRepositories, nil
}

func GetProjects() (map[string]string, error) {
	projectsMap := make(map[string]string)
	if projectClient == nil {
		// 尝试初始化客户端
		if err := initProjectClient(); err != nil {
			return nil, fmt.Errorf("华为云客户端未初始化: %w", err)
		}
	}

	var projects []v4model.ListProjectsV4ResponseBodyProjects
	offset := int32(0)
	limit := int32(100) // 每页获取100个项目，可以根据需要调整
	maxPages := 100     // 最大页数限制，防止无限循环

	for page := 0; page < maxPages; page++ {
		// 创建请求
		request := &v4model.ListProjectsV4Request{
			Offset: offset,
			Limit:  limit,
		}

		logger.Logger.Debugf("正在获取华为云CodeArts项目列表，第%d页，offset: %d, limit: %d", page+1, offset, limit)

		// 调用API获取项目列表
		response, err := projectClient.ListProjectsV4(request)
		if err != nil {
			logger.Logger.Errorf("获取华为云CodeArts项目列表失败: %v", err)
			return projectsMap, fmt.Errorf("获取华为云CodeArts项目列表失败: %w", err)
		}

		// 检查响应是否有效
		if response.Projects == nil {
			logger.Logger.Debugf("项目列表为空，结束分页")
			break
		}

		// 添加当前页的项目到结果中
		currentPageProjects := *response.Projects
		projects = append(projects, currentPageProjects...)

		logger.Logger.Debugf("当前页获取到 %d 个项目，总计已获取 %d 个项目", len(currentPageProjects), len(projects))

		// 检查是否还有更多数据
		// 如果当前页返回的项目数量小于limit，说明已经是最后一页
		if len(currentPageProjects) < int(limit) {
			logger.Logger.Debugf("已获取所有项目，结束分页")
			break
		}

		// 如果有total字段，可以用来判断是否还有更多数据
		if response.Total != nil {
			totalCount := *response.Total
			if len(projects) >= int(totalCount) {
				logger.Logger.Debugf("已获取所有项目 (总数: %d)，结束分页", totalCount)
				break
			}
		}

		// 更新offset到下一页
		offset += limit
	}

	logger.Logger.Debugf("成功获取华为云CodeArts项目列表，共 %d 个项目", len(projects))
	for _, project := range projects {
		projectsMap[*project.ProjectId] = *project.ProjectName
	}
	return projectsMap, nil
}
