package fweb

import "strings"

// parsePattern 处理pattern中的‘*’
// 例如："/a/b/*/c" -> ["a","b","*"]
func parsePattern(pattern string) []string {
	pieces := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range pieces {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}
