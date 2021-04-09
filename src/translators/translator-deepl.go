package translators

import (
	"fmt"
	"strings"
	"vnr/src/server/chrome"

	"github.com/chromedp/chromedp"
)

type DeeplTranslator struct {
	chrome *chrome.Chrome
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
