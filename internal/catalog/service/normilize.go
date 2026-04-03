package service

import "strings"

func normalizeCode(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)
	return v
}
