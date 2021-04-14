package translators

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"vnr/src/server/chrome"

	"github.com/chromedp/chromedp"
)

type DeeplTranslator struct {
	chrome *chrome.Chrome

	langsCache Langs
}

func NewDeeplTranslator(chrome *chrome.Chrome) *DeeplTranslator {
	return &DeeplTranslator{
		chrome: chrome,
	}
}

func (dt *DeeplTranslator) GetTranslation(translationOptions TranslationOptions) (TranslationResult, error) {
	translationResult := TranslationResult{TranslationOptions: translationOptions}

	url := dt.translationOptionsToUrl(translationOptions)

	actions := []chromedp.Action{
		chromedp.Navigate(url),
		chromedp.WaitVisible("lmt__rating-up"),
		chromedp.TextContent("#target-dummydiv", &translationResult.Translation),
	}
	err := dt.chrome.Run(actions...)
	if err != nil {
		return TranslationResult{}, err
	}

	translationResult.Translation = strings.TrimSpace(translationResult.Translation)
	return translationResult, nil
}

func (dt *DeeplTranslator) translationOptionsToUrl(translationOptions TranslationOptions) string {
	return fmt.Sprintf("https://www.deepl.com/translator#%s/%s/%s", translationOptions.From, translationOptions.To, translationOptions.Sentence)
}

func (dt *DeeplTranslator) GetLanguages() (Langs, error) {
	if dt.langsCache != nil {
		return dt.langsCache, nil
	}

	r, err := http.DefaultClient.Get("https://www.deepl.com")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	res, err := getLanguagesFromDeeplBody(r.Body)
	if err != nil {
		return nil, err
	}

	dt.langsCache = res
	return res, nil
}

func getLanguagesFromDeeplBody(r io.Reader) (Langs, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`dl-lang='([a-z]{2}).*>(.*)</button`)
	matches := re.FindAllSubmatch(body, -1)
	_ = matches

	res := Langs{}

	for _, match := range matches {
		lang := Lang{Name: string(match[1]), Description: string(match[2])}
		// both English (American) and English (British) have same alias
		if strings.HasPrefix(lang.Description, "English") {
			lang.Description = "English"
		}
		res = append(res, lang)
	}

	res = res.RemoveDuplicates()

	return res, nil
}
