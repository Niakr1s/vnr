package deepl

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"vnr/src/chrome"
	"vnr/src/translator"

	"github.com/chromedp/chromedp"
)

type DeeplTranslator struct {
	chrome *chrome.Chrome

	langsCache translator.Langs
}

func NewDeeplTranslator(chrome *chrome.Chrome) *DeeplTranslator {
	return &DeeplTranslator{
		chrome: chrome,
	}
}

func (dt *DeeplTranslator) GetTranslation(translationOptions translator.TranslationOptions) (translator.TranslationResult, error) {
	if dt.chrome == nil {
		return dt.getTranslationWithoutChrome(translationOptions)
	}
	return dt.getTranslationWithChrome(translationOptions)
}

func (dt *DeeplTranslator) getTranslationWithChrome(translationOptions translator.TranslationOptions) (translator.TranslationResult, error) {
	translationResult := translator.TranslationResult{TranslationOptions: translationOptions}

	url := dt.translationOptionsToUrl(translationOptions)

	actions := []chromedp.Action{
		chromedp.Navigate(url),
		chromedp.WaitVisible("lmt__rating-up"),
		chromedp.TextContent("#target-dummydiv", &translationResult.Translation),
	}
	err := dt.chrome.Run(actions...)
	if err != nil {
		return translator.TranslationResult{}, err
	}

	translationResult.Translation = strings.TrimSpace(translationResult.Translation)
	return translationResult, nil
}

func (dt *DeeplTranslator) getTranslationWithoutChrome(translationOptions translator.TranslationOptions) (translator.TranslationResult, error) {
	translationResult := translator.TranslationResult{TranslationOptions: translationOptions}

	req, err := getDeeplTranslationRpcRequest(translationOptions)
	if err != nil {
		return translationResult, err
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return translationResult, err
	}
	defer r.Body.Close()
	if r.StatusCode == http.StatusTooManyRequests {
		return translationResult, fmt.Errorf("too many requests")
	}

	translation, err := getTranslationFromDeeplJsonRpcBody(r.Body)
	if err != nil {
		return translationResult, err
	}
	translationResult.Translation = translation
	return translationResult, nil
}

func getTranslationFromDeeplJsonRpcBody(r io.Reader) (string, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	type Rpc struct {
		Result struct {
			Translations []struct {
				Beams []struct {
					PostprocessedSentence string `json:"postprocessed_sentence"`
				} `json:"beams"`
			} `json:"translations"`
		} `json:"result"`
	}

	rpc := Rpc{}
	err = json.Unmarshal(body, &rpc)
	if err != nil {
		return "", err
	}
	if len(rpc.Result.Translations) == 0 || len(rpc.Result.Translations[0].Beams) == 0 {
		return "", fmt.Errorf("no translations")
	}

	return rpc.Result.Translations[0].Beams[0].PostprocessedSentence, nil
}

func (dt *DeeplTranslator) translationOptionsToUrl(translationOptions translator.TranslationOptions) string {
	return fmt.Sprintf("https://www.deepl.com/translator#%s/%s/%s", translationOptions.From, translationOptions.To, translationOptions.Sentence)
}

func (dt *DeeplTranslator) GetLanguages() (translator.Langs, error) {
	if dt.langsCache != nil {
		return dt.langsCache, nil
	}

	// dirty fix, too lazy to fix parser
	return translator.Langs{
		translator.Lang{
			Name:        "ja",
			Description: "Japanese",
		},
		translator.Lang{
			Name:        "en",
			Description: "English",
		},
		translator.Lang{
			Name:        "ru",
			Description: "Russian",
		},
	}, nil

	// r, err := http.DefaultClient.Get("https://www.deepl.com")
	// if err != nil {
	// 	return nil, err
	// }
	// defer r.Body.Close()

	// res, err := getLanguagesFromDeeplBody(r.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// dt.langsCache = res
	// return res, nil
}

func getLanguagesFromDeeplBody(r io.Reader) (translator.Langs, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`dl-lang='([a-z]{2}).*>(.*)</button`)
	matches := re.FindAllSubmatch(body, -1)
	_ = matches

	res := translator.Langs{}

	for _, match := range matches {
		lang := translator.Lang{Name: string(match[1]), Description: string(match[2])}
		// both English (American) and English (British) have same alias
		if strings.HasPrefix(lang.Description, "English") {
			lang.Description = "English"
		}
		res = append(res, lang)
	}

	res = res.RemoveDuplicates()

	return res, nil
}
