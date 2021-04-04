package server

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"vnr/src/translators"
)

//go:embed static
var staticFiles embed.FS

//go:embed index.html
var indexHTML []byte

type ServerOptions struct {
	Port int

	Translator translators.Translator
}

func StartServer(options ServerOptions) {
	var staticFS = http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	http.Handle("/static/", fs)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var path = req.URL.Path
		log.Println("Serving request for path", path)
		w.Header().Add("Content-Type", "text/html")
		w.Write(indexHTML)
	})

	http.HandleFunc("/translate", func(w http.ResponseWriter, r *http.Request) {
		translationOptions := translationOptionsFromQuery(r.URL.Query())
		translationResult, err := options.Translator.GetTranslation(translationOptions)
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
	})

	log.Println("Listening on :3000...")
	// start the server
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
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
