package router

import (
	"auth0/pkg/authenticator"
	"auth0/pkg/middleware"
	"auth0/pkg/utils"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store
var auth *authenticator.Authenticator

func New(app *fiber.App, sessionStore *session.Store) {
	router := app.Group("/")

	auth, _ = authenticator.New()
	store = sessionStore

	router.Get("/login", loginHandler)
	router.Get("/callback", callbackHandler)
	router.Get("/user", middleware.IsAuthenticated(store), userHandler)
	router.Get("/logout", logoutHandler)
}

func loginHandler(c *fiber.Ctx) error {
	state, err := utils.GenerateRandomState()
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	sess, err := store.Get(c)
	if err != nil {

		return c.Send([]byte(err.Error()))
	}

	sess.Set("state", state)
	if err := sess.Save(); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Redirect(auth.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func callbackHandler(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	token, err := auth.Exchange(c.Context(), c.Query("code"))
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	idToken, err := auth.VerifyIDToken(c.Context(), token)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	sess.Set("access_token", token.AccessToken)
	sess.Set("profile", profile)
	if err := sess.Save(); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Redirect("/user", http.StatusTemporaryRedirect)
}

func userHandler(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Send([]byte(err.Error()))
	}

	profile := sess.Get("profile")
	return c.JSON(profile)
}

func logoutHandler(c *fiber.Ctx) error {
	logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	scheme := "http"
	if c.Context().IsTLS() {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + c.Hostname())
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	return c.Redirect(logoutUrl.String(), http.StatusTemporaryRedirect)
}
