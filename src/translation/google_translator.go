package translation

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

const googleApiEnvKey string = "GOOGLE_APPLICATION_CREDENTIALS"

type GoogleTranslator struct {
}

func (gt GoogleTranslator) Translate(input TranslationDictionary) (TranslationDictionary, error) {
	checkAPIKey()

	translated := make(map[string]string)

	for key, value := range input.Dictionary {
		translation, err := GetTranslation(key, input.OutputLang)
		if err == nil && translation.StandardKey == key {
			translated[key] = translation.Translation
		} else {
			outputText, err := translateText(value, input.OutputLang)
			if err != nil {
				return input, err
			} else {
				translation := Translation{Idiom: input.OutputLang, StandardKey: key, Translation: outputText}
				done, err := InsertTranslation(translation)
				if err != nil && !done {
					fmt.Errorf("%s", err.Error())
				}
				translated[key] = outputText
			}
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
