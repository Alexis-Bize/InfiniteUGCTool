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
