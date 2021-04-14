package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"vnr/src/translator"
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

func StartServer(options ServerOptions) {
	staticFilesRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}

	var staticFS = http.FS(staticFilesRoot)
	fs := http.FileServer(staticFS)

	http.Handle("/", fs)

	http.HandleFunc("/api/knownTranslators", knownTranslationsHandler(getTranslatorNames(options.Translators)))

	for name, translator := range options.Translators {
		http.HandleFunc(fmt.Sprintf("/api/translate/%s", name), translationHandler(translator))
		http.HandleFunc(fmt.Sprintf("/api/langs/%s", name), langsHandler(translator))
	}

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

func translationHandler(translator Translator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		translationOptions := translationOptionsFromQuery(r.URL.Query())
		log.Printf("translate start: %+v", translationOptions)
		translationResult, err := translator.GetTranslation(translationOptions)
		if err != nil {
			log.Printf("translate failure: %+v", translationOptions)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		log.Printf("translate success: %+v", translationResult)
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
