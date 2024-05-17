package internal

type Config struct {
	Title string
	Version string
}

func GetConfig() *Config {
	return &Config{
		Title: "Infinite Bookmarker",
		Version: "1.0.0",
	}
}