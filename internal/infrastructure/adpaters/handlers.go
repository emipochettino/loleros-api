package infrastructure

import (
	"github.com/emipochettino/loleros-api/internal/application"
	"github.com/gin-gonic/gin"
	"net/http"
)

//create the handler with the needed dependencies.
type RitoHandler struct {
	MatchService application.MatchService
}

func (handler RitoHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Msg: "pong",
	})
}

func (handler RitoHandler) FindMatchInfoByRegionAndSummoner(c *gin.Context) {
	region, exists := c.GetQuery("region")
	if !exists {
		c.JSON(http.StatusBadRequest, Response{
			Msg: "The parameter region is required",
		})
		return
	}
	summonerName, exists := c.GetQuery("summoner_name")
	if !exists {
		c.JSON(http.StatusBadRequest, Response{
			Msg: "The parameter summoner_name is required",
		})
		return
	}

	match, err := handler.MatchService.FindCurrentMatchByRegionAndSummonerName(region, summonerName)
	if err != nil {
		//TODO improve the error handling
		c.JSON(http.StatusInternalServerError, Response{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, match)
}
