package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CP-Payne/chirpy/internal/auth"
	"github.com/CP-Payne/chirpy/internal/database"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
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

	respondWithJSON(w, http.StatusOK, response{User{ID: user.ID, Email: user.Email}})

}
