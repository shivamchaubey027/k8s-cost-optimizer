package monitor

import (
	"context"
	"log"
	"time"

	"github.com/shivamchaubey027/k8s-cost-optimizer/internal/models"
	"github.com/shivamchaubey027/k8s-cost-optimizer/pkg/database"
	"github.com/shivamchaubey027/k8s-cost-optimizer/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func StartMonitoring() {
	ticker := time.NewTicker(10 * time.Second)

	for range ticker.C {
		podLists, err := k8s.Clientset.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Error %v", err)
		}
		for _, k8sPod := range podLists.Items {

			dbPod := models.Pod{
				Name:          k8sPod.Name,
				Namespace:     k8sPod.Namespace,
				CPURequest:    k8sPod.Spec.Containers[0].Resources.Requests.Cpu().String(),
				MemoryRequest: k8sPod.Spec.Containers[0].Resources.Requests.Memory().String(),
				IpAddress:     k8sPod.Status.PodIP,
			}

			result := database.DB.Create(&dbPod)
			if result.Error != nil {
				log.Printf("ERROR (DB Save): %v", result.Error)
			}
		}

		log.Printf("Duty Completed")
	}
}
