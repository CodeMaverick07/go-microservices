package main

import (
	"fmt"
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
	log.Println("Connected to RabbitMQ")
	//start listening for messages

	//create consumer

	//watch the queue and consume messages
}

func connect() (*amqp.Connection,error) {
	var counts int64
	var connection *amqp.Connection
	var backoff = 1 * time.Second

	for {
		c,err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			fmt.Println("Failed to connect to RabbitMQ. Retrying in ",backoff)
			counts++
		} else {
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

