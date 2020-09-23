package main

import (
	"github.com/MihaPecnik/order-maching-system/config"
	"github.com/MihaPecnik/order-maching-system/database"
	"github.com/MihaPecnik/order-maching-system/handler"
	"github.com/MihaPecnik/order-maching-system/server"
	"log"
)

func main(){
	db, err := database.NewDatabase(config.GetConfig())
	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(db)
	s := server.NewServer(h)
	log.Fatal(s.Serve(":8080"))
}