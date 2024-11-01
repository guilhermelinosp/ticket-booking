package handlers

import (
	"context"
	"database/sql"
	"ticket-booking/configs/errs"
	"ticket-booking/configs/logs"
	"ticket-booking/dtos/requests"
	"ticket-booking/dtos/responses"
	"ticket-booking/entities"
	"ticket-booking/middlewares"
	"ticket-booking/repositories"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// EventHandler handles the event routes.
type EventHandler struct {
	repository repositories.EventRepository
}

// NewEventHandler creates a new instance of EventHandler and sets up the event routes.
func NewEventHandler(router fiber.Router, repository repositories.EventRepository) {
	handler := &EventHandler{
		repository: repository,
	}

	eventRoutes := router.Group("/api/events")

	eventRoutes.Use(middlewares.Logger())

	eventRoutes.Get("/", handler.FindAll)      // Retrieve all events
	eventRoutes.Post("/", handler.Create)      // Create a new event
	eventRoutes.Get("/:id", handler.FindByID)  // Retrieve an event by ID
	eventRoutes.Put("/:id", handler.Update)    // Update an event by ID
	eventRoutes.Delete("/:id", handler.Delete) // Delete an event by ID
}

// NewContext creates a new context with a timeout of 5 seconds.
func (h *EventHandler) NewContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// FindAll retrieves all events.
func (h *EventHandler) FindAll(ctx *fiber.Ctx) error {
	context, cancel := h.NewContext()
	defer cancel()

	events, err := h.repository.FindAll(context)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "No events found")
		}

		logs.Error("EventHandler.FindAll: Failed to retrieve events", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve events")
	}

	return ctx.Status(fiber.StatusOK).JSON(&responses.BaseResponse{
		Status:  fiber.StatusOK,
		Message: "Events retrieved successfully",
		Data:    events,
	})
}

// FindByID retrieves an event by its ID.
func (h *EventHandler) FindByID(ctx *fiber.Ctx) error {
	context, cancel := h.NewContext()
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logs.Error("EventHandler.FindByID: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	event, err := h.repository.FindByID(context, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "Event not found")
		}
		logs.Error("EventHandler.FindByID: Failed to retrieve event by ID", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve events")
	}

	return ctx.Status(fiber.StatusOK).JSON(&responses.BaseResponse{
		Status:  fiber.StatusOK,
		Message: "Events retrieved successfully",
		Data:    []*entities.Event{event},
	})
}

// Create creates a new event.
func (h *EventHandler) Create(ctx *fiber.Ctx) error {
	context, cancel := h.NewContext()
	defer cancel()

	var request requests.EventRequest
	if err := ctx.BodyParser(&request); err != nil {
		logs.Error("EventHandler.Create: Failed to parse request body", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	// Validate the request
	if err := request.Validate(); err != nil {
		logs.Error("EventHandler.Create: Failed to parse request body", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	newEvent := entities.NewEvent(request.Title, request.Location, request.Date)
	createdEvent, err := h.repository.Create(context, newEvent)
	if err != nil {
		logs.Error("EventHandler.Create: Failed to create event", err)
		return errs.NewInternalServerError(ctx, "Failed to create event")
	}

	return ctx.Status(fiber.StatusCreated).JSON(&responses.BaseResponse{
		Status:  fiber.StatusCreated,
		Message: "Event created successfully",
		Data:    []*entities.Event{createdEvent},
	})
}

// Update updates an event by its ID.
func (h *EventHandler) Update(ctx *fiber.Ctx) error {
	context, cancel := h.NewContext()
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logs.Error("EventHandler.Update: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	var request requests.EventRequest
	if err := ctx.BodyParser(&request); err != nil {
		logs.Error("EventHandler.Update: Failed to parse request body", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	event, err := h.repository.FindByID(context, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "Event not found")
		}
		logs.Error("EventHandler.Update: Failed to retrieve event by ID", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve events")
	}

	if request.Title != "" {
		event.Title = request.Title
	}
	if request.Location != "" {
		event.Location = request.Location
	}
	if !request.Date.IsZero() {
		event.Date = request.Date
	}

	event.UpdatedAt = time.Now()

	updatedEvent, err := h.repository.Update(context, event)
	if err != nil {
		logs.Error("EventHandler.Update: Failed to update event", err)
		return errs.NewInternalServerError(ctx, "Failed to update event")
	}

	return ctx.Status(fiber.StatusOK).JSON(&responses.BaseResponse{
		Status:  fiber.StatusOK,
		Message: "Event updated successfully",
		Data:    []*entities.Event{updatedEvent},
	})
}

// Delete deletes an event by its ID.
func (h *EventHandler) Delete(ctx *fiber.Ctx) error {
	context, cancel := h.NewContext()
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logs.Error("EventHandler.Delete: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	event, err := h.repository.FindByID(context, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "Event not found")
		}
		logs.Error("EventHandler.Delete: Failed to retrieve event by ID", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve events")
	}

	err = h.repository.Delete(context, event.ID)
	if err != nil {
		logs.Error("EventHandler.Delete: Failed to delete event", err)
		return errs.NewInternalServerError(ctx, "Failed to delete event")
	}

	return ctx.Status(fiber.StatusOK).JSON(&responses.BaseResponse{
		Status:  fiber.StatusOK,
		Message: "Event deleted successfully",
	})
}
