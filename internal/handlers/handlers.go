package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shivamchaubey027/k8s-cost-optimizer/internal/models"
	"github.com/shivamchaubey027/k8s-cost-optimizer/pkg/database"
	"github.com/shivamchaubey027/k8s-cost-optimizer/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	cpuCostPerCoreHour  = 0.04
	memoryCostPerGBHour = 0.0416
)

type Server struct {
}

type PodAggregate struct {
	Count         int
	TotalCPUCores float64
	TotalMemoryGB float64
}
type Recommendation struct {
	PodName                string  `json:"pod_name"`
	AverageCPURequestCores float64 `json:"average_cpu_request_cores"`
	AverageMemoryRequestGB float64 `json:"average_memory_request_gb"`
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

func (s *Server) GetSavings(c *gin.Context) {
	var allPods []models.Pod
	var totalCPURequested float64
	var totalMemoryRequestedGB float64
	var totalEstimatedHourlyCost float64

	result := database.DB.Find(&allPods)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch pod data"})
		return
	}

	for _, pod := range allPods {
		var cpuCores float64
		if strings.HasSuffix(pod.CPURequest, "m") {
			milliCoresStr := strings.TrimSuffix(pod.CPURequest, "m")
			milliCores, err := strconv.ParseFloat(milliCoresStr, 64)
			if err == nil {
				cpuCores = milliCores / 1000
			}
		} else {
			cores, err := strconv.ParseFloat(pod.CPURequest, 64)
			if err == nil {
				cpuCores = cores
			}
		}
		totalCPURequested += cpuCores

		var memoryGB float64
		if strings.HasSuffix(pod.MemoryRequest, "Mi") {
			mebiBytesStr := strings.TrimSuffix(pod.MemoryRequest, "Mi")
			mebiBytes, err := strconv.ParseFloat(mebiBytesStr, 64)
			if err == nil {
				memoryGB = mebiBytes / 1024
			}
		} else if strings.HasSuffix(pod.MemoryRequest, "Gi") {
			gibiBytesStr := strings.TrimSuffix(pod.MemoryRequest, "Gi")
			gibiBytes, err := strconv.ParseFloat(gibiBytesStr, 64)
			if err == nil {
				memoryGB = gibiBytes
			}
		}
		totalMemoryRequestedGB += memoryGB

		podHourlyCost := (cpuCores * cpuCostPerCoreHour) + (memoryGB * memoryCostPerGBHour)
		totalEstimatedHourlyCost += podHourlyCost

	}
	c.JSON(200, gin.H{
		"total_pods_records":          len(allPods),
		"total_cpu_requested_cores":   fmt.Sprintf("%.2f", totalCPURequested),
		"total_memory_requested_gb":   fmt.Sprintf("%.2f GB", totalMemoryRequestedGB),
		"total_estimated_hourly_cost": fmt.Sprintf("$%.2f", totalEstimatedHourlyCost),
	})
}

func (s *Server) GetRecommendations(c *gin.Context) {
	var allPods []models.Pod

	result := database.DB.Find(&allPods)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Bad Request"})
		return
	}

	podDataMap := make(map[string]PodAggregate)
	for _, pod := range allPods {
		var cpuCores float64
		if strings.HasSuffix(pod.CPURequest, "m") {
			milliCoresStr := strings.TrimSuffix(pod.CPURequest, "m")
			milliCores, err := strconv.ParseFloat(milliCoresStr, 64)
			if err == nil {
				cpuCores = milliCores / 1000
			}
		} else {
			cores, err := strconv.ParseFloat(pod.CPURequest, 64)
			if err == nil {
				cpuCores = cores
			}
		}

		var memoryGB float64
		if strings.HasSuffix(pod.MemoryRequest, "Mi") {
			mebiBytesStr := strings.TrimSuffix(pod.MemoryRequest, "Mi")
			mebiBytes, err := strconv.ParseFloat(mebiBytesStr, 64)
			if err == nil {
				memoryGB = mebiBytes / 1024
			}
		} else if strings.HasSuffix(pod.MemoryRequest, "Gi") {
			gibiBytesStr := strings.TrimSuffix(pod.MemoryRequest, "Gi")
			gibiBytes, err := strconv.ParseFloat(gibiBytesStr, 64)
			if err == nil {
				memoryGB = gibiBytes
			}
		}

		aggregate := podDataMap[pod.Name]
		aggregate.Count++
		aggregate.TotalCPUCores += cpuCores
		aggregate.TotalMemoryGB += memoryGB

		podDataMap[pod.Name] = aggregate
	}

	var recommendations []Recommendation

	for podName, aggregateData := range podDataMap {

		if aggregateData.Count == 0 {
			continue
		}

		avgCPU := aggregateData.TotalCPUCores / float64(aggregateData.Count)
		avgMemory := aggregateData.TotalMemoryGB / float64(aggregateData.Count)

		rec := Recommendation{
			PodName:                podName,
			AverageCPURequestCores: parseFloat(fmt.Sprintf("%.3f", avgCPU)),
			AverageMemoryRequestGB: parseFloat(fmt.Sprintf("%.3f", avgMemory)),
		}

		recommendations = append(recommendations, rec)
	}

	c.JSON(200, recommendations)
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
