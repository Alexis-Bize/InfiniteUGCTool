package request

import (
	"fmt"

	"infinite-ugc-tool/configs"
)

var RequestUserAgent = fmt.Sprintf(
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/124.0.2478.109 (via %s/%s)",
	configs.GetConfig().Name,
	configs.GetConfig().Version,
)
