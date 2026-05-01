package main

import (
	//"database/sql"
	"log"
	"os"

	"github.com/go-jet/jet/generator/postgres"
	_ "github.com/lib/pq"
)

func main() {
	err := postgres.Generate("../gen",
		postgres.DBConnection{
			Host:       "/var/run/postgresql",
			Port:       5432,
			User:       "eliribble",
			Password:   "none",
			DBName:     "nidus-sync",
			SchemaName: "stadia",
			SslMode:    "disable",
		})
	if err != nil {
		log.Printf("Failed: %v", err)
		os.Exit(1)
	}
}
