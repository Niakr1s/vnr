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
	"vnr/src/translators"
)

var headlessFlag = flag.Bool("headless", true, "chrome headless mode")
var translatorFlag = flag.String("translator", translators.KnownTranslators[0], "translator")

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	cancelAndWait := func() {
		cancel()
		<-time.After(time.Millisecond * 100)
	}
	defer cancelAndWait()
	handleOsInterrupt(cancelAndWait)

	translator, err := translators.GetTranslator(*translatorFlag)
	if err != nil {
		panic(err)
	}

	translator.Init(ctx, translators.TranslatorInitOptions{
		Headless: *headlessFlag,
	})

	server.StartServer(server.ServerOptions{
		Port:       3000,
		Translator: translator,
	})

	// words := []string{
	// 	"こんいちは",
	// 	"hello",
	// 	"schwartz",
	// }
	// wg := sync.WaitGroup{}
	// wg.Add(len(words))

	// log.Printf("fetching translations...")
	// for i, word := range words {
	// 	word := word
	// 	i := i
	// 	go func() {
	// 		defer wg.Done()

	// 		translationOptions := translators.NewTranslationOptions(word)
	// 		translationOptions.To = "ru"
	// 		if i == 1 {
	// 			translationOptions.Timeont = time.Millisecond * 1000
	// 		}

	// 		translation, err := translator.GetTranslation(translationOptions)
	// 		if err != nil {
	// 			log.Printf("wrror while fetching translation, %s", err)
	// 		} else {
	// 			log.Printf("got translation: %s => %s", translation.TranslationOptions, translation.Translation)
	// 		}
	// 	}()
	// }

	// wg.Wait()
	// log.Printf("translations fetched")
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
