// Package tete defines the models used in the translation process
package translation

// TranslationDictionary struct used to input and output in the translation process
type TranslationDictionary struct {
	InputLang, OutputLang string
	Dictionary            map[string]string
	Translated            bool
}
