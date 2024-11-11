package util

import (
	"runtime"
	"strings"
)

func GetFunctionName() string {
	pc := make([]uintptr, 1) // at least 1 entry needed
	runtime.Callers(2, pc)   // 2 skips runtime.Callers and printFunctionName frames
	fn := runtime.FuncForPC(pc[0])
	return fn.Name()
}

// ConvertUrlWithAuth 把仓库httpURL转换为带认证的URL
func ConvertUrlWithAuth(url, username, password string) string {
	parts := strings.Split(url, "://")
	URL := parts[0] + "://" + username + ":" + password + "@" + parts[1]
	return URL
}
