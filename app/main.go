package main

import (
	"pramaan-chain/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		panic("Error loading .env file")
	}

	server.StartAPIServer("2961")
}
