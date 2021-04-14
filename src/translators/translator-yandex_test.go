package translators

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYandexTranslator(t *testing.T) {
	yandex := NewDeeplTranslator(makeChromeInstance(t))

	t.Run("GetLanguages", func(t *testing.T) {
		langs, err := yandex.GetLanguages()
		assert.Nil(t, err)
		assert.NotZero(t, len(langs))
	})

	t.Run("GetTranslation", func(t *testing.T) {
		translationOptions := TranslationOptions{From: "en", To: "ru", Sentence: "hello"}
		res, err := yandex.GetTranslation(translationOptions)
		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func Test_getLanguagesFromYandexBody(t *testing.T) {
	body := `
TRANSLATOR_LANGS: {"af":"Afrikaans","sq":"Albanian","am":"Amharic","ar":"Arabic"},
    `
	bodyReader := strings.NewReader(body)

	langs, err := getLanguagesFromYandexBody(bodyReader)

	assert.NoError(t, err)
	assert.Equal(t, Langs{Lang{"af", "Afrikaans"}, Lang{"sq", "Albanian"}, Lang{"am", "Amharic"}, Lang{"ar", "Arabic"}}, langs)
}
