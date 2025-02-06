package main

import (
	"errors"
	"log-service/data"
	"net/http"
)
type JSONPaload struct {
	Name string `json:"name"`
	Data string `json:"data"`

}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var 
	var requestPayload JSONPaload
    _ = app.readJSON(w,r,&requestPayload)

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
    err := app.Models.LogEntry.Insert(event)
    if err != nil {
		app.errorJSON(w,errors.New("error in inserting logentry"))
		return
	}

	resp := jsonResponse{
		Error: false,
		Message: "inserted log entry",
		
	}

	app.writeJSON(w,http.StatusAccepted,resp)

	

}