package callback

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"auth0/pkg/authenticator"
)

func Handler(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if ctx.Query("state") != session.Get("state") {
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Failed to convert an authorization code into a token.")
			return
		}

		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Redirect to logged in page.
		ctx.Redirect(http.StatusTemporaryRedirect, "/user")
	}
}

func Handle(auth *authenticator.Authenticator, store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			panic(err)
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(c.Context(), c.Query("code"))
		if err != nil {
			return c.SendStatus(http.StatusUnauthorized)
		}

		idToken, err := auth.VerifyIDToken(c.Context(), token)
		if err != nil {
			log.Println(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			log.Println(err.Error())

			return c.SendStatus(http.StatusInternalServerError)
		}

		sess.Set("access_token", token.AccessToken)
		sess.Set("profile", profile)
		if err := sess.Save(); err != nil {
			log.Println(err.Error())

			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.Redirect("/user", http.StatusTemporaryRedirect)
	}
}
