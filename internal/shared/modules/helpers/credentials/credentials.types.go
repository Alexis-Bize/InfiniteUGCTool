package credentials

type UserCredentials struct {
	Email 			string `json:"email,omitempty"`
	Password 		string `json:"password,omitempty"`
}

type SpartanTokenCredentials struct {
	Value			string `json:"value,omitempty"`
	Expiration		string `json:"expiration,omitempty"`
}

type Credentials struct {
	User			UserCredentials `json:"user,omitempty"`
	SpartanToken	SpartanTokenCredentials  `json:"spartan_token,omitempty"`
}