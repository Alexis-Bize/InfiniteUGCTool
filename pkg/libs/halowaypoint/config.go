// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
