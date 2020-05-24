// Package translation provides translators
package translation

// Translator is the interface that the translators must implement to be switchable in the translation process
type Translator interface {
	Translate(input TranslationDictionary) (TranslationDictionary, error)
}

// GetTranslator returns the current translator available
func GetTranslator() Translator {
	return &GoogleTranslator{}
}
