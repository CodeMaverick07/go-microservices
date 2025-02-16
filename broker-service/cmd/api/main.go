package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const WebPort = 8080

type Config struct{
	Rabbit *amqp.Connection
}

func main() {
	rabbitConn,err := connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v",err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	app := Config{
		Rabbit: rabbitConn,
	}
	log.Println("Starting server on port", WebPort)

	

	//define http service
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", WebPort),
		Handler: app.routes(),
	} 

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic("error on server", err)
	}

}

func connect() (*amqp.Connection,error) {
	var counts int64
	var connection *amqp.Connection
	var backoff = 1 * time.Second

	for {
		c,err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("Failed to connect to RabbitMQ. Retrying in ",backoff)
			counts++
		} else {
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}
		if counts > 5 {
			fmt.Print(err)
			return nil,err
		}

		backoff = time.Duration(math.Pow(float64(counts),2)) * time.Second
		log.Println("baking off ...")
		time.Sleep(backoff)
		continue
	}

	return connection,nil
}