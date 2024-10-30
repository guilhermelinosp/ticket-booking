package main

import (
	"log"
	"ticket-booking/configs"
	"ticket-booking/handlers"
	"ticket-booking/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create database connections using the configured functions
	reader := configs.GetReaderSqlx()
	writer := configs.GetWriterSqlx()

	// Ensure connections are closed when the application exits
	defer func() {
		if err := reader.Close(); err != nil {
			log.Println("Error closing reader:", err)
		}
	}()
	defer func() {
		if err := writer.Close(); err != nil {
			log.Println("Error closing writer:", err)
		}
	}()

	// Initialize the Fiber application
	app := fiber.New(fiber.Config{
		AppName:      "Ticket-Booking",
		ServerHeader: "Fiber",
	})

	// Initialize repositories
	eventRepository := repositories.NewEventRepository(reader, writer)

	// Set up handlers
	handlers.NewEventHandler(app, eventRepository)

	// Start the server on the specified port (defaulting to 3000)
	port := ":3000" // You can replace with your port variable from envConfig.ServerPort
	if err := app.Listen(port); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
