package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data",omitempty`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
    // Limit request body size to 1MB
    maxBytes := 1048576 
    r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

    // Create a JSON decoder
    dec := json.NewDecoder(r.Body)

    // Decode the JSON into the provided data structure
    err := dec.Decode(data)
    if err != nil {
        return err
    }

    // Ensure only one JSON value is in the request body
    err = dec.Decode(&struct{}{})
    if err != io.EOF {
        return errors.New("body must only have a single JSON value")
    }


// Valid request body:
// { "name": "John" }

// Invalid request body:
// { "name": "John" } { "name": "Jane" }

    return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
    // Convert data to JSON
    out, err := json.Marshal(data)
    if err != nil {
        return err
    } 

    // Add optional headers
    if len(headers) > 0 {
        for key, value := range headers[0] {
            w.Header()[key] = value
        }
    }

    // Set JSON content type and status
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    // Write JSON response
    _, err = w.Write(out)
    if err != nil {
        return err
    }

    return nil
}

// if we want to return an error response in JSON format, we can use the errorJSON method.

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
    // Default to 400 Bad Request, or use provided status
    statuscode := http.StatusBadRequest
    if len(status) > 0 {
        statuscode = status[0]
    }

    // Create error response payload
    var payload jsonResponse
    payload.Error = true
    payload.Message = err.Error()

	//err.Error() is a method in Go that converts an error into a string representation of the error message.

    // Write JSON error response
    return app.writeJSON(w, statuscode, payload)
}
