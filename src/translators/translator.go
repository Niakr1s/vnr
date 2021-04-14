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

type Lang struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Langs []Lang

func (l Langs) RemoveDuplicates() Langs {
	set := map[string]Lang{}

	res := Langs{}
	for _, lang := range l {
		if _, ok := set[lang.Name]; ok {
			continue
		}
		set[lang.Name] = lang
		res = append(res, lang)
	}
	return res
}

type Translator interface {
	GetTranslation(translationOptions TranslationOptions) (TranslationResult, error)
	GetLanguages() (Langs, error)
}

type GetTranslatorOptions struct {
	Chrome *chrome.Chrome
}

func GetTranslator(name string, options GetTranslatorOptions) (Translator, error) {
	switch name {
	case "deepl":
		return NewDeeplTranslator(options.Chrome), nil
	case "yandex":
		return NewYandexTranslator(options.Chrome), nil
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
