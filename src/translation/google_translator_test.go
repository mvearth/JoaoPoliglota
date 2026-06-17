package translation

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type fakeBackend struct {
	resp   []translate.Translation
	err    error
	called int
}

func (f *fakeBackend) Translate(_ context.Context, _ []string, _ language.Tag, _ *translate.Options) ([]translate.Translation, error) {
	f.called++
	return f.resp, f.err
}

func (f *fakeBackend) Close() error { return nil }

type fakeStore struct {
	get       Translation
	getErr    error
	inserted  []Translation
	insertErr error
}

func (f *fakeStore) Get(_ context.Context, _, _ string) (Translation, error) {
	return f.get, f.getErr
}

func (f *fakeStore) Insert(_ context.Context, t Translation) (bool, error) {
	if f.insertErr != nil {
		return false, f.insertErr
	}
	f.inserted = append(f.inserted, t)
	return true, nil
}

func TestTranslateCacheHit(t *testing.T) {
	store := &fakeStore{get: Translation{StandardKey: "hello", Translation: "olá"}}
	backend := &fakeBackend{}
	gt := &GoogleTranslator{backend: backend, store: store}

	out, err := gt.Translate(context.Background(), TranslationDictionary{
		OutputLang: "pt",
		Dictionary: map[string]string{"hello": "hello"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Dictionary["hello"] != "olá" {
		t.Errorf("got %q, want olá", out.Dictionary["hello"])
	}
	if backend.called != 0 {
		t.Errorf("backend called %d times on cache hit, want 0", backend.called)
	}
	if !out.Translated {
		t.Error("Translated flag should be set")
	}
}

func TestTranslateCacheMissThenStores(t *testing.T) {
	store := &fakeStore{} // Get returns zero -> miss
	backend := &fakeBackend{resp: []translate.Translation{{Text: "olá"}}}
	gt := &GoogleTranslator{backend: backend, store: store}

	out, err := gt.Translate(context.Background(), TranslationDictionary{
		OutputLang: "pt",
		Dictionary: map[string]string{"hello": "hello"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Dictionary["hello"] != "olá" {
		t.Errorf("got %q, want olá", out.Dictionary["hello"])
	}
	if backend.called != 1 {
		t.Errorf("backend called %d times, want 1", backend.called)
	}
	if len(store.inserted) != 1 || store.inserted[0].Translation != "olá" {
		t.Errorf("expected the translation to be cached, got %+v", store.inserted)
	}
}

func TestTranslateBackendErrorPropagates(t *testing.T) {
	store := &fakeStore{}
	backend := &fakeBackend{err: errors.New("api down")}
	gt := &GoogleTranslator{backend: backend, store: store}

	_, err := gt.Translate(context.Background(), TranslationDictionary{
		OutputLang: "pt",
		Dictionary: map[string]string{"hello": "hello"},
	})
	if err == nil {
		t.Fatal("expected an error when the backend fails")
	}
}

func TestTranslateInsertFailureIsNonFatal(t *testing.T) {
	store := &fakeStore{insertErr: errors.New("write failed")}
	backend := &fakeBackend{resp: []translate.Translation{{Text: "olá"}}}
	gt := &GoogleTranslator{backend: backend, store: store}

	out, err := gt.Translate(context.Background(), TranslationDictionary{
		OutputLang: "pt",
		Dictionary: map[string]string{"hello": "hello"},
	})
	if err != nil {
		t.Fatalf("cache write failure should not fail the request: %v", err)
	}
	if out.Dictionary["hello"] != "olá" {
		t.Errorf("got %q, want olá", out.Dictionary["hello"])
	}
}

func TestTranslateInvalidLanguage(t *testing.T) {
	gt := &GoogleTranslator{backend: &fakeBackend{}, store: &fakeStore{}}
	_, err := gt.Translate(context.Background(), TranslationDictionary{
		OutputLang: "not-a-language!!",
		Dictionary: map[string]string{"hello": "hello"},
	})
	if err == nil {
		t.Fatal("expected an error for an invalid target language")
	}
}
