package translators

import (
	"fmt"
	"vnr/src/server/chrome"
)

var KnownTranslators = []string{
	"deepl",
}

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
	Chrome *chrome.Chrome
}

func GetTranslator(name string, options GetTranslatorOptions) (Translator, error) {
	switch name {
	case "deepl":
		return NewDeeplTranslator(options.Chrome), nil
	default:
		return nil, fmt.Errorf("invalid translator: got: %s, expected: one of %s", name, KnownTranslators)
	}
}

func GetAllKnownTranslators(options GetTranslatorOptions) (map[string]Translator, error) {
	res := map[string]Translator{}
	for _, name := range KnownTranslators {
		translator, err := GetTranslator(name, options)
		if err != nil {
			return nil, err
		}
		res[name] = translator
	}
	return res, nil
}
