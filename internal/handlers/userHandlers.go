package handlers

import (
	"marketplace/internal/logger"
	"marketplace/internal/services"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	service services.UserService
}

func NewUserHandler(s services.UserService) *userHandler {
	return &userHandler{service: s}
}

type UserRequest struct {
	Login    string `json:"login" validate:"required,min=4,max=32,email"`
	Password string `json:"password" validate:"required,excludes= ,min=6,max=64"`
}

func (h *userHandler) Register(c echo.Context) error {

	var req UserRequest

	// Перенос json в структуру
	if err := c.Bind(&req); err != nil {
		logger.Logger.Error("failed to bind register request body")
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request"})
	}

	// Проверка на наличие всех полей
	if req.Login == "" || req.Password == "" {
		logger.Logger.Error("invalid register request")
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "login and password required"})
	}

	// Валидация
	if err := c.Validate(&req); err != nil {
		logger.Logger.Error("failed to validate register request body")
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request"})
	}

	user, err := h.service.Register(c.Request().Context(), req.Login, req.Password)
	if err != nil {
		if strings.HasPrefix(err.Error(), "password hashing error") {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *userHandler) Login(c echo.Context) error {

	var req UserRequest

	if err := c.Bind(&req); err != nil {
		logger.Logger.Error("failed to bind login request body")
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request"})
	}

	// Проверка на наличие всех полей
	if req.Login == "" || req.Password == "" {
		logger.Logger.Error("invalid login request")
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "no login or password"})
	}

	token, err := h.service.Login(c.Request().Context(), req.Login, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
