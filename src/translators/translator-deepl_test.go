package translators

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeeplTranslator(t *testing.T) {
	deepl := NewDeeplTranslator(makeChromeInstance(t))

	t.Run("GetLanguages", func(t *testing.T) {
		langs, err := deepl.GetLanguages()
		assert.Nil(t, err)
		assert.NotZero(t, len(langs))
	})

	t.Run("GetTranslation", func(t *testing.T) {
		translationOptions := TranslationOptions{From: "en", To: "ru", Sentence: "hello"}
		res, err := deepl.GetTranslation(translationOptions)
		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func Test_getLanguagesFromDeeplBody(t *testing.T) {
	body := `
<li><button class='docTrans_document__target_lang_select__entry' dl-test='doctrans-upload-lang-item' dl-lang='en-US' tabindex='0' dl-attr="onClick{'en-US'}: $0.targetLang">English (American)</button></li>
<li><button class='docTrans_document__target_lang_select__entry' dl-test='doctrans-upload-lang-item' dl-lang='en-GB' tabindex='0' dl-attr="onClick{'en-GB'}: $0.targetLang">English (British)</button></li>
<li><button class='docTrans_document__target_lang_select__entry' dl-test='doctrans-upload-lang-item' dl-lang='de-DE' tabindex='0' dl-attr="onClick{'DE'}: $0.targetLang">German</button></li>
<li><button class='docTrans_document__target_lang_select__entry' dl-test='doctrans-upload-lang-item' dl-lang='fr-FR' tabindex='0' dl-attr="onClick{'FR'}: $0.targetLang">French</button></li>
    `
	bodyReader := strings.NewReader(body)

	langs, err := getLanguagesFromDeeplBody(bodyReader)

	assert.NoError(t, err)
	assert.Equal(t, Langs{Lang{"en", "English"}, Lang{"de", "German"}, Lang{"fr", "French"}}, langs)
}
