package vcs

import (
	"testing"
)

// TestCodingVcs_GetSubGroup_RetryMechanism 测试 GetSubGroup 的重试机制
// 注意: 这是一个基础测试,验证方法不会 panic 并返回有效的 SubGroup
func TestCodingVcs_GetSubGroup_RetryMechanism(t *testing.T) {
	tests := []struct {
		name         string
		subGroupName string
		wantName     string
	}{
		{
			name:         "API调用失败时返回最小SubGroup",
			subGroupName: "test-project",
			wantName:     "test-project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codingVcs := &CodingVcs{
				SubGroupName: tt.subGroupName,
			}

			// 注意: 此测试在实际环境中会尝试调用真实 API 并可能失败
			// 这里我们主要验证:
			// 1. 方法不会 panic
			// 2. 返回的 SubGroup 至少包含 Name 字段
			// 完整的测试应该使用 mock 来模拟 API 调用

			// 由于没有有效的 API 配置,这里预期会失败并返回最小 SubGroup
			// 但不应该导致程序退出或 panic
			result := codingVcs.GetSubGroup()

			if result == nil {
				t.Fatal("GetSubGroup() 返回 nil, 预期返回非 nil 的 SubGroup")
			}

			if result.Name != tt.wantName {
				t.Errorf("GetSubGroup().Name = %v, 期望 %v", result.Name, tt.wantName)
			}

			// 验证失败时 Desc 和 Remark 为空(不嵌入错误信息)
			if result.Desc != "" && result.Desc != result.Name {
				t.Logf("注意: SubGroup.Desc = %v, 如果API成功则正常,否则应为空", result.Desc)
			}
			if result.Remark == "ERROR" {
				t.Errorf("GetSubGroup().Remark = 'ERROR', 不应该嵌入错误标记")
			}
		})
	}
}

// TestCodingVcs_GetSubGroup_NoExit 测试确保 API 失败不会导致程序退出
func TestCodingVcs_GetSubGroup_NoExit(t *testing.T) {
	codingVcs := &CodingVcs{
		SubGroupName: "non-existent-project",
	}

	// 这个调用在没有有效配置时应该失败,但不应该调用 os.Exit()
	// 我们通过 defer recover 来捕获可能的 panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetSubGroup() 发生 panic: %v", r)
		}
	}()

	result := codingVcs.GetSubGroup()

	// 验证即使失败也返回了有效的 SubGroup
	if result == nil {
		t.Fatal("API 失败时应该返回最小 SubGroup,而不是 nil")
	}

	if result.Name == "" {
		t.Error("API 失败时 SubGroup.Name 不应该为空")
	}
}
