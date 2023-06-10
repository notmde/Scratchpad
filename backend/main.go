package main

import (
	"log"

	"github.com/scratchpad-backend/server"
	"github.com/scratchpad-backend/storage"
)

func main() {
	store, err := storage.NewDBStore()

	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(":8081", store)
	s.Run()
}
