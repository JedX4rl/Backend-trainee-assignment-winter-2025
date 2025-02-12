package handler

import (
	"Backend-trainee-assignment-winter-2025/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)

	router.Use(cors.Handler(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
	}))

	//TODO add private routes

	router.Route("/api", func(r chi.Router) {
		r.Post("/auth", h.Auth)
		//r.Get("/info")
		//r.Get("/buy") //TODO: /api/buy/{item}
		//r.Post("/auth")
		//r.Post("/sendCoin")
	})
	return router
}
