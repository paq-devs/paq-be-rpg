package main

import (
	"log"
	"net/http"

	"github.com/paq-devs/paq-be-rpg/api/routes"
	"github.com/paq-devs/paq-be-rpg/config"
)

// @title Example API
// @version 1.0
// @description This is a sample server.
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @host localhost:8080
// @BasePath /api/v1
// TO-DO: Swagger documentation
func main() {
	config.Init()

	router := routes.RegisterRoutes()
	log.Fatal(http.ListenAndServe(":8080", router))
}
