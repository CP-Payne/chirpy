package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CP-Payne/chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
	polkaApiKey    string
}

func main() {
	mux := http.NewServeMux()
	godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	apiKey := os.Getenv("WEBHOOK_API")

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
		polkaApiKey:    apiKey,
	}

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	fileServer := http.FileServer(http.Dir("."))

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", http.HandlerFunc(cfg.handlerAddChirp))
	mux.HandleFunc("GET /api/chirps", http.HandlerFunc(cfg.handlerGetChirps))
	mux.HandleFunc("GET /api/chirps/{chirp_id}", http.HandlerFunc(cfg.handlerGetChirp))
	mux.HandleFunc("DELETE /api/chirps/{chirp_id}", http.HandlerFunc(cfg.handlerDeleteChirp))

	mux.HandleFunc("POST /api/users", http.HandlerFunc(cfg.handlerCreateUser))
	mux.HandleFunc("PUT /api/users", http.HandlerFunc(cfg.handlerUpdateUser))
	mux.HandleFunc("POST /api/login", http.HandlerFunc(cfg.handlerLoginUser))
	mux.HandleFunc("POST /api/refresh", http.HandlerFunc(cfg.handlerRefreshToken))
	mux.HandleFunc("POST /api/revoke", http.HandlerFunc(cfg.handlerRevokeToken))

	mux.HandleFunc("GET /admin/metrics", http.HandlerFunc(cfg.handlerMetrics))
	mux.HandleFunc("/api/reset", http.HandlerFunc(cfg.handlerReset))

	mux.HandleFunc("POST /api/polka/webhooks", http.HandlerFunc(cfg.handlerUpgradeUser))

	// mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets/"))))

	fmt.Printf("Listening on port %s\n", srv.Addr)
	err = srv.ListenAndServe()
	// err := http.ListenAndServe(":8000", corsMux)
	if err != nil {
		log.Fatal(err)
	}
}
