package domain

type League struct {
	QueueType string  `json:"queue_type"`
	Tier      string  `json:"tier"` //"MASTER"
	Rank      string  `json:"rank"` //"I"
	Wins      int     `json:"wins"`
	Losses    int     `json:"losses"`
	WinRate   float32 `json:"win_rate"`
}

type Summoner struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Level   int      `json:"level"`
	TeamId  int64    `json:"team_id"`
	Leagues []League `json:"leagues"`
}

type Match struct {
	Summoners []Summoner `json:"summoners"`
}

func NewLeague(queueType string, tier string, rank string, wins int, losses int) League {
	return League{
		QueueType: queueType,
		Tier:      tier,
		Rank:      rank,
		Wins:      wins,
		Losses:    losses,
		WinRate:   float32(wins) / float32(wins+losses),
	}
}

func NewSummoner(id string, name string, level int, teamId int64, leagues []League) Summoner {
	return Summoner{
		Id:      id,
		Name:    name,
		Level:   level,
		TeamId:  teamId,
		Leagues: leagues,
	}
}
