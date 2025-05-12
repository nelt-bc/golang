package main

import (
	"fmt"
	"tron/packages/create2"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	create2.TestCalculateAddress()
}
