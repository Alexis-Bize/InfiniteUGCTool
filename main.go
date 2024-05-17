package main

import (
	auth_service "Infinite-Bookmarker/internal/application/services/auth"
	"fmt"
	"log"
)

// WIP / DEBUG
func main() {
	// Provide a config file to store user credentials
	spartanToken, err := auth_service.Authenticate("your_email", "your_password")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(spartanToken)
}