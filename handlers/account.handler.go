package handlers

import (
	"context"
	"time"

	"ticket-booking/configs/errs"
	"ticket-booking/configs/logs"
	"ticket-booking/dtos/requests"
	"ticket-booking/dtos/responses"
	"ticket-booking/entities"
	"ticket-booking/middlewares"
	"ticket-booking/repositories"
	"ticket-booking/services"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler defines methods for handling auth routes.
type AuthHandler interface {
	SignIn(ctx *fiber.Ctx) error
	SignUp(ctx *fiber.Ctx) error
	Refresh(ctx *fiber.Ctx) error
}

// authHandler is an implementation of AuthHandler that manages authentication routes.
type authHandler struct {
	repository   repositories.AccountRepository
	tokenization services.Tokenization
	cryptography services.Cryptography
}

// SignUp handles the sign-up route, registering a new user.
func (h *authHandler) SignUp(ctx *fiber.Ctx) error {
	context, cancel := h.newContext()
	defer cancel()

	var request requests.SignUpRequest
	if err := ctx.BodyParser(&request); err != nil {
		logs.Error("AuthHandler.SignUp: Failed to parse request body", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	if err := request.Validate(); err != nil {
		logs.Error("AuthHandler.SignUp: Invalid request body", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	if _, err := h.repository.FindByEmail(context, request.Email); err == nil {
		logs.Error("AuthHandler.SignUp: Email already exists", err)
		return errs.NewBadRequest(ctx, "Email already exists")
	}

	hashedPassword, err := h.cryptography.EncryptPassword(request.Password)
	if err != nil {
		logs.Error("AuthHandler.SignUp: Failed to encrypt password", err)
		return errs.NewInternalServerError(ctx, "Failed to sign up")
	}

	newAccount := entities.NewAccount(request.Name, request.Email, hashedPassword)

	account, err := h.repository.SignUp(context, newAccount)
	if err != nil {
		logs.Error("AuthHandler.SignUp: Failed to create user", err)
		return errs.NewInternalServerError(ctx, "Failed to sign up")
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		responses.NewSignUpResponse(
			fiber.StatusCreated,
			"Sign-up successful",
			[]*entities.Account{account},
		))
}

// SignIn handles the sign-in route, authenticating a user.
func (h *authHandler) SignIn(ctx *fiber.Ctx) error {
	context, cancel := h.newContext()
	defer cancel()

	var request requests.SignInRequest
	if err := ctx.BodyParser(&request); err != nil {
		logs.Error("AuthHandler.SignIn: Failed to parse request body", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	if err := request.Validate(); err != nil {
		logs.Error("AuthHandler.SignIn: Invalid request body", err)
		return errs.NewBadRequest(ctx, "Invalid parameter")
	}

	account, err := h.repository.FindByEmail(context, request.Email)
	if err != nil {
		logs.Error("AuthHandler.SignIn: Account not found", err)
		return errs.NewBadRequest(ctx, "Account not found")
	}

	decryptedPassword, err := h.cryptography.VerifyPassword(request.Password, account.Password)
	if err != nil {
		logs.Error("AuthHandler.SignIn: Failed to verify password", err)
		return errs.NewInternalServerError(ctx, "Failed to sign in")
	}

	if !decryptedPassword {
		logs.Error("AuthHandler.SignIn: Incorrect password", err)
		return errs.NewBadRequest(ctx, "Incorrect password")
	}

	token, err := h.tokenization.GenerateToken(account.ID.String())
	if err != nil {
		logs.Error("AuthHandler.SignIn: Failed to generate token", err)
		return errs.NewInternalServerError(ctx, "Failed to sign in")
	}

	return ctx.Status(fiber.StatusOK).JSON(
		responses.NewSignInResponse(
			fiber.StatusOK,
			"Sign-in successful",
			[]*responses.TokenResponse{token},
		))
}

// Refresh handles token refresh requests.
func (h *authHandler) Refresh(ctx *fiber.Ctx) error {
	refreshToken := ctx.Get("Token")
	if refreshToken == "" {
		logs.Error("AuthHandler.Refresh: Missing refresh token in header", nil)
		return errs.NewBadRequest(ctx, "Missing refresh token")
	}

	userId := ctx.Get("UserId")
	if userId == "" {
		logs.Error("AuthHandler.Refresh: Missing user ID in header", nil)
		return errs.NewBadRequest(ctx, "Missing user ID")
	}

	// Validate the refresh token
	valid, err := h.tokenization.VerifyRefreshToken(userId, refreshToken)
	if err != nil || !valid {
		logs.Error("AuthHandler.Refresh: Invalid refresh token", err)
		return errs.NewUnauthorized(ctx, "Invalid refresh token")
	}

	// Generate a new token and refresh token
	tokenResponse, err := h.tokenization.GenerateToken(userId)
	if err != nil {
		logs.Error("AuthHandler.Refresh: Failed to generate new token", err)
		return errs.NewInternalServerError(ctx, "Failed to refresh token")
	}

	return ctx.Status(fiber.StatusOK).JSON(
		responses.NewSignInResponse(
			fiber.StatusOK,
			"Token refreshed successfully",
			[]*responses.TokenResponse{tokenResponse},
		))
}

// newContext creates a new context with a timeout of 5 seconds for database and external calls.
func (h *authHandler) newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// NewAuthHandler initializes a new instance of authHandler and sets up the auth routes.
func NewAuthHandler(router fiber.Router, repository repositories.AccountRepository, tokenization services.Tokenization, cryptography services.Cryptography) AuthHandler {
	handler := &authHandler{
		repository:   repository,
		tokenization: tokenization,
		cryptography: cryptography,
	}

	authRoutes := router.Group("/api/auth")
	authRoutes.Use(middlewares.Logger())

	authRoutes.Post("/signin", handler.SignIn)
	authRoutes.Post("/signup", handler.SignUp)
	authRoutes.Post("/refresh", handler.Refresh)

	return handler
}
