package app

import (
	"marketplace/internal/auth"
	"marketplace/internal/config"
	"marketplace/internal/db"
	"marketplace/internal/handlers"
	"marketplace/internal/logger"
	"marketplace/internal/services"
	"marketplace/internal/storage"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pressly/goose/v3"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type App struct {
	config *config.Config
}

func NewApp() *App {
	return &App{
		config: config.NewConfig(),
	}
}

func (a *App) Run() error {

	cfg := config.NewConfig()
	logger.ConfigureLogger(cfg.Logger.Level)

	database, err := db.NewDatabase(cfg)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	logger.Logger.Info("database successfully started")

	if err := goose.Up(database, "internal/db/migrations"); err != nil {
		logger.Logger.Fatalf("migrations error: %v", err)
	}
	logger.Logger.Info("migrations added successfully")

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	store := storage.NewStorage(database)
	jwtUtils := auth.NewJWTProvider(cfg.JWT.Key, cfg.JWT.Lifetime)
	ss := services.NewServices(store, jwtUtils)
	h := handlers.NewHandlers(ss)

	h.SetAPI(e, &jwtUtils)
	e.Logger.Fatal(e.Start(":" + cfg.Server.Port))

	logger.Logger.Infof("server successfully started on %s port", cfg.Server.Port)

	return nil
}
