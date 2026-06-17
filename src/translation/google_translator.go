package translation

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// translateBackend is the subset of *translate.Client used by GoogleTranslator,
// extracted as an interface so the backend can be faked in tests.
type translateBackend interface {
	Translate(ctx context.Context, inputs []string, target language.Tag, opts *translate.Options) ([]translate.Translation, error)
	Close() error
}

// translationStore is the persistence behaviour GoogleTranslator depends on.
type translationStore interface {
	Get(ctx context.Context, standardKey, idiom string) (Translation, error)
	Insert(ctx context.Context, t Translation) (bool, error)
}

// GoogleTranslator translates text via the Google Cloud Translation API,
// caching results through a translationStore.
type GoogleTranslator struct {
	backend translateBackend
	store   translationStore
}

// NewGoogleTranslator creates a translator backed by a Google Translate client.
// Credentials are resolved by the client from the environment (e.g.
// GOOGLE_APPLICATION_CREDENTIALS or Application Default Credentials). The
// returned translator owns the client and must be closed with Close.
func NewGoogleTranslator(ctx context.Context, store translationStore) (*GoogleTranslator, error) {
	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating translate client: %w", err)
	}
	return &GoogleTranslator{backend: client, store: store}, nil
}

// Close releases the underlying translate client.
func (gt *GoogleTranslator) Close() error {
	return gt.backend.Close()
}

// Translate fills in the dictionary, serving cached entries when available and
// caching newly translated ones. A cache failure is logged but does not fail
// the request.
func (gt *GoogleTranslator) Translate(ctx context.Context, input TranslationDictionary) (TranslationDictionary, error) {
	translated := make(map[string]string, len(input.Dictionary))

	for key, value := range input.Dictionary {
		cached, err := gt.store.Get(ctx, key, input.OutputLang)
		if err != nil {
			log.Printf("translation: cache lookup for %q failed: %v", key, err)
		} else if cached.StandardKey == key {
			translated[key] = cached.Translation
			continue
		}

		outputText, err := gt.translateText(ctx, value, input.OutputLang)
		if err != nil {
			return input, err
		}

		t := Translation{Idiom: input.OutputLang, StandardKey: key, Translation: outputText}
		if _, err := gt.store.Insert(ctx, t); err != nil {
			log.Printf("translation: caching %q failed: %v", key, err)
		}
		translated[key] = outputText
	}

	input.Dictionary = translated
	input.Translated = true

	return input, nil
}

func (gt *GoogleTranslator) translateText(ctx context.Context, text, targetLanguage string) (string, error) {
	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("invalid target language %q: %w", targetLanguage, err)
	}

	resp, err := gt.backend.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", err
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("no translation returned for %q", text)
	}

	return resp[0].Text, nil
}
