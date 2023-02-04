package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/utils"
)

func NewServer() *fiber.App {
	config := initConfig()
	return fiber.New(config)
}

func NewSession() *session.Store {
	return session.New(session.Config{
		Expiration:   24 * time.Hour,
		KeyLookup:    "cookie:auth-session",
		KeyGenerator: utils.UUID,
	})
}

func initConfig() fiber.Config {
	return fiber.Config{
		Prefork: true,
	}
}
