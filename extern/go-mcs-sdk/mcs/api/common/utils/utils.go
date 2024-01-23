package utils

import (
	"strings"
)

func UrlJoin(root string, parts ...string) string {
	url := root

	for _, part := range parts {
		part = strings.Trim(part, " ")
		if part == "" {
			continue
		}

		url = strings.TrimRight(url, "/") + "/" + strings.TrimLeft(part, "/")
	}

	url = strings.TrimRight(url, "/")

	return url
}

func IsStrEmpty(str *string) bool {
	return str == nil || *str == "" || strings.Trim(*str, " ") == ""
}
