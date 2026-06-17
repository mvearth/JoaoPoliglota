// Package translation provides translators and their backing store.
package translation

import "context"

// Translator is the interface that translators must implement so the
// translation backend can be swapped without touching callers.
type Translator interface {
	Translate(ctx context.Context, input TranslationDictionary) (TranslationDictionary, error)
}
