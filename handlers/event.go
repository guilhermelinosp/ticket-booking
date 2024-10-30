package handlers

import (
	"context"
	"ticket-booking/configs/logs"
	"ticket-booking/configs/validations"
	"ticket-booking/dtos/requests"
	"ticket-booking/dtos/responses"
	"ticket-booking/entities"
	"ticket-booking/repositories/interfaces"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EventHandler struct {
	repository interfaces.EventRepository
}

func (h *EventHandler) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func (h *EventHandler) FindAll(ctx *fiber.Ctx) error {
	context, cancel := h.Context()
	defer cancel()

	events, err := h.repository.FindAll(context)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&responses.BaseResponse{
			Status:  "fail",
			Message: "Failed to retrieve events",
		})
	}

	eventResponses := make([]*responses.EventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = responses.NewEventResponse(event.ID, event.Title, event.Location, event.Date)
	}

	return ctx.Status(fiber.StatusOK).JSON(&responses.BaseResponse{
		Status:  "success",
		Message: "Events retrieved successfully",
		Data:    eventResponses,
	})
}

func (h *EventHandler) FindByID(ctx *fiber.Ctx) error {
	context, cancel := h.Context()
	defer cancel()

	eventId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&responses.BaseResponse{
			Status:  "fail",
			Message: "Invalid ID parameter",
		})
	}

	event, err := h.repository.FindByID(context, eventId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(&responses.BaseResponse{
			Status:  "fail",
			Message: "Event not found",
		})
	}

	eventResponse := responses.NewEventResponse(event.ID, event.Title, event.Location, event.Date)

	return ctx.Status(fiber.StatusOK).JSON(&responses.BaseResponse{
		Status:  "success",
		Message: "Event retrieved successfully",
		Data:    []*responses.EventResponse{eventResponse},
	})
}

func (h *EventHandler) Create(ctx *fiber.Ctx) error {
	context, cancel := h.Context()
	defer cancel()

	var eventRequest *requests.EventRequest
	if err := ctx.BodyParser(&eventRequest); err != nil {
		logs.Error("Error trying to validate user info", err)
		var errs = validations.ValidateRequest(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(&responses.BaseResponse{
			Status:  "fail",
			Message: errs.Message,
		})
	}

	event := entities.NewEvent(eventRequest.Title, eventRequest.Location, eventRequest.Date)
	newEvent, err := h.repository.Create(context, event)
	if err != nil {
		logs.Error("Failed to create event", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(&responses.BaseResponse{
			Status:  "fail",
			Message: "Failed to create event",
		})
	}

	eventResponse := responses.NewEventResponse(newEvent.ID, newEvent.Title, newEvent.Location, newEvent.Date)

	return ctx.Status(fiber.StatusCreated).JSON(&responses.BaseResponse{
		Status:  "success",
		Message: "Event created successfully",
		Data:    []*responses.EventResponse{eventResponse},
	})
}

func (h *EventHandler) Update(ctx *fiber.Ctx) error {
	context, cancel := h.Context()
	defer cancel()

	eventId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&responses.BaseResponse{
			Status:  "fail",
			Message: "Invalid ID parameter",
		})
	}

	var request requests.EventRequest
	if err := ctx.BodyParser(&request); err != nil {
		logs.Error("Error parsing event request body", err)
		var validation = validations.ValidateRequest(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(&responses.BaseResponse{
			Status:  "fail",
			Message: validation.Message,
		})
	}

	event, err := h.repository.FindByID(context, eventId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(&responses.BaseResponse{
			Status:  "fail",
			Message: "Event not found",
		})
	}

	event.Title = request.Title
	event.Date = request.Date
	event.Location = request.Location
	event.UpdatedAt = time.Now()

	updatedEvent, err := h.repository.Update(context, event)
	if err != nil {
		logs.Error("Failed to update event", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(&responses.BaseResponse{
			Status:  "fail",
			Message: "Failed to update event",
		})
	}

	eventResponse := responses.NewEventResponse(updatedEvent.ID, updatedEvent.Title, updatedEvent.Location, updatedEvent.Date)

	return ctx.Status(fiber.StatusOK).JSON(&responses.BaseResponse{
		Status:  "success",
		Message: "Event updated successfully",
		Data:    []*responses.EventResponse{eventResponse},
	})
}

func (h *EventHandler) Delete(ctx *fiber.Ctx) error {
	context, cancel := h.Context()
	defer cancel()

	eventId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&responses.BaseResponse{
			Status:  "fail",
			Message: "Invalid ID parameter",
		})
	}

	event, err := h.repository.FindByID(context, eventId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(&responses.BaseResponse{
			Status:  "fail",
			Message: "Event not found",
		})
	}

	err = h.repository.Delete(context, event.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&responses.BaseResponse{
			Status:  "fail",
			Message: "Failed to delete event",
		})
	}

	return ctx.Status(fiber.StatusNoContent).JSON(&responses.BaseResponse{
		Status:  "success",
		Message: "Event deleted successfully",
	})
}

func NewEventHandler(router fiber.Router, repository interfaces.EventRepository) {
	handler := &EventHandler{
		repository: repository,
	}

	eventRoutes := router.Group("/api/events")
	eventRoutes.Get("/", handler.FindAll)
	eventRoutes.Post("/", handler.Create)
	eventRoutes.Get("/:id", handler.FindByID)
	eventRoutes.Put("/:id", handler.Update)
	eventRoutes.Delete("/:id", handler.Delete)
}
