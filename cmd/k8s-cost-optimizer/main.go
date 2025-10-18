package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shivamchaubey027/k8s-cost-optimizer/internal/models"
	"github.com/shivamchaubey027/k8s-cost-optimizer/pkg/database"
	"github.com/shivamchaubey027/k8s-cost-optimizer/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Server struct {
	pods []models.Pod
}

func (s *Server) getPods(c *gin.Context) {

	var pods []models.Pod
	result := database.DB.Find(&pods)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "No pods found"})
	}

	c.JSON(200, pods)
}

func (s *Server) createPod(c *gin.Context) {
	var newPod models.Pod

	if err := c.ShouldBindJSON(&newPod); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	result := database.DB.Create(&newPod)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to save Pod"})
		return
	}

	c.JSON(201, newPod)
}

func (s *Server) putPod(c *gin.Context) {
	idToUpdate := c.Param("id")

	var existingPod models.Pod
	findResult := database.DB.Find(&models.Pod{}, idToUpdate)
	if findResult.Error != nil {
		c.JSON(404, gin.H{"error": "error"})
		return
	}

	var updatePod models.Pod
	if err := c.ShouldBindJSON(&updatePod); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	existingPod.Name = updatePod.Name
	existingPod.Namespace = updatePod.Namespace
	existingPod.CPURequest = updatePod.CPURequest
	existingPod.MemoryRequest = updatePod.MemoryRequest

	saveResult := database.DB.Save(&existingPod)
	if saveResult.Error != nil {
		c.JSON(500, gin.H{"error": "Couldnt Save Results"})
	}

	c.JSON(200, existingPod)

}

func (s *Server) deletePod(c *gin.Context) {
	idToDelete := c.Param("id")

	result := database.DB.Delete(&models.Pod{}, idToDelete)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to delete"})
		return
	}

	c.JSON(200, gin.H{"message": "Pod deleted successfully"})

}

func (s *Server) getLivePods(c *gin.Context) {
	log.Println("--- DEBUG: Calling getLivePods ---")
	podLists, err := k8s.Clientset.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {

		log.Printf("!!! DEBUG: Error fetching pods: %v", err)
		c.JSON(500, gin.H{"error": "Bad Request"})
		return
	}
	log.Printf("--- DEBUG: Found %d pods in 'default' namespace ---", len(podLists.Items))
	c.JSON(200, podLists.Items)
}

func main() {

	database.Connect()
	k8s.Connect()

	server := &Server{
		pods: []models.Pod{},
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

	v1.DELETE("/pods/:id", server.deletePod)

	v1.PUT("/pods/:id", server.putPod)

	v1.GET("/k8s/pods/live", server.getLivePods)

	router.Run("localhost:8080")
}
