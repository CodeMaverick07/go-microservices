package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
    type requestPayload struct {
		Email   string `json:"email"`
		Password string `json:"password"`
	}

	var payload requestPayload
	err := app.readJSON(w, r, &payload)
	log.Println(payload)

	if err != nil {
		log.Println(err,1)
		app.errorJSON(w,err,http.StatusBadRequest)
		return 
	}

	user,err:= app.Models.User.GetByEmail(payload.Email)

	if err != nil {
		log.Println(err,2)
		app.errorJSON(w,errors.New("worng credentials"),http.StatusBadRequest)
		return 
	}

	 valid,err := user.PasswordMatches(payload.Password)

	 if err != nil || !valid {
		 log.Println(err,3)
		 app.errorJSON(w,errors.New("wrong credentials"),http.StatusBadRequest)
		 return
	 }
	 err = app.logRequest("authenticate",payload.Email)

	 if err != nil {
		 log.Println(err,4)
		 app.errorJSON(w,err,http.StatusBadRequest)
		 return
	 }

	 response := jsonResponse{
		Error: false,
		Message: "authenticated",
		Data: user,
	 }

	 app.writeJSON(w,http.StatusAccepted,response) 
	
}


func (app *Config) logRequest(name,data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData,_ := json.MarshalIndent(entry,"","  ")
	logServiceURL := "http://logger-service/log"

	request,err := http.NewRequest("POST",logServiceURL,bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	client := &http.Client{}
	_,err = client.Do(request)
	if err != nil {
		return err
	}
	return nil
}