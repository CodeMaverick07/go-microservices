package main

import (
	"context"
	"log-service/data"
	"time"
)



type RPCServer struct {

}

type RPCPayload struct {
	Name string 
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload,resp *string) error {
 collection := client.Database("logger").Collection("logs")
_, err := collection.InsertOne(context.TODO(), data.LogEntry{
	Name:      payload.Name,
	Data:      payload.Data,
	CreatedAt: time.Now(),
})
 if err != nil {
	 return err
 }
 *resp = "porcessed payload via RPC" + payload.Name
 return nil
}
