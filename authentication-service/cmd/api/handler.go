package main

import (
	"errors"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
    type requestPayload struct {
		Email   string `json:"email"`
		Password string `json:"password"`
	}

	var payload requestPayload
	err := app.readJSON(w, r, &payload)

	if err != nil {
		app.errorJSON(w,err,http.StatusBadRequest)
		return 
	}

	user,err:= app.Models.User.GetByEmail(payload.Email)

	if err != nil {
		app.errorJSON(w,errors.New("worng credentials"),http.StatusBadRequest)
		return 
	}

	 valid,err := user.PasswordMatches(payload.Password)

	 if err != nil || !valid {
		 app.errorJSON(w,errors.New("wrong credentials"),http.StatusBadRequest)
		 return
	 }

	 response := jsonResponse{
		Error: false,
		Message: "authenticated",
		Data: user,
	 }

	 app.writeJSON(w,http.StatusAccepted,response) 
	
}