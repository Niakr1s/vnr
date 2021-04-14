package translators

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"
	"vnr/src/server/chrome"
	"vnr/src/util"

	"github.com/stretchr/testify/assert"
)

var chromeInstance *chrome.Chrome

func init() {
	c, err := chrome.NewChrome(context.Background(), chrome.ChromeOptions{
		Headless: true,
		Timeout:  time.Second * 15,
	})
	if err != nil {
		log.Fatal(err)
	}
	chromeInstance = c
}

func TestDeeplTranslator(t *testing.T) {
	deepl := NewDeeplTranslator(chromeInstance)

	t.Run("GetLanguages", func(t *testing.T) {
		langs, err := deepl.GetLanguages()
		assert.Nil(t, err)
		assert.NotZero(t, len(langs))
		assert.True(t, util.SliceContainsString(langs, "en"))
	})
}

func Test_getLanguagesFromDeeplBody(t *testing.T) {
	body := `
<li><button class='docTrans_document__target_lang_select__entry' dl-test='doctrans-upload-lang-item' dl-lang='en-GB' tabindex='0' dl-attr="onClick{'en-GB'}: $0.targetLang">English (British)</button></li>
<li><button class='docTrans_document__target_lang_select__entry' dl-test='doctrans-upload-lang-item' dl-lang='de-DE' tabindex='0' dl-attr="onClick{'DE'}: $0.targetLang">German</button></li>
<li><button class='docTrans_document__target_lang_select__entry' dl-test='doctrans-upload-lang-item' dl-lang='fr-FR' tabindex='0' dl-attr="onClick{'FR'}: $0.targetLang">French</button></li>
    `
	bodyReader := strings.NewReader(body)

	langs, err := getLanguagesFromDeeplBody(bodyReader)

	assert.NoError(t, err)
	assert.Equal(t, []string{"en", "de", "fr"}, langs)
}
