package chrome_test

import (
	"context"
	"testing"
	"time"
	"vnr/src/server/chrome"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
)

func TestChrome(t *testing.T) {
	makeChrome := func() (*chrome.Chrome, error) {
		return chrome.NewChrome(context.Background(), chrome.ChromeOptions{
			Headless: true,
			Timeout:  time.Second * 15,
		})
	}

	t.Run("NewChrome", func(t *testing.T) {
		chrome, err := makeChrome()
		assert.Nil(t, err)
		assert.NotNil(t, chrome)
	})
	t.Run("Run", func(t *testing.T) {
		chrome, _ := makeChrome()
		err := chrome.Run(chromedp.Navigate("http://example.com"))
		assert.Nil(t, err)
	})
}
