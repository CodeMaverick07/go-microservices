package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
	Log LogPayload `json:"log,omitempty"`
}

type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
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
			app.logItem(w,requestPayload.Log)
		
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