package providers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/emipochettino/loleros-api/internal/application"
	providers "github.com/emipochettino/loleros-api/internal/infrastructure/providers/dtos"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const rateLimitExceededErrorMsg = "rate limit exceeded"

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
	url := fmt.Sprintf("%s/lol/summoner/v4/summoners/by-name/%s", r.host[region], name)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Riot-Token", r.token)

	var summonerDTO providers.SummonerDTO
	err = r.retryRequestIfLimitExceeded(func() error {
		response, err := r.client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusForbidden {
			return fmt.Errorf("rito token can be expired")
		}
		if response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("summoner not found")
		}
		if response.StatusCode == http.StatusTooManyRequests {
			return fmt.Errorf(rateLimitExceededErrorMsg)
		}
		if response.StatusCode != http.StatusOK {
			return r.handleNotOkResponse(response, url)
		}

		if err = json.NewDecoder(response.Body).Decode(&summonerDTO); err != nil {
			return err
		}

		r.cache.SetDefault(fmt.Sprintf("summoner_by_name_%s_%s", region, name), &summonerDTO)

		return nil
	})

	return &summonerDTO, err
}

func (r ritoProvider) FindSummonerByRegionAndId(region string, id string) (*providers.SummonerDTO, error) {
	if cached, isCached := r.cache.Get(fmt.Sprintf("summoner_by_id_%s_%s", region, id)); isCached {
		return cached.(*providers.SummonerDTO), nil
	}
	url := fmt.Sprintf("%s/lol/summoner/v4/summoners/%s", r.host[region], id)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Riot-Token", r.token)

	var summonerDTO providers.SummonerDTO

	err = r.retryRequestIfLimitExceeded(func() error {
		response, err := r.client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusForbidden {
			return fmt.Errorf("rito token can be expired")
		}
		if response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("summoner not found")
		}
		if response.StatusCode == http.StatusTooManyRequests {
			return fmt.Errorf(rateLimitExceededErrorMsg)
		}
		if response.StatusCode != http.StatusOK {
			return r.handleNotOkResponse(response, url)
		}

		if err = json.NewDecoder(response.Body).Decode(&summonerDTO); err != nil {
			return err
		}

		r.cache.SetDefault(fmt.Sprintf("summoner_by_id_%s_%s", region, id), &summonerDTO)

		return nil
	})

	return &summonerDTO, err
}

func (r ritoProvider) FindLeaguesByRegionAndSummonerId(region string, summonerId string) ([]providers.LeagueInfoDTO, error) {
	if cached, isCached := r.cache.Get(fmt.Sprintf("league_by_summoner_id_%s_%s", region, summonerId)); isCached {
		return cached.([]providers.LeagueInfoDTO), nil
	}
	url := fmt.Sprintf("%s/lol/league/v4/entries/by-summoner/%s", r.host[region], summonerId)
	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Riot-Token", r.token)

	var leagues []providers.LeagueInfoDTO
	err = r.retryRequestIfLimitExceeded(func() error {
		response, err := r.client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusForbidden {
			return fmt.Errorf("rito token can be expired")
		}
		if response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("leagues not found")
		}
		if response.StatusCode == http.StatusTooManyRequests {
			return fmt.Errorf(rateLimitExceededErrorMsg)
		}
		if response.StatusCode != http.StatusOK {
			return r.handleNotOkResponse(response, url)
		}

		if err = json.NewDecoder(response.Body).Decode(&leagues); err != nil {
			return err
		}

		r.cache.SetDefault(fmt.Sprintf("league_by_summoner_id_%s_%s", region, summonerId), leagues)

		return nil
	})

	return leagues, err
}

func (r ritoProvider) FindMatchBySummonerId(region string, summonerId string) (*providers.MatchDTO, error) {
	if cached, isCached := r.cache.Get(fmt.Sprintf("match_by_summoner_id_%s_%s", region, summonerId)); isCached {
		return cached.(*providers.MatchDTO), nil
	}
	url := fmt.Sprintf("%s/lol/spectator/v4/active-games/by-summoner/%s", r.host[region], summonerId)
	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Riot-Token", r.token)

	var matchDTO providers.MatchDTO
	err = r.retryRequestIfLimitExceeded(func() error {
		response, err := r.client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode == 403 {
			return fmt.Errorf("rito token can be expired")
		}
		if response.StatusCode == 404 {
			return fmt.Errorf("match not found")
		}
		if response.StatusCode == http.StatusTooManyRequests {
			return fmt.Errorf(rateLimitExceededErrorMsg)
		}
		if response.StatusCode != http.StatusOK {
			return r.handleNotOkResponse(response, url)
		}

		err = json.NewDecoder(response.Body).Decode(&matchDTO)
		if err != nil {
			return err
		}

		r.cache.SetDefault(fmt.Sprintf("match_by_summoner_id_%s_%s", region, summonerId), &matchDTO)

		return nil
	})

	return &matchDTO, err
}

func (r ritoProvider) retryRequestIfLimitExceeded(requestFunction func() error) error {
	return retry.Do(
		requestFunction,
		retry.RetryIf(isLimitExceeded),
		retry.Attempts(5),
		retry.Delay(1*time.Second),
		retry.LastErrorOnly(true),
	)
}

func (r ritoProvider) handleNotOkResponse(response *http.Response, url string) error {
	errorMsg, _ := ioutil.ReadAll(response.Body)
	log.Printf("error when requested %s. Status %d - Response %s\n",
		url,
		response.StatusCode,
		string(errorMsg),
	)
	return fmt.Errorf("uups, something went wrong")
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

func isLimitExceeded(err error) bool {
	return err.Error() == rateLimitExceededErrorMsg
}

type Cache interface {
	SetDefault(k string, x interface{})
	Get(k string) (interface{}, bool)
}
