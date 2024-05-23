package request

import "strings"

func GetBaseHeaders(extraHeaders map[string]string) map[string]string {
	headers := map[string]string{
		"User-Agent": RequestUserAgent,
		"Accept-Encoding": "identity",
	}

	for k, v := range extraHeaders {
		headers[k] = v
	}

	return headers
}

func ComputeUrl(baseUrl string, path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return baseUrl + path
}