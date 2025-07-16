package github

import (
	"ccrctl/pkg/config"
	"ccrctl/pkg/logger"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

var client *github.Client

func init() {
	if config.Cfg.GetString("source.platform") == "github" {
		// 创建一个上下文
		ctx := context.Background()

		// 替换为你的GitHub访问令牌
		token := config.Cfg.GetString("source.token")

		// 创建一个OAuth2客户端
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)

		// 创建一个GitHub客户端
		client = github.NewClient(tc)
	}
}

func GetRepos() ([]*github.Repository, error) {
	// 创建一个上下文
	ctx := context.Background()

	opt := &github.RepositoryListByAuthenticatedUserOptions{}
	var allRepos []*github.Repository

	for {
		repos, resp, err := client.Repositories.ListByAuthenticatedUser(ctx, opt)
		if err != nil {
			logger.Logger.Fatalf("Failed to list repositories: %v", err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	if config.Cfg.GetBool("migrate.exclude_github_fork") {
		var filteredRepos []*github.Repository
		for _, repo := range allRepos {
			if !*repo.Fork {
				filteredRepos = append(filteredRepos, repo)
			}
		}
		allRepos = filteredRepos
	}
	return allRepos, nil
}

func GetUserName() string {
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		logger.Logger.Fatalf("Failed to get github user: %v", err)
	}
	return user.GetName()
}

func GetReleases(owner, repo string) ([]*github.RepositoryRelease, error) {
	var allReleases []*github.RepositoryRelease
	ctx := context.Background()
	opts := &github.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	for {
		releases, resp, err := client.Repositories.ListReleases(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}
		allReleases = append(allReleases, releases...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// 基于 PublishedAt 正序排序
	sort.Slice(allReleases, func(i, j int) bool {
		return allReleases[i].PublishedAt.Time.Before(allReleases[j].PublishedAt.Time)
	})

	return allReleases, nil
}

func DownloadReleaseAsset(owner, repo string, assetID int64) ([]byte, error) {
	ctx := context.Background()
	asset, _, err := client.Repositories.DownloadReleaseAsset(ctx, owner, repo, assetID, http.DefaultClient)
	if err != nil {
		return nil, fmt.Errorf("下载附件失败: %v", err)
	}
	defer asset.Close()
	return io.ReadAll(asset)
}

type ReleaseAssetLink struct {
	URL     string
	AssetID int64
}

func ExtractDownloadLinksFromRelease(owner, repo string, releaseID int64) ([]ReleaseAssetLink, error) {
	ctx := context.Background()
	assets, _, err := client.Repositories.ListReleaseAssets(ctx, owner, repo, releaseID, nil)
	if err != nil {
		return nil, fmt.Errorf("获取附件列表失败: %v", err)
	}

	release, _, err := client.Repositories.GetRelease(ctx, owner, repo, releaseID)
	if err != nil {
		return nil, fmt.Errorf("获取release信息失败: %v", err)
	}

	re := regexp.MustCompile(`!\[.*?\]\((.*?\.(?:png|jpg|jpeg|gif|zip|tar|gz))\)`)
	matches := re.FindAllStringSubmatch(release.GetBody(), -1)
	var links []ReleaseAssetLink

	for _, match := range matches {
		if len(match) > 1 {
			url := match[1]
			for _, asset := range assets {
				if strings.Contains(url, asset.GetName()) {
					links = append(links, ReleaseAssetLink{
						URL:     url,
						AssetID: asset.GetID(),
					})
					break
				}
			}
		}
	}
	return links, nil
}

// ListUploads 获取 release 中的所有附件和图片
func ListUploads(projectID string) (files map[string]int64, err error) {
	parts := strings.Split(projectID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("无效的 projectID: %s", projectID)
	}
	owner := parts[0]
	repo := parts[1]

	ctx := context.Background()
	releases, err := GetReleases(owner, repo)
	if err != nil {
		return nil, fmt.Errorf("获取 releases 失败: %v", err)
	}

	files = make(map[string]int64)
	for _, release := range releases {
		assets, _, err := client.Repositories.ListReleaseAssets(ctx, owner, repo, release.GetID(), nil)
		if err != nil {
			return nil, fmt.Errorf("获取 release assets 失败: %v", err)
		}

		for _, asset := range assets {
			files[asset.GetName()] = asset.GetID()
		}
	}
	return files, nil
}

// DownloadFile 下载指定 ID 的文件
func DownloadFile(projectID string, fileID int64) (data []byte, err error) {
	parts := strings.Split(projectID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("无效的 projectID: %s", projectID)
	}
	owner := parts[0]
	repo := parts[1]

	return DownloadReleaseAsset(owner, repo, fileID)
}
