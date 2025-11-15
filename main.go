package main

import (
	"awesomeProject2/database"
	"awesomeProject2/item"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	dbConfig := database.NewConfig()
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(db, "./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize service
	itemService := item.NewItemService(db)

	// Initialize handler
	itemHandler := item.NewHandler(itemService)

	// Setup Gin router
	router := gin.Default()

	// API routes
	api := router.Group("/api/v1")
	{
		items := api.Group("/items")
		{
			items.GET("", itemHandler.GetAll)
			items.GET("/:id", itemHandler.GetByID)
			items.POST("", itemHandler.Create)
			items.PUT("/:id", itemHandler.Update)
			items.DELETE("/:id", itemHandler.Delete)
		}
	}

	// Start server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
