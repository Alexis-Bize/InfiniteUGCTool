package request

import (
	"fmt"
	"infinite-bookmarker/internal"
)

var RequestUserAgent = fmt.Sprintf(
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 (via %s/%s)",
	internal.GetConfig().Title,
	internal.GetConfig().Version,
)