package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	"mig_parser/database"
	"mig_parser/handlers"
	"mig_parser/parser"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Register routes
	r.GET("/currencies/latest", handlers.GetLatestCurrencies(db))
	r.GET("/currencies/average", handlers.GetAverageCurrencies(db))

	// Initialize cron job for parsing
	c := cron.New()
	c.AddFunc("0 * * * *", func() { parser.ParseCurrencies(db) }) // Run every hour
	c.Start()

	// Run initial parsing
	go parser.ParseCurrencies(db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
