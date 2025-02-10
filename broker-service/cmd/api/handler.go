package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
	Log LogPayload `json:"log,omitempty"`
	Mail MailPayload `json:"mail,omitempty"`
}
type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
type MailPayload struct {
	From string `json:"from"`
	To string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Broker(w http.ResponseWriter,r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Broker service is running",
	}
	_ = app.writeJSON(w,http.StatusOK,payload)

}

func (app *Config) HandleSubmission(w http.ResponseWriter,r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w,r,&requestPayload)

	if err != nil {
		app.errorJSON(w,err,http.StatusBadRequest)
		return
	}
 
	switch requestPayload.Action {
		case "auth":
			app.authenticate(w,requestPayload.Auth)
		
		case "log":
			app.logEventViaRabbit(w,requestPayload.Log)
		
	    case "mail":
			app.sendMail(w,requestPayload.Mail)
		
	    default:
			app.errorJSON(w,errors.New("invalid action"),http.StatusBadRequest)
	}
}

func (app *Config) logItem(w http.ResponseWriter,entry LogPayload ) {
	jsonData, _ := json.MarshalIndent(entry,"","  ")
    logServiceURL:= "http://logger-service/log"
	request ,err := http.NewRequest("POST",logServiceURL,bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w,err)
		return	
	}
	request.Header.Set("Content-Type","application/json")
	client := &http.Client{}
	response ,err := client.Do(request)
	if err != nil {
		app.errorJSON(w,err)
		return	
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w,errors.New("error calling log service"))
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "logged"
	payload.Data = entry

	app.writeJSON(w,http.StatusAccepted,payload)

}

func (app *Config) authenticate (w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a,"","  ")

	request, err := http.NewRequest("POST","http://authentication-service/authenticate",bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJSON(w,err)
		return
	}

	client := &http.Client{}
	response,err := client.Do(request)

	if err != nil {
		app.errorJSON(w,err)
		return
	}

	defer response.Body.Close()
	
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w,errors.New("invalid credentials"))
		return 
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w,errors.New("error calling auth service"))
		return 
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	 if err != nil {
		 app.errorJSON(w,err)
		 return
	 }
     if jsonFromService.Error {
		app.errorJSON(w,err,http.StatusUnauthorized )
	 }

	 var payload jsonResponse
	 payload.Error = false
	 payload.Message = "authenticated"
	 payload.Data = jsonFromService.Data 

	 app.writeJSON(w,http.StatusAccepted,payload)
	  
}

func (app *Config) sendMail(w http.ResponseWriter, m MailPayload) {
	jsonData,_ := json.MarshalIndent(m,"","  ")
	mailServiceURL := "http://mailer-service/send"
    request,err := http.NewRequest("POST",mailServiceURL,bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w,err)
		return
	}
	request.Header.Set("Content-Type","application/json")
	client := &http.Client{}
	response,err := client.Do(request)
	if err != nil {
		app.errorJSON(w,err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		log.Println(response.StatusCode)
		app.errorJSON(w,errors.New("error calling mail service"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "mail sent"
	payload.Data = m

	app.writeJSON(w,http.StatusAccepted,payload)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload ) {
	err := app.pushToQueue(l.Name,l.Data)
    if err != nil {
		app.errorJSON(w,err)
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged vai rabbitMQ" 
	app.writeJSON(w,http.StatusAccepted,payload)
}

func (app *Config) pushToQueue(name,msg string) error {
	emitter,err :=  event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}
	payload := LogPayload{
		Name: name,
		Data: msg,
	}
	log.Println(payload)
	j ,_:= json.MarshalIndent(payload,"","  ")
	err = emitter.Push(string(j),"log.INFO")
	log.Println("j: ",j)
	if err != nil {
		return err 
	}
	return nil
}