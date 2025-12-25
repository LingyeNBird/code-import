package migrate

import (
	"ccrctl/pkg/config"
	"ccrctl/pkg/vcs"
	"testing"

	"github.com/spf13/viper"
)

// MockVCS 模拟 VCS 接口用于测试
type MockVCS struct {
	repoPath string
	repoName string
	subGroup *vcs.SubGroup
}

func (m *MockVCS) GetRepoPath() string           { return m.repoPath }
func (m *MockVCS) GetRepoName() string           { return m.repoName }
func (m *MockVCS) GetSubGroup() *vcs.SubGroup    { return m.subGroup }
func (m *MockVCS) GetRepoType() string           { return "git" }
func (m *MockVCS) GetCloneUrl() string           { return "" }
func (m *MockVCS) GetUserName() string           { return "" }
func (m *MockVCS) GetToken() string              { return "" }
func (m *MockVCS) Clone() error                  { return nil }
func (m *MockVCS) GetRepoPrivate() bool          { return false }
func (m *MockVCS) GetReleases() []vcs.Releases   { return nil }
func (m *MockVCS) GetProjectID() string          { return "" }
func (m *MockVCS) GetRepoDescription() string    { return "" }
func (m *MockVCS) ListRepos() ([]vcs.VCS, error) { return nil, nil }
func (m *MockVCS) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]vcs.Attachment, error) {
	return nil, nil
}

// setupTestConfig 设置测试用的配置
func setupTestConfig(sourceRepo []string) *viper.Viper {
	cfg := viper.New()
	if len(sourceRepo) > 0 {
		cfg.Set("source.repo", sourceRepo)
	}
	return cfg
}

// TestFilterReposByConfigList_EmptyConfig 测试配置为空的情况
func TestFilterReposByConfigList_EmptyConfig(t *testing.T) {
	// 保存原有配置
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	// 设置空配置
	config.Cfg = setupTestConfig(nil)

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "org1/project1/repo1"},
		&MockVCS{repoPath: "org1/project1/repo2"},
		&MockVCS{repoPath: "org2/project2/repo3"},
	}

	result, notFoundCount := filterReposByConfigList(depotList)

	if len(result) != 3 {
		t.Errorf("空配置应返回完整列表，期望 3 个仓库，实际 %d 个", len(result))
	}
	if notFoundCount != 0 {
		t.Errorf("空配置未找到数量应为 0，实际 %d", notFoundCount)
	}
}

// TestFilterReposByConfigList_SingleRepo 测试单个仓库过滤
func TestFilterReposByConfigList_SingleRepo(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{"org1/project1/repo1"})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "org1/project1/repo1"},
		&MockVCS{repoPath: "org1/project1/repo2"},
		&MockVCS{repoPath: "org2/project2/repo3"},
	}

	result, notFoundCount := filterReposByConfigList(depotList)

	if len(result) != 1 {
		t.Errorf("应该只返回1个仓库，实际 %d 个", len(result))
	}
	if result[0].GetRepoPath() != "org1/project1/repo1" {
		t.Errorf("仓库路径不匹配，期望 org1/project1/repo1，实际 %s", result[0].GetRepoPath())
	}
	if notFoundCount != 0 {
		t.Errorf("未找到数量应为 0，实际 %d", notFoundCount)
	}
}

// TestFilterReposByConfigList_MultipleRepos 测试多个仓库过滤
func TestFilterReposByConfigList_MultipleRepos(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{
		"org1/project1/repo1",
		"org2/project2/repo3",
	})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "org1/project1/repo1"},
		&MockVCS{repoPath: "org1/project1/repo2"},
		&MockVCS{repoPath: "org2/project2/repo3"},
		&MockVCS{repoPath: "org3/project3/repo4"},
	}

	result, notFoundCount := filterReposByConfigList(depotList)

	if len(result) != 2 {
		t.Errorf("应该返回2个仓库，实际 %d 个", len(result))
	}

	resultPaths := make(map[string]bool)
	for _, r := range result {
		resultPaths[r.GetRepoPath()] = true
	}

	if !resultPaths["org1/project1/repo1"] {
		t.Error("结果应包含 org1/project1/repo1")
	}
	if !resultPaths["org2/project2/repo3"] {
		t.Error("结果应包含 org2/project2/repo3")
	}
	if resultPaths["org1/project1/repo2"] {
		t.Error("结果不应包含 org1/project1/repo2")
	}
	if notFoundCount != 0 {
		t.Errorf("未找到数量应为 0，实际 %d", notFoundCount)
	}
}

// TestFilterReposByConfigList_WithWhitespace 测试带空格的配置
func TestFilterReposByConfigList_WithWhitespace(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{
		" org1/project1/repo1 ",
		"  org2/project2/repo3  ",
	})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "org1/project1/repo1"},
		&MockVCS{repoPath: "org1/project1/repo2"},
		&MockVCS{repoPath: "org2/project2/repo3"},
	}

	result, _ := filterReposByConfigList(depotList)

	if len(result) != 2 {
		t.Errorf("应该正确处理带空格的配置，期望 2 个仓库，实际 %d 个", len(result))
	}
}

// TestFilterReposByConfigList_EmptyStrings 测试包含空字符串的配置
func TestFilterReposByConfigList_EmptyStrings(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{""})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "org1/project1/repo1"},
		&MockVCS{repoPath: "org1/project1/repo2"},
	}

	result, _ := filterReposByConfigList(depotList)

	if len(result) != 2 {
		t.Errorf("只包含空字符串的配置应返回完整列表，期望 2 个仓库，实际 %d 个", len(result))
	}
}

// TestFilterReposByConfigList_NoMatch 测试没有匹配的仓库
func TestFilterReposByConfigList_NoMatch(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{"org99/project99/repo99"})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "org1/project1/repo1"},
		&MockVCS{repoPath: "org1/project1/repo2"},
	}

	result, _ := filterReposByConfigList(depotList)

	if len(result) != 0 {
		t.Errorf("没有匹配的仓库应返回空列表，期望 0 个仓库，实际 %d 个", len(result))
	}
}

// TestFilterReposByConfigList_GitlabFormat 测试GitLab格式的仓库路径
func TestFilterReposByConfigList_GitlabFormat(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{
		"group1/subgroup1/repo1",
		"group2/repo2",
	})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "group1/subgroup1/repo1"},
		&MockVCS{repoPath: "group1/subgroup1/repo2"},
		&MockVCS{repoPath: "group2/repo2"},
		&MockVCS{repoPath: "group3/repo3"},
	}

	result, _ := filterReposByConfigList(depotList)

	if len(result) != 2 {
		t.Errorf("应该支持GitLab格式的仓库路径，期望 2 个仓库，实际 %d 个", len(result))
	}
}

// TestFilterReposByConfigList_GithubFormat 测试GitHub格式的仓库路径
func TestFilterReposByConfigList_GithubFormat(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{
		"owner1/repo1",
		"owner2/repo2",
	})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "owner1/repo1"},
		&MockVCS{repoPath: "owner1/repo2"},
		&MockVCS{repoPath: "owner2/repo2"},
	}

	result, _ := filterReposByConfigList(depotList)

	if len(result) != 2 {
		t.Errorf("应该支持GitHub格式的仓库路径，期望 2 个仓库，实际 %d 个", len(result))
	}
}

// TestFilterReposByConfigList_GongfengFormat 测试工蜂格式的仓库路径
func TestFilterReposByConfigList_GongfengFormat(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{
		"tencent/team1/project1/repo1",
		"tencent/team2/repo2",
	})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "tencent/team1/project1/repo1"},
		&MockVCS{repoPath: "tencent/team1/project1/repo2"},
		&MockVCS{repoPath: "tencent/team2/repo2"},
	}

	result, _ := filterReposByConfigList(depotList)

	if len(result) != 2 {
		t.Errorf("应该支持工蜂格式的仓库路径，期望 2 个仓库，实际 %d 个", len(result))
	}
}

// TestFilterReposByConfigList_MixedFormats 测试混合格式的仓库路径
func TestFilterReposByConfigList_MixedFormats(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{
		"owner/repo",            // GitHub格式
		"group/subgroup/repo",   // GitLab格式
		"org/team/project/repo", // 工蜂格式
	})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "owner/repo"},
		&MockVCS{repoPath: "group/subgroup/repo"},
		&MockVCS{repoPath: "org/team/project/repo"},
		&MockVCS{repoPath: "other/repo"},
	}

	result, _ := filterReposByConfigList(depotList)

	if len(result) != 3 {
		t.Errorf("应该支持混合格式的仓库路径，期望 3 个仓库，实际 %d 个", len(result))
	}
}

// TestFilterReposByConfigList_ExactMatch 测试精确匹配
func TestFilterReposByConfigList_ExactMatch(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{"org/project/repo"})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "org/project/repo"},
		&MockVCS{repoPath: "org/project/repo2"}, // 不应匹配
		&MockVCS{repoPath: "org/project2/repo"}, // 不应匹配
		&MockVCS{repoPath: "org2/project/repo"}, // 不应匹配
	}

	result, _ := filterReposByConfigList(depotList)

	if len(result) != 1 {
		t.Errorf("应该只精确匹配指定的仓库，期望 1 个仓库，实际 %d 个", len(result))
	}
	if result[0].GetRepoPath() != "org/project/repo" {
		t.Errorf("仓库路径不匹配，期望 org/project/repo，实际 %s", result[0].GetRepoPath())
	}
}

// TestFilterReposByConfigList_CaseSensitive 测试大小写敏感
func TestFilterReposByConfigList_CaseSensitive(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{"Org/Project/Repo"})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "Org/Project/Repo"},
		&MockVCS{repoPath: "org/project/repo"}, // 大小写不同，不应匹配
	}

	result, _ := filterReposByConfigList(depotList)

	if len(result) != 1 {
		t.Errorf("应该区分大小写，期望 1 个仓库，实际 %d 个", len(result))
	}
	if result[0].GetRepoPath() != "Org/Project/Repo" {
		t.Errorf("仓库路径不匹配，期望 Org/Project/Repo，实际 %s", result[0].GetRepoPath())
	}
}

// TestFilterReposByConfigList_EmptyDepotList 测试空的仓库列表
func TestFilterReposByConfigList_EmptyDepotList(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{"org/project/repo"})

	depotList := []vcs.VCS{}

	result, _ := filterReposByConfigList(depotList)

	if len(result) != 0 {
		t.Errorf("空的仓库列表应返回空列表，期望 0 个仓库，实际 %d 个", len(result))
	}
}

// TestFilterReposByConfigList_DuplicateConfig 测试配置中有重复的仓库
func TestFilterReposByConfigList_DuplicateConfig(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{
		"org/project/repo1",
		"org/project/repo1", // 重复
		"org/project/repo2",
	})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "org/project/repo1"},
		&MockVCS{repoPath: "org/project/repo2"},
	}

	result, _ := filterReposByConfigList(depotList)

	if len(result) != 2 {
		t.Errorf("重复配置应该正确处理，期望 2 个仓库，实际 %d 个", len(result))
	}
}

// TestFilterReposByConfigList_NotFound 测试未找到仓库的计数
func TestFilterReposByConfigList_NotFound(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	// 配置2个仓库，但只有1个存在
	config.Cfg = setupTestConfig([]string{
		"org/project/repo1",
		"org/project/repo-not-exist",
	})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "org/project/repo1"},
		&MockVCS{repoPath: "org/project/repo2"},
	}

	result, notFoundCount := filterReposByConfigList(depotList)

	if len(result) != 1 {
		t.Errorf("应该返回1个匹配的仓库，实际 %d 个", len(result))
	}
	if notFoundCount != 1 {
		t.Errorf("应该有1个未找到的仓库，实际 %d 个", notFoundCount)
	}
}

// TestFilterReposByConfigList_AllNotFound 测试所有仓库都未找到
func TestFilterReposByConfigList_AllNotFound(t *testing.T) {
	oldCfg := config.Cfg
	defer func() { config.Cfg = oldCfg }()

	config.Cfg = setupTestConfig([]string{
		"org/project/repo-not-exist1",
		"org/project/repo-not-exist2",
	})

	depotList := []vcs.VCS{
		&MockVCS{repoPath: "org/project/repo1"},
		&MockVCS{repoPath: "org/project/repo2"},
	}

	result, notFoundCount := filterReposByConfigList(depotList)

	if len(result) != 0 {
		t.Errorf("应该返回0个匹配的仓库，实际 %d 个", len(result))
	}
	if notFoundCount != 2 {
		t.Errorf("应该有2个未找到的仓库，实际 %d 个", notFoundCount)
	}
}
