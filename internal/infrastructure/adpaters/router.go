package infrastructure

import "github.com/gin-gonic/gin"

func NewRouter(ritoHandler RitoHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", ritoHandler.Ping)
		v1.GET("/rito/match", ritoHandler.FindMatchInfoByRegionAndSummoner)
	}

	return router
}
