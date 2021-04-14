package translators

import (
	"context"
	"log"
	"testing"
	"time"
	"vnr/src/server/chrome"
)

func makeChromeInstance(t *testing.T) *chrome.Chrome {
	c, err := chrome.NewChrome(context.Background(), chrome.ChromeOptions{
		Headless: true,
		Timeout:  time.Second * 15,
	})
	if err != nil {
		log.Fatal(err)
	}
	return c
}
