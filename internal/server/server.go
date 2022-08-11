package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"wbTest/internal/config"
)

type Server interface {
	Run() error
	App() *fiber.App
	Stop() error
}

type server struct {
	app *fiber.App
	cfg *config.Config
}

func (s *server) Run() error {
	return s.app.Listen(fmt.Sprintf(":%d", s.cfg.Port))
}

func (s *server) App() *fiber.App {
	return s.app
}

func (s server) Stop() error {
	return s.app.Shutdown()
}

func NewServer(cfg *config.Config) Server {
	app := fiber.New()
	return &server{
		app: app,
		cfg: cfg,
	}
}
