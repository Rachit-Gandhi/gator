package main

import (
	"fmt"
	"log"

	"github.com/Rachit-Gandhi/gator/internal/config"
)

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	err = config.SetUser("rachit")
	if err != nil {
		log.Fatalf("error setting user: %v", err)
	}
	fmt.Println(config)
}
