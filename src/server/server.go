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
	Port int

	Translator translators.Translator
}

func StartServer(options ServerOptions) {
	staticFilesRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}

	var staticFS = http.FS(staticFilesRoot)
	fs := http.FileServer(staticFS)

	if !isDevMode {
		http.Handle("/", fs)
	}

	http.HandleFunc("/api/translate", translationHandler(options.Translator))

	log.Println("Listening on :5322...")
	// start the server
	err = http.ListenAndServe(":5322", nil)
	if err != nil {
		log.Fatal(err)
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
