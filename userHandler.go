package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CP-Payne/chirpy/internal/auth"
	"github.com/CP-Payne/chirpy/internal/database"
)

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Password     string `json:"-"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string
	}

	type response struct {
		User
	}
	decoder := json.NewDecoder(r.Body)
	userParams := params{}
	err := decoder.Decode(&userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramters")
		return
	}
	userExist, err := cfg.DB.ValidateUserExist(userParams.Email)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if userExist {
		respondWithError(w, http.StatusConflict, "Email already exist")
		return
	}

	hashedPassword, err := auth.HashPassword(userParams.Password)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := cfg.DB.CreateUser(userParams.Email, string(hashedPassword))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		fmt.Println(err)
		return
	}
	respondWithJSON(w, http.StatusCreated, response{User{ID: user.ID, Email: user.Email, IsChirpyRed: user.IsChirpyRed}})
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	subject, issuer, err := auth.ValidateToken(token, cfg.jwtSecret)
	if err != nil || subject == "" || issuer == "" {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	if issuer == "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	decoder := json.NewDecoder(r.Body)
	userParams := params{}
	err = decoder.Decode(&userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramters")
		return
	}

	hashedPassword, err := auth.HashPassword(userParams.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	err = cfg.DB.UpdateUser(userIDInt, userParams.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			respondWithError(w, http.StatusConflict, "Email already exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user information")
		return
	}

	respondWithJSON(w, http.StatusOK, response{User: User{ID: userIDInt, Email: userParams.Email}})

}
