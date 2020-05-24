package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"joao_poliglota/handler"
)

func main() {
	router := mux.NewRouter()
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "C:\\Users\\mathv\\Sources\\JoaoPoliglota\\JoaoPoliglota-1b5ba2f03485.json")
	router.HandleFunc("/translate", handler.Translate).Methods("POST")
	http.ListenAndServe(":8000", router)
}
