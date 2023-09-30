package main

import (
	"fmt"
	"log"

	"github.com/MaximPolyaev/gofermart/internal/config"
	"github.com/MaximPolyaev/gofermart/internal/dbconn"
)

func main() {
	cfg := config.New()
	if err := cfg.Parse(); err != nil {
		log.Fatal(err)
	}

	db, err := dbconn.InitDB(*cfg.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db)
}
