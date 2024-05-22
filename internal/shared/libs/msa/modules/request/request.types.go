package msaRequest

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
