package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/nanyte25/glitchtipctl/cmd"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Run the CLI tool
	cmd.Execute()
}
