package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	Conn   *amqp.Connection
	QueueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer,error) {
	consumer := Consumer{
		Conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{},err
	}
	return consumer,nil
}

func (c *Consumer) setup() error {
	channel,err := c.Conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}


type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Consumer) Listen(topics []string) error {

	ch,err := c.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q,err := declateRandomQueue(ch)
	if err != nil {
		return err
	}

	for _,topic := range topics {
		err = ch.QueueBind(q.Name,topic,"logs_topic",false,nil)
		if err != nil {
			return err 
		}
	}

	messages,err := ch.Consume(q.Name,"",true,false,false,false,nil)
	if err != nil {
		return err 
	}

    forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body,&payload)
			go handlePayload(payload)
		}

	}()

	fmt.Printf("waiting for message [exchange, Queue] [logs_topic, %s]\n",q.Name)
    <-forever
  return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log","event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	default: 
	err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(payload Payload) error {
    // Create a map or struct that matches the logger service's expected format
    logData := map[string]string{
        "name": payload.Name,
        "data": payload.Data,
    }

    jsonData, err := json.Marshal(logData)
    if err != nil {
        return err
    }

    logServiceURL := "http://logger-service/log"
    request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    request.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return err 
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusAccepted {
        return fmt.Errorf("logger service returned status code %d", response.StatusCode)
    }

    return nil
}

