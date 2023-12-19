package main

import (
	"log"
	"payment/internal/api"
)

func main() {
	log.Println("App start")
	api.StartServer()
	log.Println("App stop")
}
