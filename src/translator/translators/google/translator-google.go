package google

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"vnr/src/translator"
)

type GoogleTranslator struct {
	langsCache translator.Langs
}

func NewGoogleTranslator() *GoogleTranslator {
	return &GoogleTranslator{}
}

func (dt *GoogleTranslator) GetTranslation(translationOptions translator.TranslationOptions) (translator.TranslationResult, error) {
	translationResult := translator.TranslationResult{TranslationOptions: translationOptions}

	rawUrl := dt.translationOptionsToUrl(translationOptions)

	r, err := http.Get(rawUrl)
	if err != nil {
		return translationResult, err
	}
	defer r.Body.Close()

	t, err := getTranslationFromGoogleBody(r.Body)
	if err != nil {
		return translationResult, err
	}
	translationResult.Translation = t

	return translationResult, nil
}

func (dt *GoogleTranslator) translationOptionsToUrl(translationOptions translator.TranslationOptions) string {
	q := url.Values{}
	q.Add("sl", translationOptions.From)
	q.Add("tl", translationOptions.To)
	q.Add("q", translationOptions.Sentence)
	return "https://translate.google.com/m?" + q.Encode()
}

func (dt *GoogleTranslator) GetLanguages() (translator.Langs, error) {
	if dt.langsCache != nil {
		return dt.langsCache, nil
	}

	r, err := http.Get("https://translate.google.com/m?mui=tl&hl=en")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	res, err := getLanguagesFromGoogleBody(r.Body)
	if err != nil {
		return nil, err
	}

	dt.langsCache = res
	return res, nil
}

func getTranslationFromGoogleBody(r io.Reader) (string, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`(?U)class="result-container">(.*)</div>`)
	match := re.FindAllSubmatch(body, -1)

	if len(match) == 0 {
		return "", fmt.Errorf("match not found")
	}
	res := html.UnescapeString(string(match[0][1]))

	return string(res), nil
}

func getLanguagesFromGoogleBody(r io.Reader) (translator.Langs, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	allLanguagesIndex := bytes.Index(body, []byte("All languages"))
	if allLanguagesIndex != -1 {
		body = body[allLanguagesIndex:]
	}

	re := regexp.MustCompile(`(?U)<a.*tl=(.*)&.*>(.*)</a>`)
	matches := re.FindAllSubmatch(body, -1)

	res := translator.Langs{}

	for _, match := range matches {
		lang := translator.Lang{Name: string(match[1]), Description: string(match[2])}
		res = append(res, lang)
	}

	res = res.RemoveDuplicates()

	return res, nil
}
