package infrastructure

// SummonerDTO dto to map answer from rito api
type MatchDTO struct {
	GameId            int64            `json:"gameId"`
	GameType          string           `json:"gameType"`
	GameMode          string           `json:"gameMode"`
	GameQueueConfigId int64            `json:"gameQueueConfigId"`
	MapId             int64            `json:"mapId"`
	GameStartTime     int64            `json:"gameStartTime"`
	Participants      []ParticipantDTO `json:"participants"`
}

type ParticipantDTO struct {
	TeamId       int64  `json:"teamId"`
	SummonerName string `json:"summonerName"`
	SummonerId   string `json:"summonerId"`
}
