package providers

import (
	"fmt"
	infrastructure "github.com/emipochettino/loleros-api/internal/infrastructure/providers/dtos"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRitoProvider(t *testing.T) {
	t.Run("Test create rito provider successfully", func(t *testing.T) {
		result, err := NewRitoProvider(map[string]string{"": ""}, "valid_token", cacheMock{})
		assert.NotNil(t, result)
		assert.Nil(t, err)
	})
}

func TestNewRitoProviderWithInvalidToken(t *testing.T) {
	t.Run("Test create rito provider with invalid token should return an error", func(t *testing.T) {
		result, err := NewRitoProvider(map[string]string{"": ""}, "", cacheMock{})
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestFindSummonerByRegionAndName(t *testing.T) {
	t.Run("Test find summoner by region and name successfully", func(t *testing.T) {
		content, err := ioutil.ReadFile("jsons/summoner_response.json")
		assert.Nil(t, err)
		server := serverMock(
			"/lol/summoner/v4/summoners/by-name/test_name",
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(content)
			})
		defer server.Close()

		provider, err := NewRitoProvider(
			map[string]string{"test_region": server.URL},
			"valid_token",
			createEmptyCache(),
		)
		assert.Nil(t, err)
		result, err := provider.FindSummonerByRegionAndName("test_region", "test_name")
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, &infrastructure.SummonerDTO{
			Id:    "flB50ZlPKdPOKSomx9Yep5FHrP-CGRdnkKHoH9nbhcLY_JxX",
			Name:  "xNibe",
			Level: 18,
		}, result)
	})
}

func TestFindSummonerByRegionAndNameWithErrors(t *testing.T) {
	tests := []struct {
		name          string
		pathFile      string
		statusCode    int
		expectedError error
	}{
		{
			"Test get summoner by region and name without access",
			"jsons/errors/forbidden_error.json",
			http.StatusForbidden,
			fmt.Errorf("rito token can be expired"),
		}, {
			"Test get summoner by region and name with non existing name",
			"jsons/errors/not_found_error.json",
			http.StatusNotFound,
			fmt.Errorf("summoner not found"),
		}, {
			"Test get summoner by region and name when rito api does not response correctly",
			"jsons/errors/internal_server_error.json",
			http.StatusInternalServerError,
			fmt.Errorf("uups, something went wrong"),
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.name), func(t *testing.T) {
			content, err := ioutil.ReadFile(tt.pathFile)
			assert.Nil(t, err)
			server := serverMock(
				"/lol/summoner/v4/summoners/by-name/test_name", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.statusCode)
					_, _ = w.Write(content)
				})
			defer server.Close()
			provider, err := NewRitoProvider(
				map[string]string{"test_region": server.URL},
				"valid_token",
				createEmptyCache(),
			)
			assert.Nil(t, err)
			_, err = provider.FindSummonerByRegionAndName("test_region", "test_name")
			assert.NotNil(t, err)
			assert.EqualValues(t, tt.expectedError, err)
		})
	}
}

func TestFindSummonerByRegionAndNameWithoutQuota(t *testing.T) {
	t.Run("Test get summoner by region and name without quota", func(t *testing.T) {
		contentTooManyRequests, err := ioutil.ReadFile("jsons/errors/too_many_requests_error.json")
		assert.Nil(t, err)
		contentSummonerResponse, err := ioutil.ReadFile("jsons/summoner_response.json")
		assert.Nil(t, err)
		requestNumber := 0
		server := serverMock(
			"/lol/summoner/v4/summoners/by-name/test_name", func(w http.ResponseWriter, r *http.Request) {
				if requestNumber < 1 {
					w.WriteHeader(http.StatusTooManyRequests)
					_, _ = w.Write(contentTooManyRequests)
				} else {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write(contentSummonerResponse)
				}
				requestNumber++
			})
		defer server.Close()
		provider, err := NewRitoProvider(
			map[string]string{"test_region": server.URL},
			"valid_token",
			createEmptyCache(),
		)
		assert.Nil(t, err)
		result, err := provider.FindSummonerByRegionAndName("test_region", "test_name")
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, &infrastructure.SummonerDTO{
			Id:    "flB50ZlPKdPOKSomx9Yep5FHrP-CGRdnkKHoH9nbhcLY_JxX",
			Name:  "xNibe",
			Level: 18,
		}, result)
	})

}

func TestFindSummonerByRegionAndId(t *testing.T) {
	t.Run("Test find summoner by region and id successfully", func(t *testing.T) {
		content, err := ioutil.ReadFile("jsons/summoner_response.json")
		assert.Nil(t, err)
		server := serverMock(
			"/lol/summoner/v4/summoners/test_id",
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(content)
			})
		defer server.Close()
		cacheMock := cacheMock{
			getMocked: func(k string) (interface{}, bool) {
				return nil, false
			},
			setDefaultMocked: func(k string, x interface{}) {
				//do nothing
			},
		}
		provider, err := NewRitoProvider(map[string]string{"test_region": server.URL}, "valid_token", cacheMock)
		assert.Nil(t, err)
		result, err := provider.FindSummonerByRegionAndId("test_region", "test_id")
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, &infrastructure.SummonerDTO{
			Id:    "flB50ZlPKdPOKSomx9Yep5FHrP-CGRdnkKHoH9nbhcLY_JxX",
			Name:  "xNibe",
			Level: 18,
		}, result)

	})
}

func TestFindSummonerByRegionAndIdWithErrors(t *testing.T) {
	tests := []struct {
		name          string
		pathFile      string
		statusCode    int
		expectedError error
	}{
		{
			"Test get summoner by region and name without access",
			"jsons/errors/forbidden_error.json",
			http.StatusForbidden,
			fmt.Errorf("rito token can be expired"),
		}, {
			"Test get summoner by region and name with non existing name",
			"jsons/errors/not_found_error.json",
			http.StatusNotFound,
			fmt.Errorf("summoner not found"),
		}, {
			"Test get summoner by region and name when rito api does not response correctly",
			"jsons/errors/internal_server_error.json",
			http.StatusInternalServerError,
			fmt.Errorf("uups, something went wrong"),
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.name), func(t *testing.T) {
			content, err := ioutil.ReadFile(tt.pathFile)
			assert.Nil(t, err)
			server := serverMock(
				"/lol/summoner/v4/summoners/test_id", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.statusCode)
					_, _ = w.Write(content)
				})
			defer server.Close()

			cacheMock := cacheMock{
				getMocked: func(k string) (interface{}, bool) {
					return nil, false
				},
				setDefaultMocked: func(k string, x interface{}) {
					//do nothing
				},
			}
			provider, err := NewRitoProvider(map[string]string{"test_region": server.URL}, "valid_token", cacheMock)
			assert.Nil(t, err)
			_, err = provider.FindSummonerByRegionAndId("test_region", "test_id")
			assert.NotNil(t, err)
			assert.EqualValues(t, tt.expectedError.Error(), err.Error())
		})
	}
}

func TestFindLeaguesByRegionAndSummonerId(t *testing.T) {
	t.Run("Test find leagues by region and summoner id successfully", func(t *testing.T) {
		content, err := ioutil.ReadFile("jsons/leagues_response.json")
		assert.Nil(t, err)
		server := serverMock(
			"/lol/league/v4/entries/by-summoner/test_id",
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(content)
			})
		defer server.Close()
		cacheMock := cacheMock{
			getMocked: func(k string) (interface{}, bool) {
				return nil, false
			},
			setDefaultMocked: func(k string, x interface{}) {
				//do nothing
			},
		}
		provider, err := NewRitoProvider(map[string]string{"test_region": server.URL}, "valid_token", cacheMock)
		assert.Nil(t, err)
		result, err := provider.FindLeaguesByRegionAndSummonerId("test_region", "test_id")
		assert.Nil(t, err)
		assert.NotNil(t, result)
		//todo assert values
	})
}

func TestFindLeaguesByRegionAndSummonerIdWithErrors(t *testing.T) {
	tests := []struct {
		name          string
		pathFile      string
		statusCode    int
		expectedError error
	}{
		{
			"Test find leagues by region and summoner id without access",
			"jsons/errors/forbidden_error.json",
			http.StatusForbidden,
			fmt.Errorf("rito token can be expired"),
		}, {
			"Test find leagues by region and summoner id with non existing summoner",
			"jsons/errors/not_found_error.json",
			http.StatusNotFound,
			fmt.Errorf("leagues not found"),
		}, {
			"Test find leagues by region and summoner id when rito api does not response correctly",
			"jsons/errors/internal_server_error.json",
			http.StatusInternalServerError,
			fmt.Errorf("uups, something went wrong"),
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.name), func(t *testing.T) {
			content, err := ioutil.ReadFile("jsons/summoner_response.json")
			assert.Nil(t, err)
			server := serverMock(
				"/lol/league/v4/entries/by-summoner/test_id",
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.statusCode)
					_, _ = w.Write(content)
				})
			defer server.Close()
			cacheMock := cacheMock{
				getMocked: func(k string) (interface{}, bool) {
					return nil, false
				},
				setDefaultMocked: func(k string, x interface{}) {
					//do nothing
				},
			}
			provider, err := NewRitoProvider(map[string]string{"test_region": server.URL}, "valid_token", cacheMock)
			assert.Nil(t, err)
			_, err = provider.FindLeaguesByRegionAndSummonerId("test_region", "test_id")
			assert.NotNil(t, err)
			assert.EqualValues(t, tt.expectedError, err)
		})
	}
}

func TestFindMatchBySummonerId(t *testing.T) {
	t.Run("Test find active game by region and summoner id successfully", func(t *testing.T) {
		content, err := ioutil.ReadFile("jsons/active_game_response.json")
		assert.Nil(t, err, "could not read file")
		server := serverMock(
			"/lol/spectator/v4/active-games/by-summoner/test_id",
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(content)
			})
		defer server.Close()
		cacheMock := cacheMock{
			getMocked: func(k string) (interface{}, bool) {
				return nil, false
			},
			setDefaultMocked: func(k string, x interface{}) {
				//do nothing
			},
		}
		provider, err := NewRitoProvider(map[string]string{"test_region": server.URL}, "valid_token", cacheMock)
		assert.Nil(t, err)
		result, err := provider.FindMatchBySummonerId("test_region", "test_id")
		assert.Nil(t, err)
		assert.NotNil(t, result)
		//todo assert values
	})
}

func TestFindMatchBySummonerIdWithErrors(t *testing.T) {
	tests := []struct {
		name          string
		pathFile      string
		statusCode    int
		expectedError error
	}{
		{
			"Test find match by region and summoner id without access",
			"jsons/errors/forbidden_error.json",
			http.StatusForbidden,
			fmt.Errorf("rito token can be expired"),
		}, {
			"Test find match by region and summoner id with non existing name",
			"jsons/errors/not_found_error.json",
			http.StatusNotFound,
			fmt.Errorf("match not found"),
		}, {
			"Test find match by region and summoner id with internal server error",
			"jsons/errors/internal_server_error.json",
			http.StatusInternalServerError,
			fmt.Errorf("uups, something went wrong"),
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.name), func(t *testing.T) {
			content, err := ioutil.ReadFile("jsons/summoner_response.json")
			assert.Nil(t, err, "could not read file")
			server := serverMock(
				"/lol/spectator/v4/active-games/by-summoner/test_id",
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.statusCode)
					_, _ = w.Write(content)
				})
			defer server.Close()
			cacheMock := cacheMock{
				getMocked: func(k string) (interface{}, bool) {
					return nil, false
				},
				setDefaultMocked: func(k string, x interface{}) {
					//do nothing
				},
			}
			provider, err := NewRitoProvider(map[string]string{"test_region": server.URL}, "valid_token", cacheMock)
			assert.Nil(t, err)
			_, err = provider.FindMatchBySummonerId("test_region", "test_id")
			assert.NotNil(t, err)
			assert.EqualValues(t, tt.expectedError, err)
		})
	}
}

func serverMock(path string, handlerFunc func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(path, handlerFunc)

	srv := httptest.NewServer(handler)

	return srv
}

func createEmptyCache() cacheMock {
	return cacheMock{
		getMocked: func(k string) (interface{}, bool) {
			return nil, false
		},
		setDefaultMocked: func(k string, x interface{}) {
			//do nothing
		},
	}
}

type cacheMock struct {
	setDefaultMocked func(k string, x interface{})
	getMocked        func(k string) (interface{}, bool)
}

func (c cacheMock) SetDefault(k string, x interface{}) {
	c.setDefaultMocked(k, x)
}

func (c cacheMock) Get(k string) (interface{}, bool) {
	return c.getMocked(k)
}
