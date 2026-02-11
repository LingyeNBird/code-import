package git

import (
	"testing"
)

func TestRemoveCredentialsFromURL(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "带凭证的URL",
			input:    "https://username:password@example.com/repo.git",
			expected: "https://example.com/repo.git",
		},
		{
			name:     "带凭证的URL和其他文本",
			input:    "推送到 https://user:secret123@git.example.org/path/repo.git 失败: 权限错误",
			expected: "推送到 https://git.example.org/path/repo.git 失败: 权限错误",
		},
		{
			name:     "多个带凭证的URL",
			input:    "尝试 https://user1:pass1@example.com 和 http://user2:pass2@example.org",
			expected: "尝试 https://example.com 和 http://example.org",
		},
		{
			name:     "不包含凭证的URL",
			input:    "https://example.com/repo.git",
			expected: "https://example.com/repo.git",
		},
		{
			name:     "包含特殊字符的凭证",
			input:    "https://user-name:p@ssw0rd!@example.com/repo.git",
			expected: "https://example.com/repo.git",
		},
		{
			name:     "实际案例",
			input:    "Locking support detected on remote \"https://cnb:feKWP2KtCzhbuvRg161ddsadasada@cnb.cool/aaa/wander3d/Wander3d.git\".",
			expected: "Locking support detected on remote \"https://cnb.cool/aaa/wander3d/Wander3d.git\".",
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
			expected: "第一行: https://example.com/repo1.git\n" +
				"第二行: http://example.org/repo2.git\n" +
				"第三行: https://example.net/repo3.git",
		},
		{
			name:     "只有用户名的URL",
			input:    "https://user@example.com/repo.git",
			expected: "https://example.com/repo.git",
		},
		{
			name:     "包含端口的URL",
			input:    "https://user:pass@example.com:8080/repo.git",
			expected: "https://example.com:8080/repo.git",
		},
		{
			name:     "包含查询参数的URL",
			input:    "https://user:pass@example.com/repo.git?branch=main&tag=v1.0",
			expected: "https://example.com/repo.git?branch=main&tag=v1.0",
		},
		{
			name:     "混合协议的URL",
			input:    "ssh://user@example.com/repo.git 和 https://user:pass@example.org/repo.git",
			expected: "ssh://user@example.com/repo.git 和 https://example.org/repo.git",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "无效URL格式",
			input:    "://example.com",
			expected: "://example.com",
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeCredentialsFromURL(tt.input)
			if result != tt.expected {
				t.Errorf("removeCredentialsFromURL() = %q, 期望 %q", result, tt.expected)
			}
		})
	}
}

func TestIsLFSObjectNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected bool
	}{
		{
			name:     "典型的 Repository or object not found 错误",
			output:   "Error downloading object: LFS: Repository or object not found",
			expected: true,
		},
		{
			name:     "小写的错误信息",
			output:   "error: repository or object not found",
			expected: true,
		},
		{
			name:     "混合大小写的错误信息",
			output:   "ERROR: Repository OR Object NOT FOUND",
			expected: true,
		},
		{
			name:     "多行输出包含错误",
			output:   "Fetching LFS objects...\nError: Repository or object not found\nfailed to download",
			expected: true,
		},
		{
			name:     "错误信息在中间位置",
			output:   "batch response: Repository or object not found: abc123.bin",
			expected: true,
		},
		{
			name:     "Object does not exist 错误（不应匹配）",
			output:   "batch response: Object does not exist on the server",
			expected: false,
		},
		{
			name:     "HTTP 404 错误（不应匹配）",
			output:   "error: RPC failed; HTTP 404 Not Found",
			expected: false,
		},
		{
			name:     "网络超时错误（不应匹配）",
			output:   "error: RPC failed; curl 28 Operation timed out",
			expected: false,
		},
		{
			name:     "权限错误（不应匹配）",
			output:   "error: Access denied. You do not have permission",
			expected: false,
		},
		{
			name:     "网络连接错误（不应匹配）",
			output:   "fatal: unable to access 'https://example.com': Could not resolve host",
			expected: false,
		},
		{
			name:     "成功信息（不应匹配）",
			output:   "Git LFS: (3 of 3 files) 15.5 MB / 15.5 MB",
			expected: false,
		},
		{
			name:     "空字符串",
			output:   "",
			expected: false,
		},
		{
			name:     "速率限制错误（不应匹配）",
			output:   "error: API rate limit exceeded",
			expected: false,
		},
		{
			name:     "部分匹配不应成功",
			output:   "Repository found but object missing",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsLFSObjectNotFoundError(tt.output)
			if result != tt.expected {
				t.Errorf("IsLFSObjectNotFoundError() = %v, 期望 %v\n输出: %q", result, tt.expected, tt.output)
			}
		})
	}
}

// 注意：FetchLFS 和 PushLFS 函数包含重试机制
// 重试配置：失败时自动重试 3 次，重试间隔为 2s、5s、10s
// 这可以有效应对以下临时故障场景：
//   - 网络抖动
//   - 服务端临时不可用
//   - API 速率限制
//   - 临时权限问题
//
// allowIncompletePush 参数的智能判断：
//   - 只有当错误确认是"源文件损坏/丢失"时才会忽略错误
//   - 对于网络、权限等可恢复的错误，即使开启该参数也会报错
//   - 这样可以最大化数据完整性，避免误用该参数导致正常文件未迁移
//
// 如需测试重试机制，建议使用集成测试或手动测试
// 单元测试中难以模拟真实的网络环境和重试场景
