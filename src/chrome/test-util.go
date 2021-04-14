package chrome

import (
	"context"
	"log"
	"testing"
	"time"
)

func MakeChromeInstance(t *testing.T) *Chrome {
	c, err := NewChrome(context.Background(), ChromeOptions{
		Headless: true,
		Timeout:  time.Second * 15,
	})
	if err != nil {
		log.Fatal(err)
	}
	return c
}
