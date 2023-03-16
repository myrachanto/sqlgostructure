package support

import "strings"

func Joins(name string) string {
	if name == "" {
		return ""
	}
	s := strings.TrimSpace(name)
	url := strings.Split(s, " ")
	return strings.Join(url, "-")
}
