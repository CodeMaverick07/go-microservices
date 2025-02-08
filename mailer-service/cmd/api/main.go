package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const webPort = 80

type Config struct {
	Mailer Mail
}

func main() {
	app := Config{
		Mailer: createMail(),
	}

	log.Printf("Starting mail service on port %d", webPort)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", webPort),
		Handler: app.routes(),	
	}
	err := srv.ListenAndServe()

	if err != nil {
		
		log.Fatalf("server failed to start: %v", err)
		log.Panic(err)
	}

}

func createMail() Mail {
port,_:= strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain: os.Getenv("MAIL_DOMAIN"),
		Host: os.Getenv("MAIL_HOST"),
		Port :port ,
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		Encryption: os.Getenv("MAIL_ENCRYPTION"),
		FromName: os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}
	return m
}