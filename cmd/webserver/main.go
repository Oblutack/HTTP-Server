package main

import (
	"log"

	"github.com/Oblutack/HTTP-Server/internal/server"
)

func main() {
	srv := server.NewServer(":8080")

	err := srv.Start()

	if err != nil{
		log.Fatal(err)
	}
}