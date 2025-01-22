package aliyun

import (
	C "ccrctl/pkg/config"
	"ccrctl/pkg/logger"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	devops20210625 "github.com/alibabacloud-go/devops-20210625/v5/client"
	"github.com/alibabacloud-go/tea/tea"
)

var client *devops20210625.Client

func init() {
	var err error
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(C.Cfg.GetString("source.ak")),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(C.Cfg.GetString("source.as")),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/devops
	config.Endpoint = tea.String(C.Cfg.GetString("source.endpoint"))
	client, err = devops20210625.NewClient(config)
	if err != nil {
		logger.Logger.Fatalf("Failed to create client: %v", err)
	}
}

func ListRepository(organizationId string) []*devops20210625.ListRepositoriesResponseBodyResult {
	var repos []*devops20210625.ListRepositoriesResponseBodyResult
	page := 1
	for {
		res, err := client.ListRepositories(&devops20210625.ListRepositoriesRequest{
			OrganizationId: tea.String(organizationId),
			PerPage:        tea.Int64(100),
			Page:           tea.Int64(int64(page)),
		})
		if err != nil {
			logger.Logger.Fatalf("Failed to list repositories: %v", err)
		}
		repos = append(repos, res.Body.Result...)
		if *res.Body.Total <= int64(page)*100 {
			break
		}
		page++
	}
	return repos
}
