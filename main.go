package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sampalm/projectapi/database"
	"github.com/sampalm/projectapi/database/migrations"
	"github.com/sampalm/projectapi/server"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	database.StartDB()
	migrations.RunMigrations(database.GetDatabase())

	server := server.NewServer()
	server.Run()
}
