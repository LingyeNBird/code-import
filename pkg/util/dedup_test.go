package util

import (
	"testing"
)

// Test_DeduplicateStringSlice_NoDuplicates 测试无重复配置
func Test_DeduplicateStringSlice_NoDuplicates(t *testing.T) {
	slice := []string{
		"org1/project1/repo1",
		"org1/project1/repo2",
		"org2/project2/repo3",
	}

	result := DeduplicateStringSlice(slice)

	if len(result.Deduplicated) != 3 {
		t.Errorf("无重复配置应返回完整列表，期望 3 项，实际 %d 项", len(result.Deduplicated))
	}
	if result.DuplicateCount != 0 {
		t.Errorf("无重复配置的 DuplicateCount 应为 0，实际 %d", result.DuplicateCount)
	}
}

// Test_DeduplicateStringSlice_WithDuplicates 测试包含重复配置
func Test_DeduplicateStringSlice_WithDuplicates(t *testing.T) {
	slice := []string{
		"org1/project1/repo1",
		"org1/project1/repo2",
		"org1/project1/repo1", // 重复
		"org2/project2/repo3",
		"org1/project1/repo2", // 重复
	}

	result := DeduplicateStringSlice(slice)

	if len(result.Deduplicated) != 3 {
		t.Errorf("去重后应返回 3 项，实际 %d 项", len(result.Deduplicated))
	}
	if result.DuplicateCount != 2 {
		t.Errorf("应有 2 个重复项，实际 %d", result.DuplicateCount)
	}

	resultSet := make(map[string]bool)
	for _, r := range result.Deduplicated {
		resultSet[r] = true
	}

	if !resultSet["org1/project1/repo1"] {
		t.Error("结果应包含 org1/project1/repo1")
	}
	if !resultSet["org1/project1/repo2"] {
		t.Error("结果应包含 org1/project1/repo2")
	}
	if !resultSet["org2/project2/repo3"] {
		t.Error("结果应包含 org2/project2/repo3")
	}
}

// Test_DeduplicateStringSlice_WithWhitespace 测试包含空格的配置
func Test_DeduplicateStringSlice_WithWhitespace(t *testing.T) {
	slice := []string{
		" org1/project1/repo1 ",
		"  org1/project1/repo1  ", // 去除空格后重复
		"org1/project1/repo2",
	}

	result := DeduplicateStringSlice(slice)

	if len(result.Deduplicated) != 2 {
		t.Errorf("去除空格后去重应返回 2 项，实际 %d 项", len(result.Deduplicated))
	}
	if result.DuplicateCount != 1 {
		t.Errorf("应有 1 个重复项，实际 %d", result.DuplicateCount)
	}
}

// Test_DeduplicateStringSlice_WithEmptyStrings 测试包含空字符串的配置
func Test_DeduplicateStringSlice_WithEmptyStrings(t *testing.T) {
	slice := []string{
		"org1/project1/repo1",
		"",
		"  ",
		"org1/project1/repo2",
		"",
	}

	result := DeduplicateStringSlice(slice)

	if len(result.Deduplicated) != 2 {
		t.Errorf("空字符串应被过滤，期望 2 项，实际 %d 项", len(result.Deduplicated))
	}
}

// Test_DeduplicateStringSlice_AllEmpty 测试全是空字符串的配置
func Test_DeduplicateStringSlice_AllEmpty(t *testing.T) {
	slice := []string{
		"",
		"  ",
		"   ",
	}

	result := DeduplicateStringSlice(slice)

	if len(result.Deduplicated) != 0 {
		t.Errorf("全是空字符串应返回空列表，期望 0 项，实际 %d 项", len(result.Deduplicated))
	}
}

// Test_DeduplicateStringSlice_EmptySlice 测试空切片
func Test_DeduplicateStringSlice_EmptySlice(t *testing.T) {
	slice := []string{}

	result := DeduplicateStringSlice(slice)

	if len(result.Deduplicated) != 0 {
		t.Errorf("空切片应返回空切片，期望 0 项，实际 %d 项", len(result.Deduplicated))
	}
}

// Test_DeduplicateStringSlice_PreserveOrder 测试保留首次出现顺序
func Test_DeduplicateStringSlice_PreserveOrder(t *testing.T) {
	slice := []string{
		"org1/project1/repo1",
		"org1/project1/repo2",
		"org1/project1/repo3",
		"org1/project1/repo2", // 重复
	}

	result := DeduplicateStringSlice(slice)

	if len(result.Deduplicated) != 3 {
		t.Errorf("去重后应返回 3 项，实际 %d 项", len(result.Deduplicated))
	}

	// 检查顺序是否保留
	if result.Deduplicated[0] != "org1/project1/repo1" {
		t.Errorf("第1项应为 org1/project1/repo1，实际 %s", result.Deduplicated[0])
	}
	if result.Deduplicated[1] != "org1/project1/repo2" {
		t.Errorf("第2项应为 org1/project1/repo2，实际 %s", result.Deduplicated[1])
	}
	if result.Deduplicated[2] != "org1/project1/repo3" {
		t.Errorf("第3项应为 org1/project1/repo3，实际 %s", result.Deduplicated[2])
	}
}

// Test_DeduplicateStringSlice_CaseSensitive 测试大小写敏感
func Test_DeduplicateStringSlice_CaseSensitive(t *testing.T) {
	slice := []string{
		"Org/Project/Repo",
		"org/project/repo", // 大小写不同，不应被视为重复
		"Org/Project/Repo", // 完全相同，应被视为重复
	}

	result := DeduplicateStringSlice(slice)

	if len(result.Deduplicated) != 2 {
		t.Errorf("大小写不同的值应被视为不同项，期望 2 项，实际 %d 项", len(result.Deduplicated))
	}
	if result.DuplicateCount != 1 {
		t.Errorf("应有 1 个重复项，实际 %d", result.DuplicateCount)
	}
}

// Test_DeduplicateStringSlice_Projects 测试项目配置去重
func Test_DeduplicateStringSlice_Projects(t *testing.T) {
	slice := []string{
		"project1",
		"project2",
		"project1", // 重复
		"project3",
	}

	result := DeduplicateStringSlice(slice)

	if len(result.Deduplicated) != 3 {
		t.Errorf("去重后应返回 3 个项目，实际 %d 个", len(result.Deduplicated))
	}
	if result.DuplicateCount != 1 {
		t.Errorf("应有 1 个重复项，实际 %d", result.DuplicateCount)
	}

	resultSet := make(map[string]bool)
	for _, r := range result.Deduplicated {
		resultSet[r] = true
	}

	if !resultSet["project1"] {
		t.Error("结果应包含 project1")
	}
	if !resultSet["project2"] {
		t.Error("结果应包含 project2")
	}
	if !resultSet["project3"] {
		t.Error("结果应包含 project3")
	}
}

// Test_DeduplicateStringSlice_MixedFormats 测试混合格式的仓库路径
func Test_DeduplicateStringSlice_MixedFormats(t *testing.T) {
	slice := []string{
		"owner/repo",            // GitHub格式
		"group/subgroup/repo",   // GitLab格式
		"org/team/project/repo", // 工蜂格式
		"owner/repo",            // 重复
		"group/subgroup/repo",   // 重复
	}

	result := DeduplicateStringSlice(slice)

	if len(result.Deduplicated) != 3 {
		t.Errorf("去重后应返回 3 项，实际 %d 项", len(result.Deduplicated))
	}
	if result.DuplicateCount != 2 {
		t.Errorf("应有 2 个重复项，实际 %d", result.DuplicateCount)
	}
}

// Test_DeduplicateStringSlice_SingleItem 测试单个配置项
func Test_DeduplicateStringSlice_SingleItem(t *testing.T) {
	slice := []string{"org/project/repo"}

	result := DeduplicateStringSlice(slice)

	if len(result.Deduplicated) != 1 {
		t.Errorf("单个配置项应返回 1 项，实际 %d 项", len(result.Deduplicated))
	}
	if result.Deduplicated[0] != "org/project/repo" {
		t.Errorf("配置项不匹配，期望 org/project/repo，实际 %s", result.Deduplicated[0])
	}
}

// Test_DeduplicateStringSlice_OnlyEmptyString 测试仅包含一个空字符串
func Test_DeduplicateStringSlice_OnlyEmptyString(t *testing.T) {
	slice := []string{""}

	result := DeduplicateStringSlice(slice)

	if len(result.Deduplicated) != 0 {
		t.Errorf("仅包含空字符串应返回空列表，期望 0 项，实际 %d 项", len(result.Deduplicated))
	}
}
