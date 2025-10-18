package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/shivamchaubey027/k8s-cost-optimizer/internal/models"
	"github.com/shivamchaubey027/k8s-cost-optimizer/pkg/database"
	"github.com/shivamchaubey027/k8s-cost-optimizer/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Server struct {
	pods []models.Pod
}

func (s *Server) GetPods(c *gin.Context) {

	var pods []models.Pod
	result := database.DB.Find(&pods)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "No pods found"})
	}

	c.JSON(200, pods)
}

func (s *Server) CreatePod(c *gin.Context) {
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

func (s *Server) PutPod(c *gin.Context) {
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

func (s *Server) DeletePod(c *gin.Context) {
	idToDelete := c.Param("id")

	result := database.DB.Delete(&models.Pod{}, idToDelete)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to delete"})
		return
	}

	c.JSON(200, gin.H{"message": "Pod deleted successfully"})

}

func (s *Server) GetLivePods(c *gin.Context) {
	podLists, err := k8s.Clientset.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Bad Request"})
		return
	}
	c.JSON(200, podLists.Items)
}
