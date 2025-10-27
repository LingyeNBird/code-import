package util

import (
    "testing"
)

func TestTrimSlash(t *testing.T) {
    // 定义测试用例
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "前后都有斜杠和空格",
            input:    " /path/to/resource/ ",
            expected: "path/to/resource",
        },
        {
            name:     "只有前面斜杠和空格",
            input:    " /path/to/resource",
            expected: "path/to/resource",
        },
        {
            name:     "只有后面斜杠和空格",
            input:    "path/to/resource/ ",
            expected: "path/to/resource",
        },
        {
            name:     "前后都有斜杠",
            input:    "/path/to/resource/",
            expected: "path/to/resource",
        },
        {
            name:     "只有前面斜杠",
            input:    "/path/to/resource",
            expected: "path/to/resource",
        },
        {
            name:     "只有后面斜杠",
            input:    "path/to/resource/",
            expected: "path/to/resource",
        },
        {
            name:     "没有斜杠只有空格",
            input:    " path/to/resource ",
            expected: "path/to/resource",
        },
        {
            name:     "没有斜杠和空格",
            input:    "path/to/resource",
            expected: "path/to/resource",
        },
        {
            name:     "只有斜杠和空格",
            input:    " / ",
            expected: "",
        },
        {
            name:     "只有一个前斜杠",
            input:    "/",
            expected: "",
        },
        {
            name:     "多个前斜杠和空格",
            input:    "  //path/to/resource",
            expected: "/path/to/resource",
        },
        {
            name:     "多个后斜杠和空格",
            input:    "path/to/resource//  ",
            expected: "path/to/resource/",
        },
        {
            name:     "前后都有多个斜杠和空格",
            input:    "  //path/to/resource//  ",
            expected: "/path/to/resource/",
        },
        {
            name:     "空字符串",
            input:    "",
            expected: "",
        },
        {
            name:     "只有空格",
            input:    "   ",
            expected: "",
        },
        {
            name:     "单个字符前后有斜杠和空格",
            input:    " /a/ ",
            expected: "a",
        },
        {
            name:     "中间有斜杠和空格",
            input:    " /path/to/resource/ ",
            expected: "path/to/resource",
        },
        {
            name:     "根路径带空格",
            input:    " / ",
            expected: "",
        },
        {
            name:     "双斜杠带空格",
            input:    " // ",
            expected: "",
        },
        {
            name:     "URL路径带空格",
            input:    " /api/v1/users/ ",
            expected: "api/v1/users",
        },
        {
            name:     "文件路径带空格",
            input:    " /home/user/documents/ ",
            expected: "home/user/documents",
        },
        {
            name:     "相对路径带空格",
            input:    " ./path/ ",
            expected: "./path",
        },
        {
            name:     "包含特殊字符和空格",
            input:    " /path-with_special.chars/ ",
            expected: "path-with_special.chars",
        },
        {
            name:     "制表符和换行符",
            input:    "\t/path/to/resource/\n",
            expected: "path/to/resource",
        },
        {
            name:     "混合空白字符",
            input:    " \t\n/path/to/resource/\t \n",
            expected: "path/to/resource",
        },
    }

    // 执行测试用例
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := TrimSlash(tt.input)
            if result != tt.expected {
                t.Errorf("TrimSlash(%q) = %q, 期望 %q", tt.input, result, tt.expected)
            }
        })
    }
}

// 基准测试
func BenchmarkTrimSlash(b *testing.B) {
    testCases := []string{
        " /path/to/resource/ ",
        " /simple/ ",
        "no-slash",
        " / ",
        "  //multiple//slashes//  ",
        " /very/long/path/with/many/segments/that/might/be/used/in/real/applications/ ",
        "\t/path/with/tabs/\t",
        " \n/path/with/newlines/\n ",
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for _, tc := range testCases {
            TrimSlash(tc)
        }
    }
}

// 表格驱动的边界测试
func TestTrimSlashEdgeCases(t *testing.T) {
    edgeCases := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "非常长的字符串带空格",
            input:    " /" + string(make([]byte, 1000)) + "/ ",
            expected: string(make([]byte, 1000)),
        },
        {
            name:     "Unicode字符带空格",
            input:    " /测试路径/ ",
            expected: "测试路径",
        },
        {
            name:     "包含内部空格",
            input:    " / path with spaces / ",
            expected: " path with spaces ",
        },
        {
            name:     "特殊字符组合带空格",
            input:    " /!@#$%^&*()/ ",
            expected: "!@#$%^&*()",
        },
        {
            name:     "多种空白字符",
            input:    " \t\r\n/path/\t\r\n ",
            expected: "path",
        },
    }

    for _, tt := range edgeCases {
        t.Run(tt.name, func(t *testing.T) {
            result := TrimSlash(tt.input)
            if result != tt.expected {
                t.Errorf("TrimSlash(%q) = %q, 期望 %q", tt.input, result, tt.expected)
            }
        })
    }
}

// 测试函数的幂等性
func TestTrimSlashIdempotent(t *testing.T) {
    testInputs := []string{
        " /path/to/resource/ ",
        " /simple ",
        " simple/ ",
        "no-slash",
        " / ",
        "",
        "   ",
        "\t/path/\t",
    }

    for _, input := range testInputs {
        t.Run("幂等性测试_"+input, func(t *testing.T) {
            first := TrimSlash(input)
            second := TrimSlash(first)
            
            if first != second {
                t.Errorf("TrimSlash 不是幂等的: TrimSlash(%q) = %q, TrimSlash(TrimSlash(%q)) = %q", 
                    input, first, input, second)
            }
        })
    }
}

// 测试与原始实现的兼容性
func TestTrimSlashCompatibility(t *testing.T) {
    // 这些测试用例确保新实现与原始实现兼容
    compatibilityTests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "原始用例1",
            input:    "/path/to/resource/",
            expected: "path/to/resource",
        },
        {
            name:     "原始用例2",
            input:    "/path/to/resource",
            expected: "path/to/resource",
        },
        {
            name:     "原始用例3",
            input:    "path/to/resource/",
            expected: "path/to/resource",
        },
        {
            name:     "原始用例4",
            input:    "path/to/resource",
            expected: "path/to/resource",
        },
        {
            name:     "原始用例5",
            input:    "/",
            expected: "",
        },
    }

    for _, tt := range compatibilityTests {
        t.Run(tt.name, func(t *testing.T) {
            result := TrimSlash(tt.input)
            if result != tt.expected {
                t.Errorf("TrimSlash(%q) = %q, 期望 %q", tt.input, result, tt.expected)
            }
        })
    }
}