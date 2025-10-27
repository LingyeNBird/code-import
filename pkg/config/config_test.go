package config

import (
	"testing"

	"ccrctl/pkg/util"
	"github.com/spf13/viper"
)

func TestProcessCNBRootOrganization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "前后都有斜杠和空格",
			input:    " /coding/ ",
			expected: "coding",
		},
		{
			name:     "只有前面斜杠",
			input:    "/coding",
			expected: "coding",
		},
		{
			name:     "只有后面斜杠",
			input:    "coding/",
			expected: "coding",
		},
		{
			name:     "只有空格",
			input:    " coding ",
			expected: "coding",
		},
		{
			name:     "没有斜杠和空格",
			input:    "coding",
			expected: "coding",
		},
		{
			name:     "复杂路径",
			input:    " /my-org/sub-org/ ",
			expected: "my-org/sub-org",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 viper 实例
			v := viper.New()
			
			// 设置初始值
			v.Set("cnb.root_organization", tt.input)
			
			// 调用处理函数
			processCNBRootOrganization(v)
			
			// 验证结果
			result := v.GetString("cnb.root_organization")
			if result != tt.expected {
				t.Errorf("processCNBRootOrganization() 处理 %q = %q, 期望 %q", tt.input, result, tt.expected)
			}
		})
	}
}

// 测试空值情况
func TestProcessCNBRootOrganizationEmptyValue(t *testing.T) {
	v := viper.New()
	
	// 不设置任何值，默认为空字符串
	processCNBRootOrganization(v)
	
	result := v.GetString("cnb.root_organization")
	if result != "" {
		t.Errorf("processCNBRootOrganization() 处理空值应该保持为空，但得到 %q", result)
	}
}

// 测试与 util.TrimSlash 的一致性
func TestProcessCNBRootOrganizationConsistency(t *testing.T) {
	testCases := []string{
		" /coding/ ",
		"/coding",
		"coding/",
		" coding ",
		"  //coding//  ",
		" / ",
		"",
		"my-org",
		" /my-org/sub-org/ ",
	}

	for _, testCase := range testCases {
		t.Run("一致性测试_"+testCase, func(t *testing.T) {
			// 使用 viper 处理
			v := viper.New()
			v.Set("cnb.root_organization", testCase)
			processCNBRootOrganization(v)
			viperResult := v.GetString("cnb.root_organization")
			
			// 直接使用 util.TrimSlash 处理
			utilResult := util.TrimSlash(testCase)
			
			// 对于空字符串，processCNBRootOrganization 不会处理，所以结果应该是原值
			if testCase == "" {
				if viperResult != testCase {
					t.Errorf("空字符串处理不一致: viper=%q, 原值=%q", viperResult, testCase)
				}
			} else {
				// 非空字符串应该与 util.TrimSlash 结果一致
				if viperResult != utilResult {
					t.Errorf("处理结果不一致: viper=%q, util.TrimSlash=%q", viperResult, utilResult)
				}
			}
		})
	}
}