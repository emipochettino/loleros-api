package infrastructure

import "github.com/gin-gonic/gin"

func Route() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", Ping)
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", Ping)
	}

	return router
}
