package server

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"net/url"
	"time"
	"vnr/src/translator"

	"github.com/tjgq/broadcast"
	"golang.design/x/clipboard"
)

//go:embed static
var staticFiles embed.FS

type Translator interface {
	GetTranslation(translationOptions translator.TranslationOptions) (translator.TranslationResult, error)
	GetLanguages() (translator.Langs, error)
}

type ServerOptions struct {
	Port string

	Translators map[string]Translator
}

type ClipboardPollHandlerOptions struct {
	Delay time.Duration
}

func StartServer(options ServerOptions) {
	mime.AddExtensionType(".js", "application/javascript")

	staticFilesRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}

	var staticFS = http.FS(staticFilesRoot)
	fs := http.FileServer(staticFS)

	http.Handle("/", fs)

	http.HandleFunc("/api/knownTranslators", knownTranslationsHandler(getTranslatorNames(options.Translators)))

	for name, translator := range options.Translators {
		http.HandleFunc(fmt.Sprintf("/api/translate/%s", name), translationHandler(name, translator))
		http.HandleFunc(fmt.Sprintf("/api/langs/%s", name), langsHandler(translator))
	}

	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/api/clipboard", clipboardHandler())
	http.HandleFunc("/api/clipboardPoll", clipboardPollHandler(ClipboardPollHandlerOptions{
		Delay: time.Second * 10,
	}))

	log.Printf("Listening on %s...", options.Port)
	// start the server
	err = http.ListenAndServe(options.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getTranslatorNames(m map[string]Translator) []string {
	res := make([]string, 0)
	for k := range m {
		res = append(res, k)
	}
	return res
}

func knownTranslationsHandler(translatorNames []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(translatorNames)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

func langsHandler(translator Translator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		langs, err := translator.GetLanguages()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		writeJson(langs, w)
	}
}

func translationHandler(translatorName string, translator Translator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		translationOptions := translationOptionsFromQuery(r.URL.Query())
		log.Printf("%s: translate start: %+v", translatorName, translationOptions)
		translationResult, err := translator.GetTranslation(translationOptions)
		if err != nil {
			log.Printf("%s: translate failure: %+v, reason: %v", translatorName, translationOptions, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		log.Printf("%s: translate success: %+v", translatorName, translationResult)
		translationResultJson, err := json.Marshal(translationResult)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(translationResultJson)
	}
}

func writeJson(obj interface{}, w http.ResponseWriter) {
	translationResultJson, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(translationResultJson)
}

func translationOptionsFromQuery(query url.Values) translator.TranslationOptions {
	sentence := query.Get("sentence")
	from := query.Get("from")
	to := query.Get("to")

	translationOptions := translator.NewTranslationOptions(sentence)
	if from != "" {
		translationOptions.From = from
	}
	if to != "" {
		translationOptions.To = to
	}
	return translationOptions
}

func clipboardHandler() http.HandlerFunc {
	type Response struct {
		Clipboard string `json:"clipboard"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		contents := clipboard.Read(clipboard.FmtText)
		res := Response{
			Clipboard: string(contents),
		}
		writeJson(res, w)
	}
}

// Returns empty string, if clipboard wasn't changed
func clipboardPollHandler(options ClipboardPollHandlerOptions) http.HandlerFunc {
	type Response struct {
		Clipboard string `json:"clipboard"`
	}
	res := Response{}

	clipboardBroadcaster := broadcast.New(0)
	go func() {
		clipboardCh := clipboard.Watch(context.TODO(), clipboard.FmtText)
		for {
			contents := <-clipboardCh
			clipboardBroadcaster.Send(contents)
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		l := clipboardBroadcaster.Listen()
		defer l.Close()

		select {
		case contents := <-l.Ch:
			res.Clipboard = string(contents.([]byte))

		case <-time.After(options.Delay):
		}

		writeJson(res, w)
	}
}
