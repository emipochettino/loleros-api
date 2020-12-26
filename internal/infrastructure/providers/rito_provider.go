package providers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/emipochettino/loleros-api/internal/application"
	providers "github.com/emipochettino/loleros-api/internal/infrastructure/providers/dtos"
	"net/http"
)

type ritoProvider struct {
	client http.Client
	token  string
	host   map[string]string
	cache  Cache
}

func (r ritoProvider) FindSummonerByRegionAndName(region string, name string) (*providers.SummonerDTO, error) {
	if cached, isCached := r.cache.Get(fmt.Sprintf("summoner_by_name_%s_%s", region, name)); isCached {
		return cached.(*providers.SummonerDTO), nil
	}
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/lol/summoner/v4/summoners/by-name/%s", r.host[region], name), nil,
	)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Riot-Token", r.token)

	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		return nil, fmt.Errorf("rito token can be expired")
	}

	if response.StatusCode == 404 {
		return nil, fmt.Errorf("summoner not found")
	}

	var summonerDTO providers.SummonerDTO
	err = json.NewDecoder(response.Body).Decode(&summonerDTO)
	if err != nil {
		return nil, err
	}

	r.cache.SetDefault(fmt.Sprintf("summoner_by_name_%s_%s", region, name), &summonerDTO)

	return &summonerDTO, nil
}

func (r ritoProvider) FindSummonerByRegionAndId(region string, id string) (*providers.SummonerDTO, error) {
	if cached, isCached := r.cache.Get(fmt.Sprintf("summoner_by_id_%s_%s", region, id)); isCached {
		return cached.(*providers.SummonerDTO), nil
	}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/lol/summoner/v4/summoners/%s", r.host[region], id), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Riot-Token", r.token)

	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		return nil, fmt.Errorf("rito token can be expired")
	}

	if response.StatusCode == 404 {
		return nil, fmt.Errorf("summoner not found")
	}

	var summonerDTO providers.SummonerDTO
	err = json.NewDecoder(response.Body).Decode(&summonerDTO)
	if err != nil {
		return nil, err
	}

	r.cache.SetDefault(fmt.Sprintf("summoner_by_id_%s_%s", region, id), &summonerDTO)

	return &summonerDTO, nil
}

func (r ritoProvider) FindLeaguesByRegionAndSummonerId(region string, summonerId string) ([]providers.LeagueInfoDTO, error) {
	if cached, isCached := r.cache.Get(fmt.Sprintf("league_by_summoner_id_%s_%s", region, summonerId)); isCached {
		return cached.([]providers.LeagueInfoDTO), nil
	}
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/lol/league/v4/entries/by-summoner/%s", r.host[region], summonerId), nil)

	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Riot-Token", r.token)

	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		return nil, fmt.Errorf("rito token can be expired")
	}
	if response.StatusCode == 404 {
		return nil, fmt.Errorf("leagues not found")
	}

	var leagues []providers.LeagueInfoDTO
	err = json.NewDecoder(response.Body).Decode(&leagues)
	if err != nil {
		return nil, err
	}

	r.cache.SetDefault(fmt.Sprintf("league_by_summoner_id_%s_%s", region, summonerId), leagues)

	return leagues, nil
}

func (r ritoProvider) FindMatchBySummonerId(region string, summonerId string) (*providers.MatchDTO, error) {
	if cached, isCached := r.cache.Get(fmt.Sprintf("match_by_summoner_id_%s_%s", region, summonerId)); isCached {
		return cached.(*providers.MatchDTO), nil
	}
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/lol/spectator/v4/active-games/by-summoner/%s", r.host[region], summonerId), nil)

	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Riot-Token", r.token)

	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		return nil, fmt.Errorf("rito token can be expired")
	}
	if response.StatusCode == 404 {
		return nil, fmt.Errorf("match not found")
	}
	if response.StatusCode >= 500 {
		return nil, fmt.Errorf("something went wrong trying to find the active match")
	}

	var matchDTO providers.MatchDTO
	err = json.NewDecoder(response.Body).Decode(&matchDTO)
	if err != nil {
		return nil, err
	}

	r.cache.SetDefault(fmt.Sprintf("match_by_summoner_id_%s_%s", region, summonerId), &matchDTO)

	return &matchDTO, nil
}

func NewRitoProvider(host map[string]string, token string, cache Cache) (application.RitoProvider, error) {
	if len(token) == 0 {
		return nil, fmt.Errorf("rito token should exist")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//TODO receive this by parameter
	//c := cache.New(30*time.Minute, 40*time.Minute)

	return ritoProvider{
		client: http.Client{Transport: tr},
		token:  token,
		host:   host,
		cache:  cache,
	}, nil
}

type Cache interface {
	SetDefault(k string, x interface{})
	Get(k string) (interface{}, bool)
}
