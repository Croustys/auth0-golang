package logout

import (
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
}
