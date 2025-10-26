package vcs

import (
	api "ccrctl/pkg/api/gitee"
	"fmt"
	"testing"
)

func TestGiteeCovertToVcs_InternalRepoHandling(t *testing.T) {
	tests := []struct {
		name            string
		repo            api.Repo
		expectedPrivate bool
		description     string
	}{
		{
			name: "普通公开仓库",
			repo: api.Repo{
				FullName:    "user/public-repo",
				Name:        "public-repo",
				HtmlUrl:     "https://gitee.com/user/public-repo",
				Private:     false,
				Internal:    false,
				Description: "公开仓库",
			},
			expectedPrivate: false,
			description:     "Private=false, Internal=false 应该保持 Private=false",
		},
		{
			name: "普通私有仓库",
			repo: api.Repo{
				FullName:    "user/private-repo",
				Name:        "private-repo",
				HtmlUrl:     "https://gitee.com/user/private-repo",
				Private:     true,
				Internal:    false,
				Description: "私有仓库",
			},
			expectedPrivate: true,
			description:     "Private=true, Internal=false 应该保持 Private=true",
		},
		{
			name: "内部仓库但未标记为私有",
			repo: api.Repo{
				FullName:    "user/internal-repo",
				Name:        "internal-repo",
				HtmlUrl:     "https://gitee.com/user/internal-repo",
				Private:     false,
				Internal:    true,
				Description: "内部仓库",
			},
			expectedPrivate: true,
			description:     "Private=false, Internal=true 应该自动设置 Private=true",
		},
		{
			name: "内部且私有仓库",
			repo: api.Repo{
				FullName:    "user/internal-private-repo",
				Name:        "internal-private-repo",
				HtmlUrl:     "https://gitee.com/user/internal-private-repo",
				Private:     true,
				Internal:    true,
				Description: "内部私有仓库",
			},
			expectedPrivate: true,
			description:     "Private=true, Internal=true 应该保持 Private=true",
		},
		{
			name: "边界情况_空描述",
			repo: api.Repo{
				FullName:    "user/empty-desc-internal",
				Name:        "empty-desc-internal",
				HtmlUrl:     "https://gitee.com/user/empty-desc-internal",
				Private:     false,
				Internal:    true,
				Description: "",
			},
			expectedPrivate: true,
			description:     "内部仓库即使描述为空也应该设置为私有",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoList := []api.Repo{tt.repo}
			vcsList := GiteeCovertToVcs(repoList)

			if len(vcsList) != 1 {
				t.Fatalf("期望返回1个VCS对象，实际返回%d个", len(vcsList))
			}

			vcs := vcsList[0].(*GiteeVcs)

			// 验证 Private 字段
			if vcs.Private != tt.expectedPrivate {
				t.Errorf("%s: 期望 Private=%v，实际 Private=%v",
					tt.description, tt.expectedPrivate, vcs.Private)
			}

			// 验证其他字段是否正确转换
			if vcs.RepoPath != tt.repo.FullName {
				t.Errorf("RepoPath 转换错误: 期望 %s，实际 %s", tt.repo.FullName, vcs.RepoPath)
			}

			if vcs.RepoName != tt.repo.Name {
				t.Errorf("RepoName 转换错误: 期望 %s，实际 %s", tt.repo.Name, vcs.RepoName)
			}

			if vcs.httpURL != tt.repo.HtmlUrl {
				t.Errorf("httpURL 转换错误: 期望 %s，实际 %s", tt.repo.HtmlUrl, vcs.httpURL)
			}

			if vcs.Desc != tt.repo.Description {
				t.Errorf("Desc 转换错误: 期望 %s，实际 %s", tt.repo.Description, vcs.Desc)
			}

			if vcs.RepoType != "git" {
				t.Errorf("RepoType 应该是 'git'，实际 '%s'", vcs.RepoType)
			}
		})
	}
}

func TestGiteeCovertToVcs_MultipleRepos(t *testing.T) {
	// 测试多个仓库的批量转换
	repoList := []api.Repo{
		{
			FullName:    "user/public-repo",
			Name:        "public-repo",
			HtmlUrl:     "https://gitee.com/user/public-repo",
			Private:     false,
			Internal:    false,
			Description: "公开仓库",
		},
		{
			FullName:    "user/internal-repo",
			Name:        "internal-repo",
			HtmlUrl:     "https://gitee.com/user/internal-repo",
			Private:     false,
			Internal:    true,
			Description: "内部仓库",
		},
		{
			FullName:    "user/private-repo",
			Name:        "private-repo",
			HtmlUrl:     "https://gitee.com/user/private-repo",
			Private:     true,
			Internal:    false,
			Description: "私有仓库",
		},
	}

	vcsList := GiteeCovertToVcs(repoList)

	if len(vcsList) != 3 {
		t.Fatalf("期望返回3个VCS对象，实际返回%d个", len(vcsList))
	}

	// 验证第一个仓库（公开）
	vcs1 := vcsList[0].(*GiteeVcs)
	if vcs1.Private != false {
		t.Error("公开仓库应该保持 Private=false")
	}

	// 验证第二个仓库（内部，应该自动设置为私有）
	vcs2 := vcsList[1].(*GiteeVcs)
	if vcs2.Private != true {
		t.Error("内部仓库应该自动设置 Private=true")
	}

	// 验证第三个仓库（私有）
	vcs3 := vcsList[2].(*GiteeVcs)
	if vcs3.Private != true {
		t.Error("私有仓库应该保持 Private=true")
	}
}

func TestGiteeCovertToVcs_EmptyRepoList(t *testing.T) {
	// 测试空仓库列表
	repoList := []api.Repo{}
	vcsList := GiteeCovertToVcs(repoList)

	if len(vcsList) != 0 {
		t.Errorf("空仓库列表应该返回空VCS列表，实际返回%d个", len(vcsList))
	}
}

func TestGiteeCovertToVcs_InternalPrivateLogic(t *testing.T) {
	// 专门测试 Internal 和 Private 字段的逻辑组合
	testCases := []struct {
		name            string
		private         bool
		internal        bool
		expectedPrivate bool
	}{
		{"false_false", false, false, false},
		{"false_true", false, true, true},   // 关键测试：Internal=true 应该导致 Private=true
		{"true_false", true, false, true},
		{"true_true", true, true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := api.Repo{
				FullName:    "test/repo",
				Name:        "repo",
				HtmlUrl:     "https://gitee.com/test/repo",
				Private:     tc.private,
				Internal:    tc.internal,
				Description: "测试仓库",
			}

			vcsList := GiteeCovertToVcs([]api.Repo{repo})
			vcs := vcsList[0].(*GiteeVcs)

			if vcs.Private != tc.expectedPrivate {
				t.Errorf("Private=%v, Internal=%v 应该得到 Private=%v，实际得到 %v",
					tc.private, tc.internal, tc.expectedPrivate, vcs.Private)
			}
		})
	}
}

// 基准测试：测试转换性能
func BenchmarkGiteeCovertToVcs(b *testing.B) {
	// 创建测试数据
	repoList := make([]api.Repo, 100)
	for i := 0; i < 100; i++ {
		repoList[i] = api.Repo{
			FullName:    fmt.Sprintf("user/repo-%d", i),
			Name:        fmt.Sprintf("repo-%d", i),
			HtmlUrl:     fmt.Sprintf("https://gitee.com/user/repo-%d", i),
			Private:     i%2 == 0,
			Internal:    i%3 == 0,
			Description: fmt.Sprintf("测试仓库 %d", i),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GiteeCovertToVcs(repoList)
	}
}

// 测试 GiteeVcs 结构体的方法
func TestGiteeVcs_Methods(t *testing.T) {
	repo := api.Repo{
		FullName:    "user/test-repo",
		Name:        "test-repo",
		HtmlUrl:     "https://gitee.com/user/test-repo",
		Private:     false,
		Internal:    true, // 内部仓库
		Description: "测试仓库描述",
	}

	vcsList := GiteeCovertToVcs([]api.Repo{repo})
	vcs := vcsList[0]

	// 测试 GetRepoPrivate 方法
	if !vcs.GetRepoPrivate() {
		t.Error("内部仓库的 GetRepoPrivate() 应该返回 true")
	}

	// 测试 GetRepoPath 方法
	if vcs.GetRepoPath() != "user/test-repo" {
		t.Errorf("GetRepoPath() 期望 'user/test-repo'，实际 '%s'", vcs.GetRepoPath())
	}

	// 测试 GetRepoName 方法
	if vcs.GetRepoName() != "test-repo" {
		t.Errorf("GetRepoName() 期望 'test-repo'，实际 '%s'", vcs.GetRepoName())
	}

	// 测试 GetRepoType 方法
	if vcs.GetRepoType() != "Git" {
		t.Errorf("GetRepoType() 期望 'Git'，实际 '%s'", vcs.GetRepoType())
	}

	// 测试 GetRepoDescription 方法
	if vcs.GetRepoDescription() != "测试仓库描述" {
		t.Errorf("GetRepoDescription() 期望 '测试仓库描述'，实际 '%s'", vcs.GetRepoDescription())
	}
}