package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type GrowthStage string

const (
	Seed      GrowthStage = "Seed"
	Seedling  GrowthStage = "Seedling"
	Young     GrowthStage = "Young"
	Youngling GrowthStage = "Youngling"
	Bloom     GrowthStage = "Bloom"
	Blooming  GrowthStage = "Blooming"
	Decayed   GrowthStage = "Decayed"
	Dead      GrowthStage = "Dead"
)

type PlantType string

const (
	FloweringPlant    PlantType = "Flowering"
	NonFloweringPlant PlantType = "NonFlowering"
)

type Plant struct {
	ID          int
	Type        PlantType
	GrowthStage GrowthStage
	HealthLevel int
	LastWatered time.Time
	LastFed     time.Time
}

var plants = make(map[int][]*Plant)
var plantMutex = sync.Mutex{}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/plant/create", createPlant)
	r.POST("/plant/action", updatePlant)
	r.DELETE("/plant/delete", deletePlant)
	r.GET("/plant/status", getPlantStatus)

	go plantDecayRoutine()

	r.Run(":8080")
}

func createPlant(c *gin.Context) {
	var req struct {
		UserID int       `json:"user_id"`
		Type   PlantType `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and type are required"})
		return
	}

	plantMutex.Lock()
	defer plantMutex.Unlock()

	newPlant := &Plant{
		ID:          len(plants[req.UserID]) + 1,
		Type:        req.Type,
		GrowthStage: Seed,
		HealthLevel: 100,
		LastWatered: time.Now(),
		LastFed:     time.Now(),
	}
	plants[req.UserID] = append(plants[req.UserID], newPlant)

	c.JSON(http.StatusCreated, newPlant)
}

func updatePlant(c *gin.Context) {
	var req struct {
		UserID  int    `json:"user_id"`
		Action  string `json:"action"`
		PlantID int    `json:"plant_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	plantMutex.Lock()
	defer plantMutex.Unlock()

	plantSlice, exists := plants[req.UserID]
	if !exists || len(plantSlice) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plant not found"})
		return
	}

	var plant *Plant
	for _, p := range plantSlice {
		if p.ID == req.PlantID {
			plant = p
			break
		}
	}

	if plant == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plant not found"})
		return
	}

	switch req.Action {
	case "water":
		plant.LastWatered = time.Now()
		plant.HealthLevel = min(plant.HealthLevel+10, 100)
	case "feed":
		plant.LastFed = time.Now()
		plant.HealthLevel = min(plant.HealthLevel+10, 100)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}

	if plant.HealthLevel > 0 {
		if plant.GrowthStage == Decayed {
			plant.GrowthStage = getGrowthStageFromHealth(plant.HealthLevel)
		}
	}

	c.JSON(http.StatusOK, plant)
}

func deletePlant(c *gin.Context) {
	var req struct {
		UserID  int `json:"user_id"`
		PlantID int `json:"plant_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	plantMutex.Lock()
	defer plantMutex.Unlock()

	plantSlice, exists := plants[req.UserID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	for i, p := range plantSlice {
		if p.ID == req.PlantID {
			plants[req.UserID] = append(plantSlice[:i], plantSlice[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Plant deleted successfully"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Plant not found"})
}

func getGrowthStageFromHealth(healthLevel int) GrowthStage {
	if healthLevel >= 80 {
		return Blooming
	} else if healthLevel >= 60 {
		return Bloom
	} else if healthLevel >= 40 {
		return Youngling
	} else if healthLevel >= 20 {
		return Young
	} else {
		return Seedling
	}
}

func getPlantStatus(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	plantMutex.Lock()
	defer plantMutex.Unlock()

	plantSlice, exists := plants[userID]
	if !exists || len(plantSlice) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plant not found"})
		return
	}

	c.JSON(http.StatusOK, plantSlice)
}

func plantDecayRoutine() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		plantMutex.Lock()
		for _, plantSlice := range plants {
			for _, plant := range plantSlice {
				if plant.HealthLevel > 0 {
					decayHealthIfNecessary(plant)
					updateGrowthStage(plant)
				}
			}
		}
		plantMutex.Unlock()
	}
}

func decayHealthIfNecessary(plant *Plant) {
	if time.Since(plant.LastWatered) > 10*time.Second || time.Since(plant.LastFed) > 10*time.Second {
		plant.HealthLevel -= 10
		if plant.HealthLevel <= 0 {
			plant.HealthLevel = 0
			plant.GrowthStage = Dead
		} else {
			plant.GrowthStage = Decayed
		}
	}
}

func updateGrowthStage(plant *Plant) {
	if plant.GrowthStage == Decayed || plant.GrowthStage == Dead {
		return
	}

	switch plant.GrowthStage {
	case Seed:
		plant.GrowthStage = Seedling
	case Seedling:
		plant.GrowthStage = Young
	case Young:
		plant.GrowthStage = Youngling
	case Youngling:
		plant.GrowthStage = Bloom
	case Bloom:
		plant.GrowthStage = Blooming
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
