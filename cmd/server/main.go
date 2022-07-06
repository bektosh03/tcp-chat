package main

import (
	"chat/pkg/server"
	"log"
)

func main() {
	s := server.New()
	if err := s.Run("127.0.0.1:8000"); err != nil {
		log.Panicf("failed to run server: %v\n", err)
	}
}
