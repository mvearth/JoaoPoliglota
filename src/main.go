package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"joao_poliglota/config"
	"joao_poliglota/handler"
	"joao_poliglota/translation"
)

func main() {
	db, err := translation.Connect(config.LoadDB())
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		// Non-fatal: the pool connects lazily, so a transient outage at
		// startup should not prevent the server from booting.
		log.Printf("database: ping failed: %v", err)
	} else {
		log.Println("database connected")
	}

	ctx := context.Background()
	translator, err := translation.NewGoogleTranslator(ctx, translation.NewRepository(db))
	if err != nil {
		log.Fatalf("translator: %v", err)
	}
	defer translator.Close()

	router := mux.NewRouter()
	router.HandleFunc("/translate", handler.New(translator).Translate).Methods("POST")

	srv := &http.Server{
		Addr:         config.HTTPAddr(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
