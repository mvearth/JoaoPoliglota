package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/translate"
	"github.com/gorilla/mux"
	"golang.org/x/text/language"
)

func main() {
	router := mux.NewRouter()
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "KEYPATH")
	router.HandleFunc("/test", test)
	http.ListenAndServe(":8000", router)
}

func test(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	lang, err := language.Parse("pt-br")
	if err != nil {
		fmt.Fprintf(w, "language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		fmt.Fprintf(w, "translate.NewClient: %v", err)
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{"Test"}, lang, nil)
	if err != nil {
		fmt.Fprintf(w, "client.Translate: %v", err)
	}
	if len(resp) == 0 {
		fmt.Fprintf(w, "Translate returned empty response to text: Test")
	}
	fmt.Fprintf(w, resp[0].Text)
}
