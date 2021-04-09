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
	var translatorFlag = flag.String("translator", translators.KnownTranslators[0], "translator")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	cancelAndWait := func() {
		cancel()
		<-time.After(time.Millisecond * 100)
	}
	defer cancelAndWait()
	handleOsInterrupt(cancelAndWait)

	chrome, err := chrome.NewChrome(ctx, chrome.ChromeOptions{
		Headless: *headlessFlag,
		Timeout:  time.Second * 15,
	})

	if err != nil {
		log.Fatalf("couldn't initialize chrome: %v", err)
	}

	translator, err := translators.GetTranslator(translators.GetTranslatorOptions{
		Translator: *translatorFlag,
		Chrome:     chrome,
	})
	if err != nil {
		panic(err)
	}

	server.StartServer(server.ServerOptions{
		Port:       env("PORT", ":5322"),
		Translator: translator,
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
