package halowaypoint

type Urls struct {
	Profile		string
	Settings 	string
	Authoring 	string
	Discovery	string
	Stats		string
}

type Config struct {
	Urls 	Urls
	Title 	string
}

func GetConfig() *Config {
	return &Config{
		Urls: Urls{
			Profile: "https://profile.svc.halowaypoint.com",
			Settings: "https://settings.svc.halowaypoint.com",
			Authoring: "https://authoring-infiniteugc.svc.halowaypoint.com",
			Discovery: "https://discovery-infiniteugc.svc.halowaypoint.com",
			Stats: "https://halostats.svc.halowaypoint.com",
		},
		Title: "hi",
	}
}
