package deepl

import (
	"strings"
	"testing"
	"vnr/src/translator"

	"github.com/stretchr/testify/assert"
)

func Test_DeeplGetTranslationWithoutChrome(t *testing.T) {
	deepl := NewDeeplTranslator()
	translationOptions := translator.TranslationOptions{From: "en", To: "ru", Sentence: "hello"}
	res, err := deepl.GetTranslation(translationOptions)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "приветствие", res.Translation)
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
	assert.Equal(t, translator.Langs{
		translator.Lang{Name: "en", Description: "English"},
		translator.Lang{Name: "de", Description: "German"},
		translator.Lang{Name: "fr", Description: "French"},
	}, langs)
}

func Test_getTranslationFromDeeplJsonRpcBody(t *testing.T) {
	body := `{
    "jsonrpc": "2.0",
    "id": 84110018,
    "result": {
        "translations": [
            {
                "beams": [
                    {
                        "postprocessed_sentence": "Здравствуйте",
                        "num_symbols": 8
                    },
                    {
                        "postprocessed_sentence": "Здравствуйте:",
                        "num_symbols": 9
                    },
                    {
                        "postprocessed_sentence": "Здравствуйте .",
                        "num_symbols": 9
                    },
                    {
                        "postprocessed_sentence": "Здравствуйте,",
                        "num_symbols": 9
                    }
                ],
                "quality": "normal"
            }
        ],
        "target_lang": "RU",
        "source_lang": "DE",
        "source_lang_is_confident": false,
        "detectedLanguages": {},
        "timestamp": 1618417753,
        "date": "20210414"
    }
}`

	translation, err := getTranslationFromDeeplJsonRpcBody(strings.NewReader(body))
	assert.Nil(t, err)
	assert.Equal(t, "Здравствуйте", translation)
}
