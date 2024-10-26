package main

import (
	"inventory-service/cmd/api"
	"log"
)

func main() {
	err := api.Run()
	if err != nil {
		log.Fatal(err)
	}
}
