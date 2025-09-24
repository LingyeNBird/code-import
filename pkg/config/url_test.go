package config

import (
	"fmt"
	"regexp"
	"testing"
)

// checkURLForTest 是独立的URL检查函数，用于测试，不依赖全局配置
func checkURLForTest(url string) error {
	// 仅允许 http:// 或 https:// 开头，后跟域名或IP地址，支持端口号，不允许包含路径
	// 支持域名格式：example.com, api.example.com
	// 支持IP地址格式：192.168.1.1, 10.0.0.1

	// IP地址模式：严格匹配IPv4地址
	ipPattern := `^https?://((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(:[0-9]+)?$`

	// 域名模式：至少包含一个点的域名，且至少有一个部分包含字母
	domainPattern := `^https?://(([A-Za-z0-9-]*[A-Za-z][A-Za-z0-9-]*\.)+[A-Za-z0-9-]*[A-Za-z][A-Za-z0-9-]*|([A-Za-z0-9-]*[A-Za-z][A-Za-z0-9-]*\.)+[A-Za-z0-9-]+|([A-Za-z0-9-]+\.)+[A-Za-z0-9-]*[A-Za-z][A-Za-z0-9-]*)(:[0-9]+)?$`

	ipRegex := regexp.MustCompile(ipPattern)
	domainRegex := regexp.MustCompile(domainPattern)

	// 先检查IP地址模式，如果匹配则直接返回成功
	if ipRegex.MatchString(url) {
		return nil
	}

	// 再检查域名模式
	if domainRegex.MatchString(url) {
		return nil
	}

	return fmt.Errorf("url %s 格式错误，必须以 'http://' 或 'https://' 开头，且只能包含域名或IP地址，不能包含路径,如 https://e.coding.net、https://cnb.cool、https://example.com:8080、http://192.168.1.1:8080", url)
}

func TestCheckURLForTest(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		// 有效的URL测试用例
		{
			name:    "valid https URL without port",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "valid http URL without port",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "valid https URL with port",
			url:     "https://example.com:8080",
			wantErr: false,
		},
		{
			name:    "valid http URL with port",
			url:     "http://example.com:9000",
			wantErr: false,
		},
		{
			name:    "valid subdomain without port",
			url:     "https://a.example.com",
			wantErr: false,
		},
		{
			name:    "valid subdomain with port",
			url:     "https://a.example.com:9000",
			wantErr: false,
		},
		{
			name:    "valid multi-level subdomain",
			url:     "https://api.v1.example.com",
			wantErr: false,
		},
		{
			name:    "valid multi-level subdomain with port",
			url:     "https://api.v1.example.com:3000",
			wantErr: false,
		},
		{
			name:    "valid domain with hyphen",
			url:     "https://my-domain.com",
			wantErr: false,
		},
		{
			name:    "valid domain with hyphen and port",
			url:     "https://my-domain.com:8080",
			wantErr: false,
		},
		{
			name:    "valid coding example",
			url:     "https://e.coding.net",
			wantErr: false,
		},
		{
			name:    "valid cnb example",
			url:     "https://cnb.cool",
			wantErr: false,
		},

		// 无效的URL测试用例
		{
			name:    "invalid - no protocol",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "invalid - ftp protocol",
			url:     "ftp://example.com",
			wantErr: true,
		},
		{
			name:    "invalid - with path",
			url:     "https://example.com/path",
			wantErr: true,
		},
		{
			name:    "invalid - with query",
			url:     "https://example.com?query=1",
			wantErr: true,
		},
		{
			name:    "invalid - with fragment",
			url:     "https://example.com#fragment",
			wantErr: true,
		},
		{
			name:    "invalid - no domain",
			url:     "https://",
			wantErr: true,
		},
		{
			name:    "invalid - no TLD",
			url:     "https://example",
			wantErr: true,
		},
		{
			name:    "invalid - empty string",
			url:     "",
			wantErr: true,
		},
		{
			name:    "invalid - only protocol",
			url:     "https://",
			wantErr: true,
		},
		{
			name:    "invalid - invalid port format",
			url:     "https://example.com:abc",
			wantErr: true,
		},
		{
			name:    "invalid - port with slash",
			url:     "https://example.com:8080/",
			wantErr: true,
		},
		{
			name:    "invalid - multiple ports",
			url:     "https://example.com:8080:9000",
			wantErr: true,
		},
		{
			name:    "valid - port zero",
			url:     "https://example.com:0",
			wantErr: false, // 技术上端口0是有效的，虽然不常用
		},
		{
			name:    "valid - very high port",
			url:     "https://example.com:99999",
			wantErr: false, // 正则表达式不检查端口范围，只检查格式
		},
		{
			name:    "invalid - underscore in domain",
			url:     "https://example_domain.com",
			wantErr: true,
		},
		{
			name:    "invalid - space in URL",
			url:     "https://example .com",
			wantErr: true,
		},
		{
			name:    "invalid - trailing slash",
			url:     "https://example.com/",
			wantErr: true,
		},
		{
			name:    "valid - IP address",
			url:     "http://192.168.1.1",
			wantErr: false,
		},
		{
			name:    "valid - IP address with port",
			url:     "http://192.168.1.1:8080",
			wantErr: false,
		},
		{
			name:    "valid - IP address https",
			url:     "https://192.168.1.1",
			wantErr: false,
		},
		{
			name:    "valid - IP address https with port",
			url:     "https://192.168.1.1:8080",
			wantErr: false,
		},
		{
			name:    "valid - localhost IP",
			url:     "http://127.0.0.1",
			wantErr: false,
		},
		{
			name:    "valid - localhost IP with port",
			url:     "http://127.0.0.1:3000",
			wantErr: false,
		},
		{
			name:    "valid - private IP range 10.x",
			url:     "http://10.0.0.1",
			wantErr: false,
		},
		{
			name:    "valid - private IP range 10.x with port",
			url:     "http://10.0.0.1:9000",
			wantErr: false,
		},
		{
			name:    "valid - edge IP 0.0.0.0",
			url:     "http://0.0.0.0",
			wantErr: false,
		},
		{
			name:    "valid - edge IP 255.255.255.255",
			url:     "http://255.255.255.255",
			wantErr: false,
		},
		{
			name:    "invalid - IP with invalid octet 256",
			url:     "http://192.168.1.256",
			wantErr: true,
		},
		{
			name:    "invalid - IP with invalid octet 999",
			url:     "http://192.168.999.1",
			wantErr: true,
		},
		{
			name:    "invalid - incomplete IP",
			url:     "http://192.168.1",
			wantErr: true,
		},
		{
			name:    "invalid - IP with 5 octets",
			url:     "http://192.168.1.1.1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkURLForTest(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkURLForTest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckURLForTest_EdgeCases(t *testing.T) {
	// 测试边界情况
	edgeCases := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "single character domain parts",
			url:     "https://a.b.c",
			wantErr: false,
		},
		{
			name:    "long domain name",
			url:     "https://very-long-subdomain-name.very-long-domain-name.com",
			wantErr: false,
		},
		{
			name:    "domain with numbers",
			url:     "https://api2.example123.com",
			wantErr: false,
		},
		{
			name:    "domain with numbers and port",
			url:     "https://api2.example123.com:8080",
			wantErr: false,
		},
		{
			name:    "minimum valid port",
			url:     "https://example.com:1",
			wantErr: false,
		},
		{
			name:    "common ports",
			url:     "https://example.com:443",
			wantErr: false,
		},
		{
			name:    "http default port explicitly",
			url:     "http://example.com:80",
			wantErr: false,
		},
		{
			name:    "https default port explicitly",
			url:     "https://example.com:443",
			wantErr: false,
		},
	}

	for _, tt := range edgeCases {
		t.Run(tt.name, func(t *testing.T) {
			err := checkURLForTest(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkURLForTest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// 基准测试
func BenchmarkCheckURLForTest(b *testing.B) {
	testURL := "https://example.com:8080"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checkURLForTest(testURL)
	}
}

func BenchmarkCheckURLForTest_WithoutPort(b *testing.B) {
	testURL := "https://example.com"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checkURLForTest(testURL)
	}
}

func BenchmarkCheckURLForTest_ComplexDomain(b *testing.B) {
	testURL := "https://api.v1.subdomain.example.com:9000"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checkURLForTest(testURL)
	}
}
