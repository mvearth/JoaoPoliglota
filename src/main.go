package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"joao_poliglota/handler"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/translate", handler.Translate).Methods("POST")
	http.ListenAndServe(":8000", router)
}
