package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"joao_poliglota/translation"
)

type fakeTranslator struct {
	out    translation.TranslationDictionary
	err    error
	called bool
}

func (f *fakeTranslator) Translate(_ context.Context, in translation.TranslationDictionary) (translation.TranslationDictionary, error) {
	f.called = true
	if f.err != nil {
		return in, f.err
	}
	return f.out, nil
}

func post(t *testing.T, h *Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/translate", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.Translate(rec, req)
	return rec
}

func TestTranslateOK(t *testing.T) {
	ft := &fakeTranslator{out: translation.TranslationDictionary{
		OutputLang: "pt",
		Dictionary: map[string]string{"hello": "olá"},
		Translated: true,
	}}
	h := New(ft)

	rec := post(t, h, `{"outputLang":"pt","dictionary":{"hello":"hello"}}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("content-type = %q, want application/json", ct)
	}
	var got translation.TranslationDictionary
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if got.Dictionary["hello"] != "olá" {
		t.Errorf("dictionary[hello] = %q, want olá", got.Dictionary["hello"])
	}
}

func TestTranslateBadRequest(t *testing.T) {
	cases := map[string]string{
		"invalid json":     `{not json`,
		"unknown field":    `{"outputLang":"pt","dictionary":{"a":"b"},"bogus":1}`,
		"missing lang":     `{"dictionary":{"a":"b"}}`,
		"empty dictionary": `{"outputLang":"pt","dictionary":{}}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			ft := &fakeTranslator{}
			rec := post(t, New(ft), body)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
			}
			if ft.called {
				t.Error("translator should not be called on a bad request")
			}
		})
	}
}

func TestTranslateInternalError(t *testing.T) {
	ft := &fakeTranslator{err: context.DeadlineExceeded}
	rec := post(t, New(ft), `{"outputLang":"pt","dictionary":{"a":"b"}}`)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
