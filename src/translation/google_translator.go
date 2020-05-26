package translation

import (
	"context"
	"os"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

const googleApiEnvKey string = "GOOGLE_APPLICATION_CREDENTIALS"

// GoogleTranslator encapsulate the translation process used by GoogleAPI
type GoogleTranslator struct {
}

// Translate using the GoogleAPI
func (gt GoogleTranslator) Translate(input TranslationDictionary) (TranslationDictionary, error) {
	checkAPIKey()

	translated := make(map[string]string)

	for key, value := range input.Dictionary {
		outputText, err := translateText(value, input.OutputLang)
		if err != nil {
			return input, err
		} else {
			translated[key] = outputText
		}
	}

	input.Dictionary = translated

	input.Translated = true

	return input, nil
}

func translateText(text, targetLanguage string) (string, error) {
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", err
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
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
