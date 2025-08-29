package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	type cleaned_body struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	bodysplit := strings.Split(params.Body, " ")

	for i := 0; i < len(bodysplit); i++ {
		body := strings.ToLower(bodysplit[i])
		if body == "kerfuffle" || body == "sharbert" || body == "fornax" {
			bodysplit[i] = "****"
		}
	}
	body := strings.Join(bodysplit, " ")
	cleanBody := &cleaned_body{}
	cleanBody.Cleaned_body = body
	respondWithJSON(w, 200, cleanBody)
}
