package translation

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/text/language"
)

const googleApiEnvKey string = "GOOGLE_APPLICATION_CREDENTIALS"

// GoogleTranslator encapsulate the translation process used by GoogleAPI
type GoogleTranslator struct {
}

// Translate using the GoogleAPI
func (gt GoogleTranslator) Translate(input TranslationDictionary) (TranslationDictionary, error) {
	checkAPIKey()

	var translated map[string]string

	for index, element := range input.Dictionary {
		outputText, err := translate(input.InputLang, input.OutputLang, element)
		if err != nil {

		}

		append(translated[index], element)
	}

	input.Dictionary = translated

	return input, nil
}

func translate(text, inputLang, outputLang string) (string, error) {
	ctx := context.Background()

	lang, err := language.Parse("pt-br")
	if err != nil {
		fmt.Fprintf(w, "language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{"Test"}, lang, nil)
	if err != nil {
		return "", err
	}

	return resp[0].Text, nil
}

func checkAPIKey() {
	if os.Getenv(googleApiEnvKey) == "" {
		os.Setenv(googleApiEnvKey, "C:\\Users\\mathv\\Sources\\JoaoPoliglota\\JoaoPoliglota-1b5ba2f03485.json")
	}
}
