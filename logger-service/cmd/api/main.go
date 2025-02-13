package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"

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

type Config struct{
	Models data.Models
}

func main() {
  mongoClient, err := connectToMongo()

  if err != nil {
	panic(err)
  }
  log.Println("connected to mongo")
  client = mongoClient

  ctx , cancel := context.WithCancel(context.Background())
  defer cancel()

  defer func() {
	if err := client.Disconnect(ctx); err != nil {
	  panic(err)
	}
  }()

  app := Config{
	Models: data.New(client),

  }
  err = rpc.Register(new(RPCServer))
  if err != nil {
	panic(err)
	  }
  go app.rpcListen()
  go app.gRPCListen()
  log.Println("starting port on ",WebPort)
  srv := &http.Server{
	  Addr: fmt.Sprintf(":%d", WebPort),
	Handler: app.routes(),
}

err = srv.ListenAndServe()
if err != nil {
	panic(err)
}
}

func (app *Config) rpcListen() error {
	log.Println("starting rpc server on ",rpcPort)
	listen,err := net.Listen("tcp",fmt.Sprintf("0.0.0.0:%d",rpcPort))
    if err != nil {
		return err
	}
	defer listen.Close()
	for {
		rpcConn,err := listen.Accept()
		if err != nil {
			log.Println("error accepting connection",err)
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
		
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