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

// TicketHandler defines methods for handling Ticket routes.
type TicketHandler interface {
	FindAll(ctx *fiber.Ctx) error
	FindByID(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	Validate(ctx *fiber.Ctx) error
}

// TicketHandler handles the Ticket routes.
type ticketHandler struct {
	ticketRepo repositories.TicketRepository
	eventRepo  repositories.EventRepository
}

// NewContext creates a new context with a timeout of 5 seconds.
func (t *ticketHandler) newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// Valider implements TicketHandler.
func (t *ticketHandler) Validate(ctx *fiber.Ctx) error {
	context, cancel := t.newContext()
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logs.Error("TicketHandler.Validate: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	var request requests.EventRequest
	if err := ctx.BodyParser(&request); err != nil {
		logs.Error("TicketHandler.Update: Failed to parse request body", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	event, err := t.ticketRepo.FindByID(context, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "Ticket not found")
		}
		logs.Error("TicketHandler.Validate: Failed to retrieve ticket by ID", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve tickets")
	}

	if event.Entered {
		return errs.NewBadRequest(ctx, "Ticket already validated")
	}

	event.Entered = true
	event.UpdatedAt = time.Now()

	updatedEvent, err := t.ticketRepo.Validate(context, event)
	if err != nil {
		logs.Error("TicketHandler.Validate: Failed to update ticket", err)
		return errs.NewInternalServerError(ctx, "Failed to update ticket")
	}

	return ctx.Status(fiber.StatusOK).JSON(
		responses.NewTicketResponse(
			fiber.StatusOK,
			"Ticket updated successfully",
			[]*entities.Ticket{updatedEvent},
		))
}

// Create implements TicketHandler.
func (t *ticketHandler) Create(ctx *fiber.Ctx) error {
	context, cancel := t.newContext()
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logs.Error("TicketHandler.Create: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	event, err := t.eventRepo.FindByID(context, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "Event not found")
		}
		logs.Error("TicketHandler.Create: Failed to retrieve event by ID", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve events")
	}

	createdTicket, err := t.ticketRepo.Create(context, entities.NewTicket(event.ID))
	if err != nil {
		logs.Error("TicketHandler.Create: Failed to create Ticket", err)
		return errs.NewInternalServerError(ctx, "Failed to create Ticket")
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		responses.NewTicketResponse(
			fiber.StatusCreated,
			"Ticket created successfully",
			[]*entities.Ticket{createdTicket},
		))
}

// Delete implements TicketHandler.
func (t *ticketHandler) Delete(ctx *fiber.Ctx) error {
	context, cancel := t.newContext()
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logs.Error("TicketHandler.Delete: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	event, err := t.ticketRepo.FindByID(context, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "Ticket not found")
		}
		logs.Error("TicketHandler.Delete: Failed to retrieve event by ID", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve tickets")
	}

	err = t.ticketRepo.Delete(context, event.ID)
	if err != nil {
		logs.Error("TicketHandler.Delete: Failed to delete ticket", err)
		return errs.NewInternalServerError(ctx, "Failed to delete ticket")
	}

	return ctx.Status(fiber.StatusOK).JSON(
		responses.NewTicketResponse(
			fiber.StatusOK,
			"Ticket deleted successfully",
			nil,
		))
}

// FindAll implements TicketHandler.
func (t *ticketHandler) FindAll(ctx *fiber.Ctx) error {
	context, cancel := t.newContext()
	defer cancel()

	events, err := t.ticketRepo.FindAll(context)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "No tickets found")
		}

		logs.Error("TicketHandler.FindAll: Failed to retrieve tickets", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve tickets")
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.NewTicketResponse(
		fiber.StatusOK,
		"Tickets retrieved successfully",
		events,
	))
}

// FindByID implements TicketHandler.
func (t *ticketHandler) FindByID(ctx *fiber.Ctx) error {
	context, cancel := t.newContext()
	defer cancel()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logs.Error("TicketHandler.FindByID: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	event, err := t.ticketRepo.FindByID(context, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "Ticket not found")
		}
		logs.Error("TicketHandler.FindByID: Failed to retrieve ticket by ID", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve tickets")
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.NewTicketResponse(
		fiber.StatusOK,
		"Ticket retrieved successfully",
		[]*entities.Ticket{event},
	))
}

// NewTicketHandler creates a new instance of TicketHandler and sets up the Ticket routes.
func NewTicketHandler(router fiber.Router, ticketRepo repositories.TicketRepository, eventRepo repositories.EventRepository) TicketHandler {
	handler := &ticketHandler{
		ticketRepo: ticketRepo,
		eventRepo:  eventRepo,
	}

	TicketRoutes := router.Group("/api/tickets")

	TicketRoutes.Use(middlewares.Logger())

	TicketRoutes.Get("/", handler.FindAll)      // Retrieve all Tickets
	TicketRoutes.Post("/:id", handler.Create)   // Create a new Ticket
	TicketRoutes.Get("/:id", handler.FindByID)  // Retrieve an Ticket by ID
	TicketRoutes.Delete("/:id", handler.Delete) // Delete an Ticket by ID
	TicketRoutes.Put("/:id", handler.Validate)  // Validate a ticket

	return handler
}
