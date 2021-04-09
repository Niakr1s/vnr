package translators_test

import (
	"context"
	"log"
	"testing"
	"time"
	"vnr/src/server/chrome"
	"vnr/src/translators"
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
	deepl := translators.NewDeeplTranslator(chromeInstance)

	t.Run("GetLanguages", func(t *testing.T) {
		langs, err := deepl.GetLanguages()
		assert.Nil(t, err)
		assert.NotZero(t, len(langs))
		assert.True(t, util.SliceContainsString(langs, "en"))
	})
}
