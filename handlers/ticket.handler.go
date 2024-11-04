package handlers

import (
	"context"
	"database/sql"
	"strings"
	"ticket-booking/configs/errs"
	"ticket-booking/configs/logs"
	"ticket-booking/dtos/requests"
	"ticket-booking/dtos/responses"
	"ticket-booking/entities"
	"ticket-booking/middlewares"
	"ticket-booking/repositories"
	"ticket-booking/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TicketHandler interface {
	FindAll(ctx *fiber.Ctx) error
	FindByID(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	Validate(ctx *fiber.Ctx) error
}

type ticketHandler struct {
	ticketRepo   repositories.TicketRepository
	eventRepo    repositories.EventRepository
	tokenization services.Tokenization
	cryptography services.Cryptography
}

func (t *ticketHandler) newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func (t *ticketHandler) Validate(ctx *fiber.Ctx) error {
	context, cancel := t.newContext()
	defer cancel()

	token := strings.TrimPrefix(ctx.Get("Authorization"), "Bearer ")
	accountID, err := t.tokenization.GetAccountID(token)
	if err != nil {
		logs.Error("TicketHandler.Validate: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

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

	ticket, err := t.ticketRepo.FindByID(context, id, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "Ticket not found")
		}
		logs.Error("TicketHandler.Validate: Failed to retrieve ticket by ID", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve tickets")
	}

	if ticket.Entered {
		return errs.NewBadRequest(ctx, "Ticket already validated")
	}

	ticket.Entered = true
	ticket.UpdatedAt = time.Now()

	t.ticketRepo.Validate(context, ticket)

	return ctx.Status(fiber.StatusNoContent).JSON(
		responses.NewTicketResponse(
			fiber.StatusNoContent,
			"Ticket updated successfully",
			[]*entities.Ticket{},
		))
}

func (t *ticketHandler) Create(ctx *fiber.Ctx) error {
	context, cancel := t.newContext()
	defer cancel()

	token := strings.TrimPrefix(ctx.Get("Authorization"), "Bearer ")
	accountID, err := t.tokenization.GetAccountID(token)
	if err != nil {
		logs.Error("TicketHandler.Create: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

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

	createdTicket, err := t.ticketRepo.Create(context, entities.NewTicket(event.ID, accountID))
	if err != nil {
		logs.Error("TicketHandler.Create: Failed to create Ticket", err)
		return errs.NewInternalServerError(ctx, "Failed to create Ticket")
	}

	createdTicket.Event = event

	return ctx.Status(fiber.StatusCreated).JSON(
		responses.NewTicketResponse(
			fiber.StatusCreated,
			"Ticket created successfully",
			[]*entities.Ticket{createdTicket},
		))
}

func (t *ticketHandler) Delete(ctx *fiber.Ctx) error {
	context, cancel := t.newContext()
	defer cancel()

	token := strings.TrimPrefix(ctx.Get("Authorization"), "Bearer ")
	accountID, err := t.tokenization.GetAccountID(token)
	if err != nil {
		logs.Error("TicketHandler.Create: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logs.Error("TicketHandler.Delete: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	ticket, err := t.ticketRepo.FindByID(context, id, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "Ticket not found")
		}
		logs.Error("TicketHandler.Delete: Failed to retrieve event by ID", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve tickets")
	}

	err = t.ticketRepo.Delete(context, ticket.ID, ticket.AccountID)
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

func (t *ticketHandler) FindAll(ctx *fiber.Ctx) error {
	context, cancel := t.newContext()
	defer cancel()

	token := strings.TrimPrefix(ctx.Get("Authorization"), "Bearer ")
	accountID, err := t.tokenization.GetAccountID(token)
	if err != nil {
		logs.Error("TicketHandler.Create: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	tickets, err := t.ticketRepo.FindAll(context, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "Ticket not found")
		}

		logs.Error("TicketHandler.FindAll: Failed to retrieve tickets", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve tickets")
	}

	if len(tickets) == 0 {
		return errs.NewNotFound(ctx, "Ticket not found")
	}

	for _, ticket := range tickets {
		event, err := t.eventRepo.FindByID(context, ticket.EventID)
		if err != nil {
			if err == sql.ErrNoRows {
				return errs.NewNotFound(ctx, "Event not found")
			}
			logs.Error("TicketHandler.FindAll: Failed to retrieve event for ticket", err)
			return errs.NewInternalServerError(ctx, "Failed to retrieve events for tickets")
		}

		ticket.Event = event
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.NewTicketResponse(
		fiber.StatusOK,
		"Tickets retrieved successfully",
		tickets,
	))
}

func (t *ticketHandler) FindByID(ctx *fiber.Ctx) error {
	context, cancel := t.newContext()
	defer cancel()

	token := strings.TrimPrefix(ctx.Get("Authorization"), "Bearer ")
	accountID, err := t.tokenization.GetAccountID(token)
	if err != nil {
		logs.Error("TicketHandler.Create: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		logs.Error("TicketHandler.FindByID: Invalid ID parameter", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	ticket, err := t.ticketRepo.FindByID(context, id, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "Ticket not found")
		}
		logs.Error("TicketHandler.FindByID: Failed to retrieve ticket by ID", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve tickets")
	}

	event, err := t.eventRepo.FindByID(context, ticket.EventID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFound(ctx, "Event not found")
		}
		logs.Error("TicketHandler.FindByID: Failed to retrieve event by ID", err)
		return errs.NewInternalServerError(ctx, "Failed to retrieve events")
	}

	ticket.Event = event

	return ctx.Status(fiber.StatusOK).JSON(responses.NewTicketResponse(
		fiber.StatusOK,
		"Ticket retrieved successfully",
		[]*entities.Ticket{ticket},
	))
}

func NewTicketHandler(router fiber.Router, ticketRepo repositories.TicketRepository, eventRepo repositories.EventRepository, tokenization services.Tokenization) TicketHandler {
	handler := &ticketHandler{
		ticketRepo:   ticketRepo,
		eventRepo:    eventRepo,
		tokenization: tokenization,
	}

	ticketRoutes := router.Group("/api/tickets")

	ticketRoutes.Use(middlewares.Logger())
	ticketRoutes.Use(middlewares.Auth(tokenization))

	ticketRoutes.Get("/", handler.FindAll)   
	ticketRoutes.Post("/:id", handler.Create)   // Create a new Ticket
	ticketRoutes.Get("/:id", handler.FindByID)  // Retrieve an Ticket by ID
	ticketRoutes.Delete("/:id", handler.Delete) // Delete an Ticket by ID
	ticketRoutes.Put("/:id", handler.Validate)  // Validate a ticket

	return handler
}
