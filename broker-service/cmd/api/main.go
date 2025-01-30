package main

import (
	"fmt"
	"log"
	"net/http"
)

const WebPort = 80

type Config struct{}

func main() {
	app := Config{}
	log.Println("Starting server on port", WebPort)

	//define http service
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", WebPort),
		Handler: app.routes(),
	} 

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic("error on server", err)
	}

}