package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	WebPort = 80
	rpcPort = 5001
	mongoURL= "mongodb://mongo:27017"
	gRpcPort = 50001
)

var client *mongo.Client

type Config struct{}

func main() {
  mongoClient, err := connectToMongo()

  if err != nil {
	panic(err)
  }
  client = mongoClient

  ctx , cancel := context.WithCancel(context.Background())
  defer cancel()

  defer func() {
	if err := client.Disconnect(ctx); err != nil {
	  panic(err)
	}
  }()

}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
	  Username: "admin",
	  Password: "password",
	})
	
	c, err := mongo.Connect( clientOptions)
	if err != nil {
	  log.Println("Error connecting to MongoDB", err)
	  return nil, err
	}
	
	// Verify the connection
	err = c.Ping(context.TODO(), nil)
	if err != nil {
	  log.Println("Failed to ping MongoDB", err)
	  return nil, err
	}
	
	return c, nil
  }