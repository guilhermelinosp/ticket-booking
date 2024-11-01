package main

import (
	"ticket-booking/configs"
	"ticket-booking/configs/logs"
	"ticket-booking/handlers"
	"ticket-booking/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		logs.Fatal("Error loading .env file", err)
	}

	// Create database connections using the configured functions
	reader := configs.GetReaderSqlx()
	writer := configs.GetWriterSqlx()

	// Ensure connections are closed when the application exits
	defer func() {
		if err := reader.Close(); err != nil {
			logs.Error("Error closing reader", err)
		}
	}()
	defer func() {
		if err := writer.Close(); err != nil {
			logs.Error("Error closing writer", err)
		}
	}()

	// Initialize the Fiber application
	app := fiber.New(fiber.Config{
		AppName:      "Ticket-Booking",
		ServerHeader: "Fiber",
	})

	// Initialize repositories
	eventRepo := repositories.NewEventRepository(reader, writer)
	ticketRepo := repositories.NewTicketRepository(reader, writer)

	// Set up handlers
	handlers.NewEventHandler(app, eventRepo)
	handlers.NewTicketHandler(app, ticketRepo, eventRepo)

	// Start the server on the specified port (defaulting to 3000)
	port := ":3000" // You can replace with your port variable from envConfig.ServerPort
	logs.Info("Starting server on port", zap.String("port", port))
	if err := app.Listen(port); err != nil {
		logs.Fatal("Error starting server", err)
	}
}
