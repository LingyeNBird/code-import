package vcs

import (
	"ccrctl/pkg/api/coding"
	"ccrctl/pkg/config"
	"testing"
)

func TestCodingVcs_GetCloneUrl(t *testing.T) {
	// 模拟测试数据
	testRepo := coding.Depots{
		HttpsUrl:    "https://coding.example.com/project/repo.git",
		SshUrl:      "git@coding.example.com:project/repo.git",
		Name:        "repo",
		ProjectName: "project",
	}

	codingVcs := &CodingVcs{
		httpURL:  testRepo.HttpsUrl,
		sshURL:   testRepo.SshUrl,
		RepoPath: "project/repo",
		RepoName: "repo",
	}

	tests := []struct {
		name     string
		sshMode  bool
		expected string
	}{
		{
			name:     "HTTP模式",
			sshMode:  false,
			expected: "https://coding:test-token@coding.example.com/project/repo.git",
		},
		{
			name:     "SSH模式",
			sshMode:  true,
			expected: "git@coding.example.com:project/repo.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置配置
			config.Cfg.Set("migrate.ssh", tt.sshMode)
			config.Cfg.Set("source.token", "test-token")

			result := codingVcs.GetCloneUrl()

			if tt.sshMode {
				// SSH 模式应该直接返回 SSH URL
				if result != tt.expected {
					t.Errorf("SSH模式下 GetCloneUrl() = %v, 期望 %v", result, tt.expected)
				}
			} else {
				// HTTP 模式应该包含认证信息
				if !contains(result, "coding") || !contains(result, "test-token") {
					t.Errorf("HTTP模式下 GetCloneUrl() = %v, 应该包含用户名和token", result)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
