package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func IsAuthenticated(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Send([]byte(err.Error()))
		}
		
		if sess.Get("profile") == nil {
			return c.Redirect("/", http.StatusSeeOther)
		}
		return c.Next()
	}
}
