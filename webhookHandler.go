package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CP-Payne/chirpy/internal/auth"
	"github.com/CP-Payne/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error with apikey")
		return
	}
	if cfg.polkaApiKey != apiKey {
		respondWithError(w, http.StatusUnauthorized, "error with apikey")
		return
	}

	decoder := json.NewDecoder(r.Body)
	reqParams := params{}
	err = decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramters")
		return
	}
	if reqParams.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, http.StatusText(http.StatusOK))
		return
	}

	err = cfg.DB.UpgradeUser(reqParams.Data.UserID)

	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, "User upgraded")

}
