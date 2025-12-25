package migrate

import (
	"testing"
)

// Test_deduplicateWithLog 测试带日志的去重功能
func Test_deduplicateWithLog(t *testing.T) {
	// 测试有重复项的情况
	slice := []string{
		"org1/project1/repo1",
		"org1/project1/repo2",
		"org1/project1/repo1", // 重复
		"org2/project2/repo3",
		"org1/project1/repo2", // 重复
	}

	result := deduplicateWithLog(slice, "source.repo")

	if len(result) != 3 {
		t.Errorf("去重后应返回 3 项，实际 %d 项", len(result))
	}

	// 测试无重复项的情况（不应输出日志）
	slice2 := []string{
		"org1/project1/repo1",
		"org1/project1/repo2",
		"org2/project2/repo3",
	}

	result2 := deduplicateWithLog(slice2, "source.repo")

	if len(result2) != 3 {
		t.Errorf("无重复配置应返回完整列表，期望 3 项，实际 %d 项", len(result2))
	}
}
