package main

import (
	"fmt"
	"listener-service/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	//try to connect to rabbitmq
	rabbitConn,err := connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v",err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	
	//start listening for messages
	log.Printf("listening for and consuming RabbitMQ messages...")

	//create consumer
    consumer,err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err )
	}

	//watch the queue and consume messages
    err = consumer.Listen([]string{"log.INFO","log.ERROR","log.WARNING"})
    if err != nil {
		panic(err )
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

