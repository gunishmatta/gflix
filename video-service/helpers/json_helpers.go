package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ResponseStrategy interface {
	Write(w http.ResponseWriter, payload interface{}) error
}

type JSONStrategy struct {
}

type PlainTextStrategy struct{}

func (j JSONStrategy) Write(w http.ResponseWriter, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		return err
	}
	return nil
}

func (s PlainTextStrategy) Write(w http.ResponseWriter, payload interface{}) error {
	_, err := fmt.Fprint(w, payload)
	if err != nil {
		return err
	}
	return nil
}

func Response(w http.ResponseWriter, statusCode int, payload interface{}, strategy ResponseStrategy) {
	w.WriteHeader(statusCode)
	err := strategy.Write(w, payload)
	if err != nil {
		log.Println("Failed to serialize response:", err)
	}
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	Response(w, statusCode, payload, JSONStrategy{})
}

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	Response(w, statusCode, message, PlainTextStrategy{})
}
