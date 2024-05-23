package halowaypoint

type UserProfileGamerpic struct {
	Small	string `json:"small"`
	Medium	string `json:"medium"`
	Large	string `json:"large"`
	Xlarge	string `json:"xlarge"`
}
type UserProfileResponse struct {
	Xuid		string `json:"xuid"`
	Gamertag 	string `json:"gamertag"`
	Gamerpic 	UserProfileGamerpic `json:"gamerpic"`
}

type SpartanToken struct {
	Value 		string `json:"value"`
	Expiration 	string `json:"expiration"`
}
