package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/luqmanshaban/shipex/cmd"
)


func main() {
	err := godotenv.Load(); if err != nil {
		log.Fatalf("Error loading .env %v",err)
	}
	cmd.Execute()
}