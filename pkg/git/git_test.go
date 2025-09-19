package git

import (
	"testing"
)

func TestMaskSensitiveInfo(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "带凭证的URL",
			input:    "https://username:password@example.com/repo.git",
			expected: "https://****:****@example.com/repo.git",
		},
		{
			name:     "带凭证的URL和其他文本",
			input:    "推送到 https://user:secret123@git.example.org/path/repo.git 失败: 权限错误",
			expected: "推送到 https://****:****@git.example.org/path/repo.git 失败: 权限错误",
		},
		{
			name:     "多个带凭证的URL",
			input:    "尝试 https://user1:pass1@example.com 和 http://user2:pass2@example.org",
			expected: "尝试 https://****:****@example.com 和 http://****:****@example.org",
		},
		{
			name:     "不包含凭证的URL",
			input:    "https://example.com/repo.git",
			expected: "https://example.com/repo.git",
		},
		{
			name:     "包含特殊字符的凭证",
			input:    "https://user-name:p@ssw0rd!@example.com/repo.git",
			expected: "https://****:****@example.com/repo.git",
		},
		{
			name:     "实际案例",
			input:    "Locking support detected on remote \"https://cnb:feKWP2KtCzhbuvRg16164FHmDsC@cnb.cool/liujiboy/wander3d/Wander3d.git\".",
			expected: "Locking support detected on remote \"https://****:****@cnb.cool/liujiboy/wander3d/Wander3d.git\".",
		},
		{
			name:     "无URL的文本",
			input:    "这是一段没有URL的文本",
			expected: "这是一段没有URL的文本",
		},
		{
			name: "多行带凭证的URL",
			input: "第一行: https://user1:pass1@example.com/repo1.git\n" +
				"第二行: http://admin:secret@example.org/repo2.git\n" +
				"第三行: https://test:complex@p@ss@example.net/repo3.git",
			expected: "第一行: https://****:****@example.com/repo1.git\n" +
				"第二行: http://****:****@example.org/repo2.git\n" +
				"第三行: https://****:****@example.net/repo3.git",
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskSensitiveInfo(tt.input)
			if result != tt.expected {
				t.Errorf("maskSensitiveInfo() = %q, 期望 %q", result, tt.expected)
			}
		})
	}
}
