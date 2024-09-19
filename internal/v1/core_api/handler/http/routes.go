package http

import (
	"github.com/gofiber/fiber/v2"
	"main/config"
	"main/internal/v1/core_api/domain/ports"
)

// MapRoutes Auth Domain routes
func MapRoutes(cfg *config.Config, h ports.IHandlers, router fiber.Router) {

	// Örnek Route:

	/*
		jwts := jwtware.New(jwtware.Config{
			SigningKey: []byte(cfg.Server.APP_SECRET),
		})
		//region Common

	*/
	common := router.Group("/auth")
	common.Post("/users", h.CreateUser)           // /auth/users POST
	common.Patch("/users/:user_id", h.UpdateUser) // /auth/users/10 PATCH
}
