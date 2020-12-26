package application
/*
import (
	"fmt"
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
	FindCurrentMatchByRegionAndSummonerName(region string, summonerName string) ([]application.SummonerDTO, error)
}

func (m matchService) FindCurrentMatchByRegionAndSummonerName(region string, summonerName string) ([]application.SummonerDTO, error) {
	start := time.Now()
	summonerDTO, err := m.ritoProvider.FindSummonerByRegionAndName(region, summonerName)
	if err != nil {
		return nil, err
	}

	matchDTO, err := m.ritoProvider.FindMatchBySummonerId(region, summonerDTO.Id)
	if err != nil {
		return nil, err
	}

	//TODO make this for async
	var wg sync.WaitGroup
	wg.Add(10)
	var summoners []application.SummonerDTO
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
			summoner := providers.MapToSummonerModel(*summonerDTO, participant, leaguesDTO)
			m.mu.Lock()
			summoners = append(summoners, summoner)
			m.mu.Unlock()
		}(participant)
	}
	wg.Wait()
	print(fmt.Sprintf("\ntime: %.2f seconds\n", time.Now().Sub(start).Seconds()))
	return summoners, nil
}

func NewMatchService(provider RitoProvider) MatchService {
	return matchService{ritoProvider: provider, mu: &sync.Mutex{}}
}

 */