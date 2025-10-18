package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shivamchaubey027/k8s-cost-optimizer/internal/monitor"
	"github.com/shivamchaubey027/k8s-cost-optimizer/internal/routes"
	"github.com/shivamchaubey027/k8s-cost-optimizer/pkg/database"
	"github.com/shivamchaubey027/k8s-cost-optimizer/pkg/k8s"
)

func main() {

	database.Connect()
	k8s.Connect()

	router := gin.Default()

	routes.SetupRouter(router)
	go monitor.StartMonitoring()
	router.Run("localhost:8080")
}
