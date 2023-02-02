package main

import (
	"encoding/gob"
	"log"

	//"net/http"

	"github.com/joho/godotenv"

	"auth0/cmd/server"
	"auth0/pkg/authenticator"
	"auth0/pkg/callback"
	"auth0/pkg/login"
	"auth0/pkg/logout"
	"auth0/pkg/middleware"
	"auth0/pkg/user"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	gob.Register(map[string]interface{}{})
	app := server.NewServer()
	store := server.NewSession()

	app.Get("/login", login.Handle(auth, store))
	app.Get("/callback", callback.Handle(auth, store))
	app.Get("/user", middleware.IsAuthenticated(store), user.Handler(store))
	app.Get("/logout", logout.Handler())

	log.Print("Server listening on http://localhost:3000/")
	log.Fatal(app.Listen("0.0.0.0:3000"))
}
