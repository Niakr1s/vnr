package google

import (
	"strings"
	"testing"
	"vnr/src/translator"

	"github.com/stretchr/testify/assert"
)

func TestGoogleTranslator(t *testing.T) {
	google := NewGoogleTranslator()

	t.Run("GetLanguages", func(t *testing.T) {
		langs, err := google.GetLanguages()
		assert.Nil(t, err)
		assert.NotZero(t, len(langs))
	})

	t.Run("GetTranslation", func(t *testing.T) {
		translationOptions := translator.TranslationOptions{From: "ja", To: "ru", Sentence: "こんにちは世界"}
		res, err := google.GetTranslation(translationOptions)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "Привет мир", res.Translation)
	})
}

func Test_getTranslationFromGoogleBody(t *testing.T) {
	body := `</style></head><body><div class="root-container"><div class="header"><div class="logo-image"></div><div class="logo-text">Перевести</div></div><div class="languages-container"><div class="sl-and-tl"><a href="./m?sl=en&amp;tl=ru&amp;q=hello&amp;mui=sl&amp;hl=ru">английский</a> → <a href="./m?sl=en&amp;tl=ru&amp;q=hello&amp;mui=tl&amp;hl=ru">русский</a></div><div class="swap-container"><a href="./m?sl=ru&amp;tl=en&amp;q=%D0%9F%D1%80%D0%B8%D0%B2%D0%B5%D1%82&amp;hl=ru">Переключить</a></div></div><div class="input-container"><form action="/m"><input type="hidden" name="sl" value="en"><input type="hidden" name="tl" value="ru"><input type="hidden" name="hl" value="ru"><input type="text" aria-label="Исходный текст" name="q" class="input-field" maxlength="2048" value="hello"><div class="translate-button-container"><input type="submit" value="Перевести" class="translate-button"></div></form></div><div class="result-container">Привет</div><div class="links-container"><ul><li><a href="https://www.google.com/m?hl=ru">Главная страница Google</a></li><li><a href="https://www.google.com/tools/feedback/survey/xhtml?productId=95112&hl=ru">Отправить отзыв</a></li><li><a href="https://www.google.com/intl/ru/policies">Политика конфиденциальности и Условия использования</a></li><li><a href="./full">Полная версия</a></li></ul></div></div></body></html>`
	bodyReader := strings.NewReader(body)

	translation, err := getTranslationFromGoogleBody(bodyReader)

	assert.NoError(t, err)
	assert.Equal(t, "Привет", translation)
}

func Test_getLanguagesFromGoogleBody(t *testing.T) {
	body := `
            <div class="language-list-header">Recent languages</div>
            <div class="language-item">
                <a href="./m?sl&amp;tl=ru&amp;hl=en">Russian</a>
            </div>
            <div class="language-item">
                <a href="./m?sl&amp;tl=en&amp;hl=en">English</a>
            </div>
            <div class="language-list-header">All languages</div>
            <div class="language-item">
                <a href="./m?sl&amp;tl=af&amp;hl=en">Afrikaans</a>
            </div>
            <div class="language-item">
                <a href="./m?sl&amp;tl=sq&amp;hl=en">Albanian</a>
            </div>
            <div class="language-item">
                <a href="./m?sl&amp;tl=zh-CN&amp;hl=en">Chinese (Simplified)</a>
            </div>
            <div class="language-item">
                <a href="./m?sl&amp;tl=zh-TW&amp;hl=en">Chinese (Traditional)</a>
            </div>
    `
	bodyReader := strings.NewReader(body)

	langs, err := getLanguagesFromGoogleBody(bodyReader)

	assert.NoError(t, err)
	assert.Equal(t, translator.Langs{
		translator.Lang{Name: "af", Description: "Afrikaans"},
		translator.Lang{Name: "sq", Description: "Albanian"},
		translator.Lang{Name: "zh-CN", Description: "Chinese (Simplified)"},
		translator.Lang{Name: "zh-TW", Description: "Chinese (Traditional)"},
	}, langs)
}
