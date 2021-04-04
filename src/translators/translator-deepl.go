package translators

import (
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/chromedp"
)

type DeeplTranslator struct {
	*BaseChromeTranslator
}

func NewDeeplTranslator() *DeeplTranslator {
	return &DeeplTranslator{
		BaseChromeTranslator: &BaseChromeTranslator{},
	}
}

func (dt *DeeplTranslator) GetTranslation(translationOptions TranslationOptions) (TranslationResult, error) {
	taskCtx, cancel := dt.getCtx()
	defer cancel()
	taskCtx, cancel = context.WithTimeout(taskCtx, translationOptions.Timeont)
	defer cancel()

	translationResult := TranslationResult{TranslationOptions: translationOptions}

	url := dt.translationOptionsToUrl(translationOptions)

	actions := []chromedp.Action{
		chromedp.Navigate(url),
		chromedp.WaitVisible("lmt__rating-up"),
		chromedp.TextContent("#target-dummydiv", &translationResult.Translation),
	}
	err := chromedp.Run(taskCtx, actions...)

	translationResult.Translation = strings.TrimSpace(translationResult.Translation)

	return translationResult, err
}

func (dt *DeeplTranslator) translationOptionsToUrl(translationOptions TranslationOptions) string {
	return fmt.Sprintf("https://www.deepl.com/translator#%s/%s/%s", translationOptions.From, translationOptions.To, translationOptions.Sentence)
}
