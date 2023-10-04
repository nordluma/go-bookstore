package main

import (
	"log"
	"sync"

	"github.com/nordluma/go-bookstore/config"
	"github.com/nordluma/go-bookstore/server"
)

func main() {
	log.Println("Starting library server")

	log.Println("Initializing configs")
	err := config.InitConfig("bookstore", nil)
	if err != nil {
		log.Fatalf("Failed to read config: %v\n", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Println("Starting HTTP server")
		err := server.StartHTTPServer()
		if err != nil {
			log.Fatalf("Could not start HTTP server: %v\n", err)
		}

		log.Println("HTTP server gracefully shut down")
	}()
	wg.Wait()

	log.Println("Server stopped")
}
