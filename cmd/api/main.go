package main

import (
	"fmt"
	"go-server/internal/config"
	"go-server/internal/db"
	"go-server/internal/server"
	"log"
)
func main(){
	config,err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration %v",err)
	}
	fmt.Println("Config loaded successfully:", config)

	database,err := db.Connect(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer func (){
		if err := db.Disconnect(database); err != nil {
			log.Fatalf("Failed to disconnect from database: %v", err)
		}
	}()
	fmt.Println("Database connected successfully")

	router := server.NewRouter()

	addr := fmt.Sprintf(":%s", config.ServerPort)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}