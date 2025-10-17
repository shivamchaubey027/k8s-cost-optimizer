package database

import (
	"fmt"
	"log"

	"github.com/shivamchaubey027/k8s-cost-optimizer/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=user password=password dbname=k8s_cost_optimizer port=5432 sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect db", err)
	}
	fmt.Println("Successfully connected")

	DB.AutoMigrate(&models.Pod{})
}
