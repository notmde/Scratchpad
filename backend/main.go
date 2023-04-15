package main

import (
	"log"

	"github.com/scratchpad-backend/api"
	"github.com/scratchpad-backend/storage"
)

func main() {
	store, err := storage.NewDBStore()

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(":8081", store)

	server.Run()
}
