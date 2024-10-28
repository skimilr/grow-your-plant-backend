package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Plant struct {
	ID          uint       `json:"id"`
	UserID      int        `json:"user_id"`
	GrowthStage string     `json:"growth_stage"`
	HealthLevel int        `json:"health_level"`
	LastWatered *time.Time `json:"last_watered"`
	LastFed     *time.Time `json:"last_fed"`
}

var db *gorm.DB

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	dsn := os.Getenv("DATABASE_URL")
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Plant{})

	// Set up Gin router
	r := gin.Default()

	// Define routes here (API endpoints)

	// Get plant status
	r.GET("/plant/status", getPlantStatus)
	// Water or feed the plant
	r.POST("/plant/action", updatePlant)

	// Start the server
	r.Run(":8080")
}

func getPlantStatus(c *gin.Context) {
	userId := c.Query("user_id") // Retrieve the user ID from the request

	var plant Plant
	if err := db.Where("user_id = ?", userId).First(&plant).Error; err != nil {
		c.JSON(404, gin.H{"error": "Plant not found"})
		return
	}

	c.JSON(200, plant)
}

func updatePlant(c *gin.Context) {
	var action struct {
		UserID int    `json:"user_id"`
		Action string `json:"action"` // "water" or "feed"
	}
	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var plant Plant
	if err := db.Where("user_id = ?", action.UserID).First(&plant).Error; err != nil {
		c.JSON(404, gin.H{"error": "Plant not found"})
		return
	}

	// Update the plant based on the action
	if action.Action == "water" {
		now := time.Now()        // Get the current time
		plant.LastWatered = &now // Assign the address of now
		plant.HealthLevel += 1   // Increase health level
	} else if action.Action == "feed" {
		now := time.Now()      // Get the current time
		plant.LastFed = &now   // Assign the address of now
		plant.HealthLevel += 2 // More health from feeding
	}

	db.Save(&plant) // Save the updated plant
	c.JSON(200, plant)
}
