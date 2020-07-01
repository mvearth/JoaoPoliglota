package translation

type TranslationDictionary struct {
	InputLang, OutputLang string
	Dictionary            map[string]string
	Translated            bool
}

type Translation struct {
	ID                              int
	Idiom, StandardKey, Translation string
}
