package server

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"vnr/src/translators"
)

//go:embed static
var staticFiles embed.FS

type ServerOptions struct {
	Port string

	Translator translators.Translator
}

func StartServer(options ServerOptions) {
	staticFilesRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}

	var staticFS = http.FS(staticFilesRoot)
	fs := http.FileServer(staticFS)

	http.Handle("/", fs)

	http.HandleFunc("/api/knownTranslators", knownTranslationsHandler())

	http.HandleFunc("/api/translate", translationHandler(options.Translator))

	log.Printf("Listening on %s...", options.Port)
	// start the server
	err = http.ListenAndServe(options.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func knownTranslationsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(translators.KnownTranslators)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

func translationHandler(translator translators.Translator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		translationOptions := translationOptionsFromQuery(r.URL.Query())
		translationResult, err := translator.GetTranslation(translationOptions)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		translationResultJson, err := json.Marshal(translationResult)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(translationResultJson)
	}
}

func translationOptionsFromQuery(query url.Values) translators.TranslationOptions {
	sentence := query.Get("sentence")
	from := query.Get("from")
	to := query.Get("to")

	translationOptions := translators.NewTranslationOptions(sentence)
	if from != "" {
		translationOptions.From = from
	}
	if to != "" {
		translationOptions.To = to
	}
	return translationOptions
}
