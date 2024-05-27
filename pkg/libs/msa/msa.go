package msa

type LiveClientAuthOptions struct {
	ClientID		string
	Scope			string
	ResponseType	string
	RedirectURI		string
	State			string
}

type LivePreAuthResponse struct {
	Cookie	string
	Matches	LivePreAuthMatchedParameters
}

type LivePreAuthMatchedParameters struct {
	PPFT	string
	URLPost	string
}

type LiveCredentials struct {
	Email		string
	Password	string
}

type AuthStrategy struct {
	Data			string 	`json:"data"`
	Type			int		`json:"type"`
	Display			string 	`json:"display"`
	OtcEnabled		bool	`json:"otcEnabled"`
	OtcSent			bool	`json:"otcSent"`
	IsLost			bool	`json:"isLost"`
	IsSleeping		bool	`json:"isSleeping"`
	IsSADef			bool	`json:"isSADef"`
	IsVoiceDef		bool	`json:"isVoiceDef"`
	IsVoiceOnly		bool	`json:"isVoiceOnly"`
	PushEnabled		bool	`json:"pushEnabled"`
}