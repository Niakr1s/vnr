package translators

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"vnr/src/server/chrome"
	"vnr/src/util"

	"github.com/chromedp/chromedp"
)

type DeeplTranslator struct {
	chrome *chrome.Chrome

	langsCache []string
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

func (dt *DeeplTranslator) GetLanguages() ([]string, error) {
	if dt.langsCache != nil {
		return dt.langsCache, nil
	}

	r, err := http.DefaultClient.Get("https://www.deepl.com")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`dl-lang='([a-z]{2})`)
	matches := re.FindAllSubmatch(body, -1)
	_ = matches

	res := []string{}

	for _, match := range matches {
		res = append(res, string(match[1]))
	}

	res = util.RemoveDuplicates(res)
	dt.langsCache = res
	return res, nil
}
