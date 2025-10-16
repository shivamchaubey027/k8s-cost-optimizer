package main

import "github.com/gin-gonic/gin"

type Pod struct {
	Name          string `json:"name"`
	Namespace     string `json:"nameSpace"`
	CPURequest    string `json:"cpuRequest"`
	MemoryRequest string `json:"memoryRequest"`
	IpAddress     string `json:"ipAddress"`
}

type Server struct {
	pods []Pod
}

func (s *Server) getPods(c *gin.Context) {
	c.JSON(200, s.pods)
}

func (s *Server) createPod(c *gin.Context) {
	var newPod Pod

	if err := c.ShouldBindJSON(&newPod); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	s.pods = append(s.pods, newPod)

	c.JSON(201, newPod)
}

func main() {
	server := &Server{
		pods: []Pod{
			{Name: "webapp-1", Namespace: "production", CPURequest: "500m", MemoryRequest: "1Gi"},
			{Name: "database-1", Namespace: "production", CPURequest: "4", MemoryRequest: "8Gi"},
			{Name: "cache-1", Namespace: "production", CPURequest: "1", MemoryRequest: "2Gi"},
			{Name: "api-gateway-1", Namespace: "ingress", CPURequest: "1", MemoryRequest: "1Gi"},
		},
	}

	router := gin.Default()

	health := router.Group("/health")

	v1 := router.Group("/api/v1")

	v1.GET("/pods", server.getPods)

	v1.POST("/pods", server.createPod)

	health.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.Run("localhost:8080")
}
