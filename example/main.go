package main

import (
	"log"
	lc "github.com/amila-ku/go-locust-client"
)

const (
	 hostURL = "http://localhost:8089"
	 users   = 5
	 hatchRate = 1
)

func main(){
	client, err := lc.New(hostURL)
	if err != nil {
		log.Print("first fail")
		log.Fatal(err)
	}

	log.Print("initialized client")

	_, err = client.GenerateLoad(users, hatchRate)
	if err != nil {
		log.Print("failed generating")
		log.Fatal(err)
	}

}