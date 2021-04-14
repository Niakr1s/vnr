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

// YandexTranslator won't work because of anti-bot system. :(
type YandexTranslator struct {
	chrome *chrome.Chrome

	langsCache Langs
}

func NewYandexTranslator(chrome *chrome.Chrome) *YandexTranslator {
	return &YandexTranslator{
		chrome: chrome,
	}
}

func (dt *YandexTranslator) GetTranslation(translationOptions TranslationOptions) (TranslationResult, error) {
	translationResult := TranslationResult{TranslationOptions: translationOptions}

	url := dt.translationOptionsToUrl(translationOptions)

	actions := []chromedp.Action{
		chromedp.Navigate(url),
		chromedp.WaitVisible("#externalContent"),
		chromedp.TextContent(".translation-word", &translationResult.Translation),
	}
	err := dt.chrome.Run(actions...)
	if err != nil {
		return TranslationResult{}, err
	}

	translationResult.Translation = strings.TrimSpace(translationResult.Translation)
	return translationResult, nil
}

func (dt *YandexTranslator) translationOptionsToUrl(translationOptions TranslationOptions) string {
	return fmt.Sprintf("https://translate.yandex.com/?lang=%s-%s&text=%s", translationOptions.From, translationOptions.To, translationOptions.Sentence)
}

func (dt *YandexTranslator) GetLanguages() (Langs, error) {
	if dt.langsCache != nil {
		return dt.langsCache, nil
	}

	r, err := http.DefaultClient.Get("https://translate.yandex.com")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	res, err := getLanguagesFromYandexBody(r.Body)
	if err != nil {
		return nil, err
	}

	dt.langsCache = res
	return res, nil
}

func getLanguagesFromYandexBody(r io.Reader) (Langs, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`TRANSLATOR_LANGS:\s*{(.*)}`)
	matches := re.FindAllSubmatch(body, 1)
	if len(matches) == 0 || len(matches[0]) < 2 {
		return nil, fmt.Errorf("match not found")
	}

	langStr := string(matches[0][1])

	res, err := parseYandexLangString(langStr)
	if err != nil {
		return nil, err
	}
	res = res.RemoveDuplicates()

	return res, nil
}

// parseYandexLangString parses str of kind  ""az":"Азербайджанский","sq":"Албанский",..." to Langs
func parseYandexLangString(str string) (Langs, error) {
	re := regexp.MustCompile(`"(.*)":"(.*)"`)

	str = strings.TrimSpace(str)
	splitted := strings.Split(str, ",")

	res := Langs{}
	for _, langStr := range splitted {
		match := re.FindAllStringSubmatch(langStr, -1)
		if len(match) == 0 || len(match[0]) != 3 {
			return nil, fmt.Errorf("couldn't parse string %s", langStr)
		}
		lang := Lang{Name: match[0][1], Description: match[0][2]}
		res = append(res, lang)
	}
	return res, nil
}
