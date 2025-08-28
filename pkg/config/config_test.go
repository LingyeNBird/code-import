package config

import (
	"strings"
	"testing"
)

func TestCheckCommonToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
		desc    string
	}{
		// 有效的 token 测试用例
		{
			name:    "有效的字母数字组合",
			token:   "abc123",
			wantErr: false,
			desc:    "纯字母数字组合应该通过",
		},
		{
			name:    "有效的下划线token",
			token:   "test_token",
			wantErr: false,
			desc:    "包含下划线的token应该通过",
		},
		{
			name:    "有效的中划线token",
			token:   "my-token",
			wantErr: false,
			desc:    "包含中划线的token应该通过",
		},
		{
			name:    "有效的点号token",
			token:   "token.v1",
			wantErr: false,
			desc:    "包含点号的token应该通过（新增功能）",
		},
		{
			name:    "有效的多点号token",
			token:   "api.token.123",
			wantErr: false,
			desc:    "包含多个点号的token应该通过",
		},
		{
			name:    "有效的混合字符token",
			token:   "test.api_token-v1",
			wantErr: false,
			desc:    "包含点号、下划线、中划线的混合token应该通过",
		},
		{
			name:    "有效的以点号结尾token",
			token:   "token.",
			wantErr: false,
			desc:    "以点号结尾的token应该通过",
		},
		{
			name:    "有效的以点号开头token",
			token:   ".token",
			wantErr: false,
			desc:    "以点号开头的token应该通过",
		},
		{
			name:    "有效的纯大写字母",
			token:   "ABCDEF",
			wantErr: false,
			desc:    "纯大写字母应该通过",
		},
		{
			name:    "有效的混合大小写",
			token:   "AbC123",
			wantErr: false,
			desc:    "混合大小写应该通过",
		},

		// 无效的 token 测试用例
		{
			name:    "无效的特殊字符@",
			token:   "token@123",
			wantErr: true,
			desc:    "包含@符号的token应该失败",
		},
		{
			name:    "无效的空格",
			token:   "token space",
			wantErr: true,
			desc:    "包含空格的token应该失败",
		},
		{
			name:    "无效的特殊字符#",
			token:   "token#123",
			wantErr: true,
			desc:    "包含#符号的token应该失败",
		},
		{
			name:    "无效的特殊字符$",
			token:   "token$123",
			wantErr: true,
			desc:    "包含$符号的token应该失败",
		},
		{
			name:    "无效的特殊字符%",
			token:   "token%123",
			wantErr: true,
			desc:    "包含%符号的token应该失败",
		},
		{
			name:    "无效的特殊字符&",
			token:   "token&123",
			wantErr: true,
			desc:    "包含&符号的token应该失败",
		},
		{
			name:    "无效的括号",
			token:   "token(123)",
			wantErr: true,
			desc:    "包含括号的token应该失败",
		},
		{
			name:    "无效的中文字符",
			token:   "token中文",
			wantErr: true,
			desc:    "包含中文字符的token应该失败",
		},
		{
			name:    "无效的特殊字符/",
			token:   "token/123",
			wantErr: true,
			desc:    "包含斜杠的token应该失败",
		},
		{
			name:    "无效的特殊字符\\",
			token:   "token\\123",
			wantErr: true,
			desc:    "包含反斜杠的token应该失败",
		},

		// 边界情况
		{
			name:    "空字符串",
			token:   "",
			wantErr: true,
			desc:    "空字符串应该失败",
		},
		{
			name:    "单个有效字符",
			token:   "a",
			wantErr: false,
			desc:    "单个有效字符应该通过",
		},
		{
			name:    "单个点号",
			token:   ".",
			wantErr: false,
			desc:    "单个点号应该通过",
		},
		{
			name:    "多个连续点号",
			token:   "...",
			wantErr: false,
			desc:    "多个连续点号应该通过",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkCommonToken(tt.token)

			if tt.wantErr && err == nil {
				t.Errorf("checkCommonToken() 应该返回错误，但没有返回。测试用例: %s, token: %q", tt.desc, tt.token)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("checkCommonToken() 不应该返回错误，但返回了: %v。测试用例: %s, token: %q", err, tt.desc, tt.token)
			}

			// 验证错误信息包含正确的提示
			if tt.wantErr && err != nil {
				expectedSubstring := "source.token 包含非法字符"
				if !strings.Contains(err.Error(), expectedSubstring) {
					t.Errorf("错误信息应该包含 %q，但实际错误信息是: %v", expectedSubstring, err)
				}
			}
		})
	}
}

// 注意：这里我们使用 strings.Contains 来检查错误信息

// 基准测试
func BenchmarkCheckCommonToken(b *testing.B) {
	testTokens := []string{
		"abc123",
		"test_token",
		"my-token",
		"token.v1",
		"api.token.123",
		"test.api_token-v1",
		"invalid@token",
		"invalid token",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, token := range testTokens {
			checkCommonToken(token)
		}
	}
}

// 测试正则表达式性能
func BenchmarkCheckCommonTokenSingleToken(b *testing.B) {
	token := "test.api_token-v1.123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checkCommonToken(token)
	}
}
