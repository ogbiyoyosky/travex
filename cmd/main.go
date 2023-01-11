package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	connection "github.com/ogbiyoyosky/travex/db"
	"github.com/ogbiyoyosky/travex/routes"
)

func main() {

	l := log.New(os.Stdout, "travex-api ", log.LstdFlags)

	connection.Connect()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.SetupRoutes(app)
	// start the server

	//seed.RunSeedData()
	go func() {
		l.Println("Starting server on port 9090")

		app.Listen(":9090")
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	app.Shutdown()
}
