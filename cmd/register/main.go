package main

import (
	"github.com/joho/godotenv"
	"github.com/kaanaktas/openbanking-registration/internal/client"
	"github.com/kaanaktas/openbanking-registration/pkg/register"
	"log"
	"os"
)

func init() {
	_ = godotenv.Load()
	client.LoadCache()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := client.NewEchoEngine()

	// Routes
	e.GET("/:aspsp/register", register.Register)

	log.Printf("starting server at :%s", port)

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("error while starting server at :%s, %v", port, err)
	}
}
