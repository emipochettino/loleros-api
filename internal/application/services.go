package application

import (
	"fmt"
	"github.com/emipochettino/loleros-api/internal/domain"
	providers "github.com/emipochettino/loleros-api/internal/infrastructure/providers/dtos"
	"log"
	"sync"
	"time"
)

type matchService struct {
	ritoProvider RitoProvider
	mu           *sync.Mutex
}

type MatchService interface {
	FindCurrentMatchByRegionAndSummonerName(region string, summonerName string) (*domain.Match, error)
}

func (m matchService) FindCurrentMatchByRegionAndSummonerName(region string, summonerName string) (*domain.Match, error) {
	start := time.Now()
	summonerDTO, err := m.ritoProvider.FindSummonerByRegionAndName(region, summonerName)
	if err != nil {
		return nil, err
	}

	matchDTO, err := m.ritoProvider.FindMatchBySummonerId(region, summonerDTO.Id)
	if err != nil {
		return nil, err
	}

	//TODO make this async
	//new approach: create waitGroup, add 1 in each iteration before go func, results add into other chan, wait
	summoners := make(chan domain.Summoner, len(matchDTO.Participants))
	var wg sync.WaitGroup
	// add the number of summoners in the match
	wg.Add(len(matchDTO.Participants))
	for _, participant := range matchDTO.Participants {
		go func(participant providers.ParticipantDTO) {
			defer wg.Done()
			summonerDTO, err := m.ritoProvider.FindSummonerByRegionAndId(region, participant.SummonerId)
			if err != nil {
				log.Println(err)
				return
			}
			leaguesDTO, err := m.ritoProvider.FindLeaguesByRegionAndSummonerId(region, summonerDTO.Id)
			if err != nil {
				return
			}

			var leagues []domain.League
			for _, leagueDTO := range leaguesDTO {
				leagues = append(
					leagues,
					domain.NewLeague(
						leagueDTO.QueueType,
						leagueDTO.Tier,
						leagueDTO.Rank,
						leagueDTO.Wins,
						leagueDTO.Losses,
					))
			}

			summoner := domain.NewSummoner(
				summonerDTO.Id,
				summonerDTO.Name,
				summonerDTO.Level,
				participant.TeamId,
				leagues,
			)
			summoners <- summoner
		}(participant)
	}
	wg.Wait()
	close(summoners)
	summonerSlice := make([]domain.Summoner, 0)
	for summoner := range summoners {
		summonerSlice = append(summonerSlice, summoner)
	}

	print(fmt.Sprintf("\ntime: %.2f seconds\n", time.Now().Sub(start).Seconds()))
	return &domain.Match{Summoners: summonerSlice}, nil
}

func NewMatchService(provider RitoProvider) MatchService {
	return matchService{ritoProvider: provider, mu: &sync.Mutex{}}
}
