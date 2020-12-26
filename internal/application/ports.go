package application

import providers "github.com/emipochettino/loleros-api/internal/infrastructure/providers/dtos"

type RitoProvider interface {
	FindSummonerByRegionAndName(region string, name string) (*providers.SummonerDTO, error)
	FindMatchBySummonerId(region string, summonerId string) (*providers.MatchDTO, error)
	FindSummonerByRegionAndId(region string, id string) (*providers.SummonerDTO, error)
	FindLeaguesByRegionAndSummonerId(region string, summonerId string) ([]providers.LeagueInfoDTO, error)
}
