package main

import (
	"fmt"
	"log"
	"net/http"

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
