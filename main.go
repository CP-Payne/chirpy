package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	corsMux := middlewareCors(mux)

	srv := http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))

	fmt.Printf("Listening on port %s\n", srv.Addr)
	err := srv.ListenAndServe()
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
