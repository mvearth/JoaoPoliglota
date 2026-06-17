package translation

// TranslationDictionary is the request/response payload of the translate
// endpoint: a map of stable keys to text, translated from InputLang to
// OutputLang.
type TranslationDictionary struct {
	InputLang  string            `json:"inputLang"`
	OutputLang string            `json:"outputLang"`
	Dictionary map[string]string `json:"dictionary"`
	Translated bool              `json:"translated"`
}

// Translation is a single cached translation row.
type Translation struct {
	ID          int    `json:"id"`
	Idiom       string `json:"idiom"`
	StandardKey string `json:"standardKey"`
	Translation string `json:"translation"`
}
