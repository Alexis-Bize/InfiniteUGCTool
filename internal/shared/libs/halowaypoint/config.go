package halowaypoint

type urls struct {
	Profile		string
	Settings 	string
	Authoring 	string
}

type Config struct {
	Urls 	urls
	Title 	string
}

func GetConfig() *Config {
	return &Config{
		Urls: urls{
			Profile: "https://profile.svc.halowaypoint.com",
			Settings: "https://settings.svc.halowaypoint.com",
			Authoring: "https://authoring-infiniteugc.svc.halowaypoint.com",
		},
		Title: "hi",
	}
}
