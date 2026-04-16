package app

import (
	"ChipHouse/cmd/internal/config"
	"ChipHouse/pkg/logging"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type App struct {
	config *config.Config
	logger *logging.Logger
}

func NewApp(config *config.Config, logger *logging.Logger) (App, error) {
	logger.Println("Router initializing")
	router := httprouter.New()

	logger.Println("Swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	return App{
		config: config,
		logger: logger,
	}, nil
}
