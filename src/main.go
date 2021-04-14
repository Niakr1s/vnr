package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vnr/src/server"
	"vnr/src/server/chrome"
	"vnr/src/translators"
)

func main() {
	var headlessFlag = flag.Bool("headless", true, "chrome headless mode")
	var withChrome = flag.Bool("withChrome", true, "should app use chrome")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	cancelAndWait := func() {
		cancel()
		<-time.After(time.Millisecond * 100)
	}
	defer cancelAndWait()
	handleOsInterrupt(cancelAndWait)

	var c *chrome.Chrome
	if *withChrome {
		log.Println("initializing chrome...")
		chromeInstance, err := chrome.NewChrome(ctx, chrome.ChromeOptions{
			Headless: *headlessFlag,
			Timeout:  time.Second * 15,
		})
		if err != nil {
			log.Fatalf("couldn't initialize chrome: %v", err)
		}
		c = chromeInstance
	}

	transl, err := translators.GetAllKnownTranslators(translators.GetTranslatorOptions{
		Chrome: c,
	})
	if err != nil {
		log.Fatalf("couldn't get translators: %v", err)
	}

	server.StartServer(server.ServerOptions{
		Port:        env("PORT", ":5322"),
		Translators: transl,
	})
}

func env(k string, defaultV string) string {
	v, ok := os.LookupEnv(k)
	if !ok {
		return defaultV
	}
	return v
}

func handleOsInterrupt(fns ...func()) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		for _, fn := range fns {
			log.Printf("inside handleOsInterrupt")
			fn()
		}
		<-time.After(time.Millisecond * 100)
		os.Exit(0)
	}()
}
