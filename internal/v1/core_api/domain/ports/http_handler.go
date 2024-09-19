package ports

import "github.com/gofiber/fiber/v2"

// IHandlers Core Domain HTTP handler interface
type IHandlers interface {
	CreateUser(c *fiber.Ctx) error
	GetUserById(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}
