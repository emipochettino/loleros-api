package infrastructure

// SummonerDTO dto to map answer from rito api
type SummonerDTO struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"summonerLevel"`
}
