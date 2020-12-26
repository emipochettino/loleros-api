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
		content, err := ioutil.ReadFile("json/summoner_response.json")
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
			"json/forbbiden_error.json",
			http.StatusForbidden,
			fmt.Errorf("rito token can be expired"),
		}, {
			"Test get summoner by region and name with non existing name",
			"json/not_found_error.json",
			http.StatusNotFound,
			fmt.Errorf("summoner not found"),
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

//TODO CUANDO SE QUEDA SIN CUOTA
func TestFindSummonerByRegionAndNameWithout(t *testing.T) {

}

func TestFindSummonerByRegionAndId(t *testing.T) {
	t.Run("Test find summoner by region and id successfully", func(t *testing.T) {
		server := serverMock("/lol/summoner/v4/summoners/test_id", summonerMock)
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
		mockFunc      func(w http.ResponseWriter, r *http.Request)
		expectedError error
	}{
		{
			"Test find summoner by region and id without access",
			forbiddenMock,
			fmt.Errorf("rito token can be expired"),
		}, {
			"Test find summoner by region and id with non existing name",
			notFoundMock,
			fmt.Errorf("summoner not found"),
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.name), func(t *testing.T) {
			server := serverMock("/lol/summoner/v4/summoners/test_id", tt.mockFunc)
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
			assert.EqualValues(t, tt.expectedError, err)
		})
	}
}

func TestFindLeaguesByRegionAndSummonerId(t *testing.T) {
	t.Run("Test find leagues by region and summoner id successfully", func(t *testing.T) {
		server := serverMock("/lol/league/v4/entries/by-summoner/test_id", leaguesMock)
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
		mockFunc      func(w http.ResponseWriter, r *http.Request)
		expectedError error
	}{
		{
			"Test find leagues by region and summoner id without access",
			forbiddenMock,
			fmt.Errorf("rito token can be expired"),
		}, {
			"Test find leagues by region and summoner id with non existing summoner",
			notFoundMock,
			fmt.Errorf("leagues not found"),
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.name), func(t *testing.T) {
			server := serverMock("/lol/league/v4/entries/by-summoner/test_id", tt.mockFunc)
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
		server := serverMock("/lol/spectator/v4/active-games/by-summoner/test_id", activeGameMock)
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
		mockFunc      func(w http.ResponseWriter, r *http.Request)
		expectedError error
	}{
		{
			"Test find match by region and summoner id without access",
			forbiddenMock,
			fmt.Errorf("rito token can be expired"),
		}, {
			"Test find match by region and summoner id with non existing name",
			notFoundMock,
			fmt.Errorf("match not found"),
		}, {
			"Test find match by region and summoner id with internal server error",
			internalServerErrorMock,
			fmt.Errorf("something went wrong trying to find the active match"),
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.name), func(t *testing.T) {
			server := serverMock("/lol/spectator/v4/active-games/by-summoner/test_id", tt.mockFunc)
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

func summonerMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(getSummonerJson))
}

func leaguesMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(getLeaguesJson))
}

func activeGameMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(getActiveGameJson))
}

func forbiddenMock(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(forbiddenError))
}

func notFoundMock(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(notFoundError))
}

func internalServerErrorMock(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(internalServerError))
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

const (
	getSummonerJson = `{
							"id": "flB50ZlPKdPOKSomx9Yep5FHrP-CGRdnkKHoH9nbhcLY_JxX",
							"accountId": "QqLuDVPvar_I4xjQgnGZUeDhSLwvutmQ0fFOsu6j06xRIES5e2azgG2Y",
							"puuid": "vaRYQXkiSj5D8mJUw8dFYdsktMLRYLABaF4SEnA9C69PX_yuJjsM9C1kdiCjvvVHnxpU9dUWvjqM8A",
							"name": "xNibe",
							"profileIconId": 3542,
							"revisionDate": 1585079006000,
							"summonerLevel": 18
						}`
	getLeaguesJson = `[
						  {
							"leagueId": "b8c67709-6857-484e-a75d-aae615c7b776",
							"queueType": "RANKED_FLEX_SR",
							"tier": "GOLD",
							"rank": "I",
							"summonerId": "5Nn5hoqMZjWtssygP7bJl0fnCZneGrO90_TSS02olMXG9gM",
							"summonerName": "Prodigium",
							"leaguePoints": 100,
							"wins": 34,
							"losses": 21,
							"veteran": false,
							"inactive": false,
							"freshBlood": false,
							"hotStreak": false,
							"miniSeries": {
							  "target": 3,
							  "wins": 1,
							  "losses": 1,
							  "progress": "LWNNN"
							}
						  },
						  {
							"leagueId": "5b76481e-0327-453c-b6a1-1daf4dc6252d",
							"queueType": "RANKED_SOLO_5x5",
							"tier": "DIAMOND",
							"rank": "II",
							"summonerId": "5Nn5hoqMZjWtssygP7bJl0fnCZneGrO90_TSS02olMXG9gM",
							"summonerName": "Prodigium",
							"leaguePoints": 81,
							"wins": 523,
							"losses": 524,
							"veteran": false,
							"inactive": false,
							"freshBlood": false,
							"hotStreak": false
						  }
						]`
	getActiveGameJson = `{
							  "gameId": 3620211084,
							  "mapId": 11,
							  "gameMode": "CLASSIC",
							  "gameType": "MATCHED_GAME",
							  "gameQueueConfigId": 420,
							  "participants": [
								{
								  "teamId": 100,
								  "spell1Id": 11,
								  "spell2Id": 14,
								  "championId": 120,
								  "profileIconId": 3587,
								  "summonerName": "lIIlIIlIIlIIl",
								  "bot": false,
								  "summonerId": "NY4NS-f3ikxVVXrDg4Z6AdntugcQeBpSBSl9JTwy5WAXcY8",
								  "gameCustomizationObjects": [],
								  "perks": {
									"perkIds": [
									  8230,
									  8275,
									  8234,
									  8232,
									  8143,
									  8135,
									  5005,
									  5008,
									  5002
									],
									"perkStyle": 8200,
									"perkSubStyle": 8100
								  }
								},
								{
								  "teamId": 100,
								  "spell1Id": 4,
								  "spell2Id": 7,
								  "championId": 22,
								  "profileIconId": 4025,
								  "summonerName": "Boy Wonder",
								  "bot": false,
								  "summonerId": "iSCZESEpDckfhyV9-EtEmMPtIeJ7E0-QG3hcsvYovbQZdW4",
								  "gameCustomizationObjects": [],
								  "perks": {
									"perkIds": [
									  8008,
									  8009,
									  9104,
									  8014,
									  8304,
									  8410,
									  5005,
									  5008,
									  5002
									],
									"perkStyle": 8000,
									"perkSubStyle": 8300
								  }
								},
								{
								  "teamId": 100,
								  "spell1Id": 14,
								  "spell2Id": 4,
								  "championId": 39,
								  "profileIconId": 1665,
								  "summonerName": "tales of omero",
								  "bot": false,
								  "summonerId": "dTg0l3G0-B84dwZTq-cWzFCfDhT2MuU7ElHkyKw3sA9Sh1s",
								  "gameCustomizationObjects": [],
								  "perks": {
									"perkIds": [
									  8010,
									  9111,
									  9104,
									  8299,
									  8345,
									  8352,
									  5005,
									  5008,
									  5002
									],
									"perkStyle": 8000,
									"perkSubStyle": 8300
								  }
								},
								{
								  "teamId": 100,
								  "spell1Id": 14,
								  "spell2Id": 4,
								  "championId": 44,
								  "profileIconId": 4777,
								  "summonerName": "Mystic Polar",
								  "bot": false,
								  "summonerId": "fZb1PCKKa1KOjCrylJdGgqFnJo68A7gdh5PXM884QZRNEjo",
								  "gameCustomizationObjects": [],
								  "perks": {
									"perkIds": [
									  8465,
									  8463,
									  8473,
									  8453,
									  8009,
									  9105,
									  5005,
									  5003,
									  5002
									],
									"perkStyle": 8400,
									"perkSubStyle": 8000
								  }
								},
								{
								  "teamId": 100,
								  "spell1Id": 4,
								  "spell2Id": 12,
								  "championId": 8,
								  "profileIconId": 4777,
								  "summonerName": "without regrets",
								  "bot": false,
								  "summonerId": "Gk-SK7RjE8MvulMciYI6foykYhue81TSZClnyKsE3HxIcn0",
								  "gameCustomizationObjects": [],
								  "perks": {
									"perkIds": [
									  8230,
									  8275,
									  8210,
									  8236,
									  8304,
									  8347,
									  5007,
									  5002,
									  5002
									],
									"perkStyle": 8200,
									"perkSubStyle": 8300
								  }
								},
								{
								  "teamId": 200,
								  "spell1Id": 4,
								  "spell2Id": 12,
								  "championId": 134,
								  "profileIconId": 22,
								  "summonerName": "LainFangMu",
								  "bot": false,
								  "summonerId": "lsoW6vTSAHN52fHsd3XNbHLnYytL03N2yK88WgZmC4YiC9ua",
								  "gameCustomizationObjects": [],
								  "perks": {
									"perkIds": [
									  8230,
									  8226,
									  8210,
									  8236,
									  8304,
									  8345,
									  5005,
									  5008,
									  5002
									],
									"perkStyle": 8200,
									"perkSubStyle": 8300
								  }
								},
								{
								  "teamId": 200,
								  "spell1Id": 12,
								  "spell2Id": 4,
								  "championId": 236,
								  "profileIconId": 603,
								  "summonerName": "Katanivarina",
								  "bot": false,
								  "summonerId": "3RmZkYKMYkV6UdO9umImcvlSP1hFFaYgY7T78BOHSE0jmbM",
								  "gameCustomizationObjects": [],
								  "perks": {
									"perkIds": [
									  8005,
									  8009,
									  9104,
									  8299,
									  8304,
									  8345,
									  5005,
									  5008,
									  5001
									],
									"perkStyle": 8000,
									"perkSubStyle": 8300
								  }
								},
								{
								  "teamId": 200,
								  "spell1Id": 11,
								  "spell2Id": 4,
								  "championId": 254,
								  "profileIconId": 4676,
								  "summonerName": "Booky Boos",
								  "bot": false,
								  "summonerId": "13BihA3huqxFviwUfWRwfl3-qTEooxnpLJvhB-C8tg-TvzE",
								  "gameCustomizationObjects": [],
								  "perks": {
									"perkIds": [
									  8010,
									  9111,
									  9105,
									  8014,
									  8233,
									  8236,
									  5005,
									  5008,
									  5002
									],
									"perkStyle": 8000,
									"perkSubStyle": 8200
								  }
								},
								{
								  "teamId": 200,
								  "spell1Id": 4,
								  "spell2Id": 7,
								  "championId": 222,
								  "profileIconId": 4626,
								  "summonerName": "UF Dan",
								  "bot": false,
								  "summonerId": "l5mqKTnFwzu3prBf06icM-HmVbBuaNYnKNF_sL0tnmo_g-g",
								  "gameCustomizationObjects": [],
								  "perks": {
									"perkIds": [
									  8008,
									  8009,
									  9103,
									  8014,
									  8275,
									  8236,
									  5005,
									  5008,
									  5002
									],
									"perkStyle": 8000,
									"perkSubStyle": 8200
								  }
								},
								{
								  "teamId": 200,
								  "spell1Id": 4,
								  "spell2Id": 14,
								  "championId": 117,
								  "profileIconId": 4471,
								  "summonerName": "ˆºˆˆºˆˆºˆ",
								  "bot": false,
								  "summonerId": "FtVZJ62xCok2e1hEFQEfU1EdJ_PGHtWBUXmqG1HvzbDsjXQ",
								  "gameCustomizationObjects": [],
								  "perks": {
									"perkIds": [
									  8214,
									  8226,
									  8233,
									  8236,
									  8345,
									  8347,
									  5007,
									  5008,
									  5002
									],
									"perkStyle": 8200,
									"perkSubStyle": 8300
								  }
								}
							  ],
							  "observers": {
								"encryptionKey": "Bu9hDYLI00dqhmf39F8gcvyA2zXgrIvT"
							  },
							  "platformId": "NA1",
							  "bannedChampions": [
								{
								  "championId": 875,
								  "teamId": 100,
								  "pickTurn": 1
								},
								{
								  "championId": 142,
								  "teamId": 100,
								  "pickTurn": 2
								},
								{
								  "championId": 360,
								  "teamId": 100,
								  "pickTurn": 3
								},
								{
								  "championId": 412,
								  "teamId": 100,
								  "pickTurn": 4
								},
								{
								  "championId": 84,
								  "teamId": 100,
								  "pickTurn": 5
								},
								{
								  "championId": 555,
								  "teamId": 200,
								  "pickTurn": 6
								},
								{
								  "championId": 11,
								  "teamId": 200,
								  "pickTurn": 7
								},
								{
								  "championId": 107,
								  "teamId": 200,
								  "pickTurn": 8
								},
								{
								  "championId": 81,
								  "teamId": 200,
								  "pickTurn": 9
								},
								{
								  "championId": 104,
								  "teamId": 200,
								  "pickTurn": 10
								}
							  ],
							  "gameStartTime": 1602886365464,
							  "gameLength": 926
							}`
	forbiddenError      = `{"message": "forbidden"}`
	notFoundError       = `{"message": "not found"}`
	internalServerError = `{"message": "internal server error"}`
)
