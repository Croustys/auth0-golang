package main

import (
	"encoding/gob"
	"log"

	"github.com/joho/godotenv"

	"auth0/cmd/server"
	"auth0/internal/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	gob.Register(map[string]interface{}{})
	app := server.NewServer()
	store := server.NewSession()

	router.New(app, store)

	log.Print("Server listening on http://localhost:3000/")
	log.Fatal(app.Listen("0.0.0.0:3000"))
}
