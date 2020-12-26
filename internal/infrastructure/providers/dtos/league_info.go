package infrastructure

type LeagueInfoDTO struct {
	QueueType string `json:"queueType"`
	Tier      string `json:"tier"` //"MASTER"
	Rank      string `json:"rank"` //"I"
	Wins      int `json:"wins"`
	Losses    int `json:"losses"`
}
