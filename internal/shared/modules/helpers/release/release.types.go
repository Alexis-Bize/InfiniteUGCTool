package release

type Release struct {
	TagName string `json:"tag_name"`
}

type Releases []Release
