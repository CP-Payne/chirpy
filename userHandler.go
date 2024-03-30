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
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Token    string `json:"token,omitempty"`
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
	respondWithJSON(w, http.StatusCreated, response{User{ID: user.ID, Email: user.Email}})
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int64 `json:"expires_in_seconds,omitempty"`
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
	user := database.User{}
	if userExist {
		user, err = cfg.DB.GetUserByEmail(userParams.Email)
		if err != nil {
			fmt.Println(err)
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

	}
	err = auth.CheckPasswordHash(userParams.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect credentials")
		return
	}

	// Using int64 for if we want to increase the 24 hour expire limit in the future
	expiresInSeconds := int64(0)
	if userParams.ExpiresInSeconds != nil {
		expiresInSeconds = *userParams.ExpiresInSeconds
	}
	const maxExpirationInSeconds = 24 * 60 * 60
	if expiresInSeconds <= 0 || expiresInSeconds > maxExpirationInSeconds {
		expiresInSeconds = maxExpirationInSeconds
	}

	// If expiration time is set to 0 or greater than 24 hours, it will default to 24 hours
	token, err := auth.CreateToken(user.ID, cfg.jwtSecret, expiresInSeconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, response{User{ID: user.ID, Email: user.Email, Token: token}})

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

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
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
