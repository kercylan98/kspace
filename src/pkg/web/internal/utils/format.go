package utils

import "strings"

// FormatUrlPathCharacter 为缺失 URL 地址前缀的 path 填充 "/" 并返回
func FormatUrlPathCharacter(path string) string {
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}
