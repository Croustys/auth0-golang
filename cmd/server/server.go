package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/utils"
)

func NewServer() *fiber.App {
	return fiber.New()
}

func NewSession() *session.Store {
	return session.New(session.Config{
		Expiration:   24 * time.Hour,
		KeyLookup:    "cookie:auth-session",
		KeyGenerator: utils.UUID,
	})
}
