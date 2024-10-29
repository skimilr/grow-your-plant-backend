package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Plant struct {
	gorm.Model
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

	// Add CORS middleware here
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Content-Type"},
	}))

	// Get plant status
	r.GET("/plant/status", getPlantStatus)
	// Create the plant
	r.POST("/plant/create", createPlant)
	// Water or feed the plant
	r.POST("/plant/action", updatePlant)

	// Start the server
	r.Run(":8080")
}

func getPlantStatus(c *gin.Context) {
	// Get user_id from query parameters
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		log.Println("Missing user_id in query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// Convert user_id to int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		log.Println("Invalid user_id in query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	// Debug statement to log the user_id received
	log.Printf("Received user_id: %d", userID)

	var plant Plant
	if err := db.Where("user_id = ?", userID).First(&plant).Error; err != nil {
		log.Printf("Plant not found for user_id: %d", userID)
		c.JSON(404, gin.H{"error": "Plant not found"})
		return
	}

	c.JSON(200, plant)
}

func createPlant(c *gin.Context) {
	var req struct {
		UserID int `json:"user_id"`
	}

	// Debug statement to log the incoming request body
	log.Printf("Incoming request to create plant: %v", c.Request.Body)

	if err := c.ShouldBindJSON(&req); err != nil || req.UserID <= 0 {
		log.Println("Invalid user_id in request for plant creation")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// Debug statement to log the user_id received
	log.Printf("Received user_id for plant creation: %d", req.UserID)

	// Check if the plant already exists
	var existingPlant Plant
	if err := db.Where("user_id = ?", req.UserID).First(&existingPlant).Error; err == nil {
		log.Printf("Plant already exists for user_id: %d", req.UserID)
		c.JSON(http.StatusConflict, gin.H{"error": "Plant already exists"})
		return
	}

	// Create a new plant with initial values
	newPlant := Plant{
		UserID:      req.UserID,
		GrowthStage: "1",
		HealthLevel: 100,
		LastWatered: nil,
		LastFed:     nil,
	}
	if err := db.Create(&newPlant).Error; err != nil {
		log.Printf("Error creating plant: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plant"})
		return
	}

	c.JSON(http.StatusCreated, newPlant)
}

func updatePlant(c *gin.Context) {
	var action struct {
		UserID int    `json:"user_id"`
		Action string `json:"action"`
	}

	// Debug statement to log the incoming request body
	log.Printf("Incoming request to update plant: %v", c.Request.Body)

	if err := c.ShouldBindJSON(&action); err != nil {
		log.Println("Invalid request for plant update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Debug statement to log the user_id and action received
	log.Printf("Received user_id: %d, action: %s", action.UserID, action.Action)

	var plant Plant
	if err := db.Where("user_id = ?", action.UserID).First(&plant).Error; err != nil {
		log.Printf("Plant not found for user_id: %d", action.UserID)
		c.JSON(404, gin.H{"error": "Plant not found"})
		return
	}

	now := time.Now()
	if action.Action == "water" {
		plant.LastWatered = &now
		if plant.HealthLevel < 10 {
			plant.HealthLevel++
		}
	} else if action.Action == "feed" {
		plant.LastFed = &now
		if plant.HealthLevel < 10 {
			plant.HealthLevel += 2
		}
	}

	if err := db.Save(&plant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plant"})
		return
	}

	c.JSON(http.StatusOK, plant)
}
