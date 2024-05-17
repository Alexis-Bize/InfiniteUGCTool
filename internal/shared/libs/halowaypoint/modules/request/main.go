package halowaypoint_request

import "Infinite-Bookmarker/internal/shared"

func GetBaseHeaders(extraHeaders map[string]string) map[string]string {
	headers := map[string]string{
		"User-Agent": shared.RequestUserAgent,
		"Accept-Encoding": "identity",
	}

	for k, v := range extraHeaders {
		headers[k] = v
	}

	return headers
}