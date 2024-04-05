package api

import (
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Init(router fiber.Router)
}

type Server struct {
	handlers []Handler
}

func NewServer(handlers ...Handler) *Server {
	return &Server{
		handlers: handlers,
	}
}

func (s *Server) Init(router fiber.Router) {
	for _, handler := range s.handlers {
		handler.Init(router)
	}
}
