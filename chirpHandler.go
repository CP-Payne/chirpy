package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/CP-Payne/chirpy/internal/auth"
	"github.com/CP-Payne/chirpy/internal/database"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodGet {
	// 	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodGet {
	// 	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	// 	return
	// }
	// TODO
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// hits := fmt.Sprintf("Hits: %d", cfg.fileserverHits)
	fmt.Fprintf(w, `<html>
	  <body>
	    <h1>Welcome, Chirpy Admin</h1>
	    <p>Chirpy has been visited %d times!</p>
	  </body>
	</html>`, cfg.fileserverHits)

	// w.Write([]byte(hits))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	// TODO
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	// Optional parameter
	authorIDParam, ok := r.URL.Query()["author_id"]
	var chirps []database.Chirp
	var err error
	if ok && len(authorIDParam[0]) > 0 {
		// Convert author_id from string to int
		authorID, err := strconv.Atoi(authorIDParam[0])
		if err != nil {
			http.Error(w, "Invalid author_id parameter", http.StatusBadRequest)
			return
		}
		// Fetch chirps for the given author ID
		chirps, err = cfg.DB.GetChirpsByAuthorID(authorID)
	} else {
		// Fetch all chirps if no author_id is provided
		chirps, err = cfg.DB.GetChirps()
	}
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
	}

	respondWithJSON(w, http.StatusOK, chirps)

}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirp_id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id")
		return
	}

	chirp, err := cfg.DB.GetChirp(intId)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if chirp.ID == 0 {
		respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return

	}

	respondWithJSON(w, http.StatusOK, chirp)

}

func (cfg *apiConfig) handlerAddChirp(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body string `json:"body"`
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

	if issuer != "chirpy-access" {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	userIdInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID from token")
	}

	decoder := json.NewDecoder(r.Body)
	chirpText := params{}
	err = decoder.Decode(&chirpText)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramters")
		return
	}

	if len(chirpText.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleandedChirp := cleanChirp(chirpText.Body, badWords)

	chirp, err := cfg.DB.CreateChirp(cleandedChirp, userIdInt)
	if err != nil {
		fmt.Println(err)
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirp_id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id")
		return
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

	if issuer != "chirpy-access" {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	chirp, err := cfg.DB.GetChirp(intId)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if chirp.ID == 0 {
		respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	// Compare chirp author with subject
	authorID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}
	if authorID != chirp.AuthorID {
		respondWithError(w, http.StatusForbidden, "You do not own this resource")
		return
	}
	err = cfg.DB.DeleteChirp(intId)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, database.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	respondWithJSON(w, http.StatusOK, http.StatusText(http.StatusOK))
}

func cleanChirp(chirpText string, badWords map[string]struct{}) (cleanedChirp string) {
	textSlice := strings.Split(chirpText, " ")
	for i, word := range textSlice {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			textSlice[i] = "****"
		}
	}
	return strings.Join(textSlice, " ")
}
