package main

import (
	"log"

	"github.com/MaximPolyaev/gofermart/internal/config"
)

func main() {
	cfg := config.New()
	if err := cfg.Parse(); err != nil {
		log.Fatal(err)
	}
}
