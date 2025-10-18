package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shivamchaubey027/k8s-cost-optimizer/internal/handlers"
)

func SetupRouter(router *gin.Engine) {
	server := &handlers.Server{}

	health := router.Group("/health")
	{
		health.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	v1 := router.Group("/api/v1")
	{

		v1.GET("/pods", server.GetPods)

		v1.POST("/pods", server.CreatePod)

		v1.DELETE("/pods/:id", server.DeletePod)

		v1.PUT("/pods/:id", server.PutPod)

		v1.GET("/k8s/pods/live", server.GetLivePods)
	}

}
