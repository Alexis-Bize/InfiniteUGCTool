package halowaypointRequest

type userProfileGamerpic struct {
	Small	string `json:"small"`
	Medium	string `json:"medium"`
	Large	string `json:"large"`
	Xlarge	string `json:"xlarge"`
} 
type UserProfileResponse struct {
	Xuid		string `json:"xuid"`
	Gamertag 	string `json:"gamertag"`
	Gamerpic 	userProfileGamerpic `json:"gamerpic"`
}