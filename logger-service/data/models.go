package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LogEntry struct {
	ID string `bson:"_id,omitempty" json:"id,omitempty"`
	Name string `bson:"name" json:"name"`
	Data string `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`


}

type Models struct {
	LogEntry LogEntry
}

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

func (l *LogEntry) Insert(entry LogEntry) error {
     collection := client.Database("logs").Collection("logs")
	 _, err := collection.InsertOne(context.TODO(), LogEntry{
		 Name: entry.Name,
		 Data: entry.Data,
		 CreatedAt: time.Now(),
		 UpdatedAt: time.Now(),
	 })

	 if err != nil {
		 panic(err)
		 return err 
	 }
	 return nil
}