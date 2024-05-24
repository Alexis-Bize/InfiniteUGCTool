package halowaypoint

import "time"

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

type MatchSpectateResponse struct {
	FilmStatusBond								int `json:"FilmStatusBond"`
	CustomData struct {
		FilmLength 								int `json:"FilmLength"`
		Chunks []struct {
			Index 								int `json:"Index"`
			ChunkStartTimeOffsetMilliseconds	int `json:"ChunkStartTimeOffsetMilliseconds"`
			DurationMilliseconds				int `json:"DurationMilliseconds"`
			ChunkSize							int `json:"ChunkSize"`
			FileRelativePath					string `json:"FileRelativePath"`
			ChunkType							int `json:"ChunkType"`
		} `json:"Chunks"`
		HasGameEnded							bool `json:"HasGameEnded"`
		ManifestRefreshSeconds					int `json:"ManifestRefreshSeconds"`
		MatchID									string `json:"MatchId"`
		FilmMajorVersion						int `json:"FilmMajorVersion"`
	} `json:"CustomData"`
	BlobStoragePathPrefix						string `json:"BlobStoragePathPrefix"`
	AssetID										string `json:"AssetId"`
}

// Partial
type MatchStatsResponse struct {
	MatchID	string `json:"MatchId"`
	MatchInfo struct {
		StartTime						time.Time `json:"StartTime"`
		EndTime							time.Time `json:"EndTime"`
		Duration						string `json:"Duration"`
		LifecycleMode					int `json:"LifecycleMode"`
		GameVariantCategory				int `json:"GameVariantCategory"`
		LevelID							string `json:"LevelId"`
		MapVariant struct {
			AssetKind					int `json:"AssetKind"`
			AssetID						string `json:"AssetId"`
			VersionID					string `json:"VersionId"`
		} `json:"MapVariant"`
		UgcGameVariant struct {
			AssetKind					int `json:"AssetKind"`
			AssetID						string `json:"AssetId"`
			VersionID					string `json:"VersionId"`
		} `json:"UgcGameVariant"`
		ClearanceID 					string `json:"ClearanceId"`
		Playlist struct {
			AssetKind					int `json:"AssetKind"`
			AssetID						string `json:"AssetId"`
			VersionID					string `json:"VersionId"`
		} `json:"Playlist"`
		PlaylistExperience 				int `json:"PlaylistExperience"`
		PlaylistMapModePair struct {
			AssetKind					int `json:"AssetKind"`
			AssetID						string `json:"AssetId"`
			VersionID					string `json:"VersionId"`
		} `json:"PlaylistMapModePair"`
		SeasonID						string `json:"SeasonId"`
		PlayableDuration				string `json:"PlayableDuration"`
		TeamsEnabled					bool `json:"TeamsEnabled"`
		TeamScoringEnabled				bool `json:"TeamScoringEnabled"`
		GameplayInteraction				int `json:"GameplayInteraction"`
	} `json:"MatchInfo"`
}

type NewSessionResponse struct {
	SessionID	string `json:"SessionId"`
	AssetID		string `json:"AssetId"`
}
