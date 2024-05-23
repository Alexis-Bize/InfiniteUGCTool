package identity

type UserCredentials struct {
	Email		string `json:"email,omitempty"`
	Password	string `json:"password,omitempty"`
}

type SpartanTokenDetails struct {
	Value		string `json:"value,omitempty"`
	Expiration	string `json:"expiration,omitempty"`
}

type XboxNetworkIdentity struct {
	Xuid		string `json:"xuid,omitempty"`
	Gamertag	string `json:"gamertag,omitempty"`
}

type Identity struct {
	User			UserCredentials		`json:"user,omitempty"`
	SpartanToken	SpartanTokenDetails `json:"spartan_token,omitempty"`
	XboxNetwork		XboxNetworkIdentity	`json:"xbox_network,omitempty"`
}