package vcs

import (
	api "ccrctl/pkg/api/gitea"
	"testing"
)

// TestGiteaCovertToVcs_VisibilityMapping 测试 Gitea 仓库可见性映射逻辑
func TestGiteaCovertToVcs_VisibilityMapping(t *testing.T) {
	tests := []struct {
		name            string
		repo            api.Repo
		expectedPrivate bool
		description     string
	}{
		{
			name: "公开仓库",
			repo: api.Repo{
				FullName:    "user/public-repo",
				Name:        "public-repo",
				CloneUrl:    "https://gitea.example.com/user/public-repo.git",
				Private:     false,
				Internal:    false,
				Description: "公开仓库",
			},
			expectedPrivate: false,
			description:     "Private=false, Internal=false 应该返回 false",
		},
		{
			name: "私有仓库",
			repo: api.Repo{
				FullName:    "user/private-repo",
				Name:        "private-repo",
				CloneUrl:    "https://gitea.example.com/user/private-repo.git",
				Private:     true,
				Internal:    false,
				Description: "私有仓库",
			},
			expectedPrivate: true,
			description:     "Private=true, Internal=false 应该返回 true",
		},
		{
			name: "内部仓库",
			repo: api.Repo{
				FullName:    "user/internal-repo",
				Name:        "internal-repo",
				CloneUrl:    "https://gitea.example.com/user/internal-repo.git",
				Private:     false,
				Internal:    true,
				Description: "内部仓库",
			},
			expectedPrivate: true,
			description:     "Private=false, Internal=true 应该返回 true (映射为 private)",
		},
		{
			name: "内部且私有仓库",
			repo: api.Repo{
				FullName:    "user/internal-private-repo",
				Name:        "internal-private-repo",
				CloneUrl:    "https://gitea.example.com/user/internal-private-repo.git",
				Private:     true,
				Internal:    true,
				Description: "内部私有仓库",
			},
			expectedPrivate: true,
			description:     "Private=true, Internal=true 应该返回 true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoList := []api.Repo{tt.repo}
			vcsList := GiteaCovertToVcs(repoList)

			if len(vcsList) != 1 {
				t.Fatalf("期望返回1个VCS对象，实际返回%d个", len(vcsList))
			}

			vcs := vcsList[0].(*GiteaVcs)

			// 验证 Internal 字段是否正确传递
			if vcs.Internal != tt.repo.Internal {
				t.Errorf("Internal 字段映射错误: 期望 %v，实际 %v",
					tt.repo.Internal, vcs.Internal)
			}

			// 验证 Private 字段是否正确传递
			if vcs.Private != tt.repo.Private {
				t.Errorf("Private 字段映射错误: 期望 %v，实际 %v",
					tt.repo.Private, vcs.Private)
			}

			// 验证 GetRepoPrivate() 方法返回值
			actualPrivate := vcs.GetRepoPrivate()
			if actualPrivate != tt.expectedPrivate {
				t.Errorf("%s: GetRepoPrivate() 期望返回 %v，实际返回 %v",
					tt.description, tt.expectedPrivate, actualPrivate)
			}

			// 验证其他字段是否正确转换
			if vcs.RepoPath != tt.repo.FullName {
				t.Errorf("RepoPath 转换错误: 期望 %s，实际 %s", tt.repo.FullName, vcs.RepoPath)
			}

			if vcs.RepoName != tt.repo.Name {
				t.Errorf("RepoName 转换错误: 期望 %s，实际 %s", tt.repo.Name, vcs.RepoName)
			}

			if vcs.httpURL != tt.repo.CloneUrl {
				t.Errorf("httpURL 转换错误: 期望 %s，实际 %s", tt.repo.CloneUrl, vcs.httpURL)
			}

			if vcs.Desc != tt.repo.Description {
				t.Errorf("Desc 转换错误: 期望 %s，实际 %s", tt.repo.Description, vcs.Desc)
			}
		})
	}
}

// TestGiteaVcs_GetRepoPrivate 测试 GetRepoPrivate 方法的各种组合
func TestGiteaVcs_GetRepoPrivate(t *testing.T) {
	testCases := []struct {
		name     string
		private  bool
		internal bool
		expected bool
	}{
		{"false_false", false, false, false},
		{"false_true", false, true, true}, // 关键测试：Internal=true 应该返回 true
		{"true_false", true, false, true},
		{"true_true", true, true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vcs := &GiteaVcs{
				Private:  tc.private,
				Internal: tc.internal,
			}

			result := vcs.GetRepoPrivate()
			if result != tc.expected {
				t.Errorf("Private=%v, Internal=%v: 期望 GetRepoPrivate()=%v，实际=%v",
					tc.private, tc.internal, tc.expected, result)
			}
		})
	}
}

// TestGiteaCovertToVcs_MultipleRepos 测试批量转换多个仓库
func TestGiteaCovertToVcs_MultipleRepos(t *testing.T) {
	repoList := []api.Repo{
		{
			FullName:    "user/public-repo",
			Name:        "public-repo",
			CloneUrl:    "https://gitea.example.com/user/public-repo.git",
			Private:     false,
			Internal:    false,
			Description: "公开仓库",
		},
		{
			FullName:    "user/internal-repo",
			Name:        "internal-repo",
			CloneUrl:    "https://gitea.example.com/user/internal-repo.git",
			Private:     false,
			Internal:    true,
			Description: "内部仓库",
		},
		{
			FullName:    "user/private-repo",
			Name:        "private-repo",
			CloneUrl:    "https://gitea.example.com/user/private-repo.git",
			Private:     true,
			Internal:    false,
			Description: "私有仓库",
		},
	}

	vcsList := GiteaCovertToVcs(repoList)

	if len(vcsList) != 3 {
		t.Fatalf("期望返回3个VCS对象，实际返回%d个", len(vcsList))
	}

	// 验证第一个仓库（公开）
	vcs1 := vcsList[0].(*GiteaVcs)
	if vcs1.GetRepoPrivate() != false {
		t.Error("公开仓库的 GetRepoPrivate() 应该返回 false")
	}
	if vcs1.Private != false || vcs1.Internal != false {
		t.Error("公开仓库的 Private 和 Internal 字段都应该是 false")
	}

	// 验证第二个仓库（内部，应该映射为私有）
	vcs2 := vcsList[1].(*GiteaVcs)
	if vcs2.GetRepoPrivate() != true {
		t.Error("内部仓库的 GetRepoPrivate() 应该返回 true")
	}
	if vcs2.Private != false || vcs2.Internal != true {
		t.Error("内部仓库的 Private 应该是 false，Internal 应该是 true")
	}

	// 验证第三个仓库（私有）
	vcs3 := vcsList[2].(*GiteaVcs)
	if vcs3.GetRepoPrivate() != true {
		t.Error("私有仓库的 GetRepoPrivate() 应该返回 true")
	}
	if vcs3.Private != true || vcs3.Internal != false {
		t.Error("私有仓库的 Private 应该是 true，Internal 应该是 false")
	}
}

// TestGiteaCovertToVcs_EmptyRepoList 测试空仓库列表
func TestGiteaCovertToVcs_EmptyRepoList(t *testing.T) {
	repoList := []api.Repo{}
	vcsList := GiteaCovertToVcs(repoList)

	if len(vcsList) != 0 {
		t.Errorf("空仓库列表应该返回空VCS列表，实际返回%d个", len(vcsList))
	}
}

// TestGiteaVcs_Methods 测试 GiteaVcs 结构体的各种方法
func TestGiteaVcs_Methods(t *testing.T) {
	repo := api.Repo{
		FullName:    "org/team/test-repo",
		Name:        "test-repo",
		CloneUrl:    "https://gitea.example.com/org/team/test-repo.git",
		Private:     false,
		Internal:    true, // 内部仓库
		Description: "测试仓库描述",
	}

	vcsList := GiteaCovertToVcs([]api.Repo{repo})
	vcs := vcsList[0]

	// 测试 GetRepoPrivate 方法
	if !vcs.GetRepoPrivate() {
		t.Error("内部仓库的 GetRepoPrivate() 应该返回 true")
	}

	// 测试 GetRepoPath 方法
	if vcs.GetRepoPath() != "org/team/test-repo" {
		t.Errorf("GetRepoPath() 期望 'org/team/test-repo'，实际 '%s'", vcs.GetRepoPath())
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

	// 测试 GetSubGroup 方法
	subGroup := vcs.GetSubGroup()
	if subGroup.Name != "org/team" {
		t.Errorf("GetSubGroup().Name 期望 'org/team'，实际 '%s'", subGroup.Name)
	}
}

// TestGiteaVcs_GetSubGroup 测试 GetSubGroup 方法的各种路径情况
func TestGiteaVcs_GetSubGroup(t *testing.T) {
	testCases := []struct {
		name             string
		repoPath         string
		expectedSubGroup string
	}{
		{
			name:             "两层路径",
			repoPath:         "user/repo",
			expectedSubGroup: "user",
		},
		{
			name:             "三层路径",
			repoPath:         "org/team/repo",
			expectedSubGroup: "org/team",
		},
		{
			name:             "单层路径",
			repoPath:         "repo",
			expectedSubGroup: "",
		},
		{
			name:             "多层嵌套路径",
			repoPath:         "org/dept/team/project/repo",
			expectedSubGroup: "org/dept/team/project",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vcs := &GiteaVcs{
				RepoPath: tc.repoPath,
			}

			subGroup := vcs.GetSubGroup()
			if subGroup.Name != tc.expectedSubGroup {
				t.Errorf("路径 '%s': 期望子组织 '%s'，实际 '%s'",
					tc.repoPath, tc.expectedSubGroup, subGroup.Name)
			}
		})
	}
}
