package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	userParams := params{}
	err := decoder.Decode(&userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramters")
		return
	}

	user, err := cfg.DB.CreateUser(userParams.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		fmt.Println(err)
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}
