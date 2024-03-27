package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/CP-Payne/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	mux := http.NewServeMux()

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	fileServer := http.FileServer(http.Dir("."))

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("/api/reset", http.HandlerFunc(cfg.handlerReset))
	mux.HandleFunc("POST /api/chirps", http.HandlerFunc(cfg.handlerAddChirp))
	mux.HandleFunc("GET /api/chirps", http.HandlerFunc(cfg.handlerGetChirps))
	mux.HandleFunc("GET /api/chirps/{chirp_id}", http.HandlerFunc(cfg.handlerGetChirp))

	mux.HandleFunc("GET /admin/metrics", http.HandlerFunc(cfg.handlerMetrics))

	// mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets/"))))

	fmt.Printf("Listening on port %s\n", srv.Addr)
	err = srv.ListenAndServe()
	// err := http.ListenAndServe(":8000", corsMux)
	if err != nil {
		log.Fatal(err)
	}
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

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
	respChirps, err := cfg.DB.GetChirps()
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
	}

	respondWithJSON(w, http.StatusOK, respChirps)

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

	decoder := json.NewDecoder(r.Body)
	chirpText := params{}
	err := decoder.Decode(&chirpText)
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

	chirp, err := cfg.DB.CreateChirp(cleandedChirp)
	if err != nil {
		fmt.Println(err)
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Printf("Responing with 5xx error: %s", message)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
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
