package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CP-Payne/chirpy/internal/auth"
	"github.com/CP-Payne/chirpy/internal/database"
)

// Using int64 for if we want to increase the 24 hour expire limit in the future
// Expiration in seconds
const accessTokenExpiration = 60 * 60
const refreshTokenExpiration = 60 * 24 * 60 * 60

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	// If expiration time is set to 0 or greater than 24 hours, it will default to 24 hours
	accessToken, err := auth.CreateToken(user.ID, cfg.jwtSecret, accessTokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	refreshToken, err := auth.CreateRefreshToken(user.ID, cfg.jwtSecret, refreshTokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	err = cfg.DB.AddTokenToDB(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, response{User{ID: user.ID, Email: user.Email, IsChirpyRed: user.IsChirpyRed, Token: accessToken, RefreshToken: refreshToken}})

}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
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

	if issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Provide a refresh token")
		return
	}
	revokedStatus, err := cfg.DB.IsTokenRevoked(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")

		return
	}

	if revokedStatus {
		respondWithError(w, http.StatusUnauthorized, "please log in again")
		return
	}

	userIdInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID from token")
	}

	accessToken, err := auth.CreateToken(userIdInt, cfg.jwtSecret, accessTokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, response{Token: accessToken})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Success string `json:"success"`
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

	if issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Provide a refresh token")
		return
	}
	revokedStatus, err := cfg.DB.IsTokenRevoked(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if revokedStatus {
		respondWithError(w, http.StatusConflict, "Token already revoked")
		return
	}

	err = cfg.DB.RevokeRefreshToken(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
	}
	respondWithJSON(w, http.StatusOK, response{Success: "Token revoked"})
}
