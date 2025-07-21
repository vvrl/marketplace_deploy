package handlers

import (
	"marketplace/internal/auth"
	"marketplace/internal/middlewares"
	"marketplace/internal/services"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	UserHandler *userHandler
	AdHandler   *adHandler
}

func NewHandlers(s *services.Services) *Handlers {
	return &Handlers{
		UserHandler: NewUserHandler(s.UserService),
		AdHandler:   NewAdHandler(s.AdService),
	}
}

func (h *Handlers) SetAPI(e *echo.Echo, jwtUtils *auth.JWTProvider) {
	e.Use(middlewares.JwtMiddleware(*jwtUtils))

	// User methods
	e.POST("/register", h.UserHandler.Register)
	e.GET("/login", h.UserHandler.Login)

	// Advertisement methods
	e.POST("/postAd", h.AdHandler.PostAd)
	e.GET("/getList", h.AdHandler.GetAdList)
}
