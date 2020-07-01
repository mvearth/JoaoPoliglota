package handler

import (
	"encoding/json"
	"joao_poliglota/translation"
	"log"
	"net/http"
)

func Translate(w http.ResponseWriter, r *http.Request) {
	translator := translation.GetTranslator()
	decoder := json.NewDecoder(r.Body)

	var trInput translation.TranslationDictionary
	err := decoder.Decode(&trInput)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(err)
	} else {
		w.Header().Set("Content-type", "application/json")

		trOutput, err := translator.Translate(trInput)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode(trOutput)
		}
	}
}
