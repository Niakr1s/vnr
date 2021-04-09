package translators

import (
	"fmt"
	"vnr/src/server/chrome"
)

type TranslationOptions struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Sentence string `json:"sentence"`
}

func NewTranslationOptions(sentence string) TranslationOptions {
	return TranslationOptions{
		From:     "auto",
		To:       "auto",
		Sentence: sentence,
	}
}

type TranslationResult struct {
	TranslationOptions
	Translation string `json:"translation"`
}

type Translator interface {
	GetTranslation(translationOptions TranslationOptions) (TranslationResult, error)
}

type GetTranslatorOptions struct {
	Translator string
	Chrome     *chrome.Chrome
}

func GetTranslator(options GetTranslatorOptions) (Translator, error) {
	switch options.Translator {
	case "deepl":
		return NewDeeplTranslator(options.Chrome), nil
	default:
		return nil, fmt.Errorf("invalid translator: got: %s, expected: one of %s", options.Translator, KnownTranslators)
	}
}
