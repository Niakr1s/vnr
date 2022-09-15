package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"vnr/src/translator"
	"vnr/src/util"
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
		http.HandleFunc(fmt.Sprintf("/api/translate/%s", name), translationHandler(name, translator))
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

func translationHandler(translatorName string, transl Translator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		translationOptions := translationOptionsFromQuery(r.URL.Query())
		log.Printf("%s: translate start: %+v", translatorName, translationOptions)

		var translationResult translator.TranslationResult
		var err error
		if translationOptions.Single {
			translationResult, err = getTranslation(transl, translationOptions)
		} else {
			translationResult, err = getTranslationSplitBySentences(transl, translationOptions)
		}

		if err != nil {
			log.Printf("%s: translate failure: %+v, reason: %v", translatorName, translationOptions, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		log.Printf("%s: translate success: %+v", translatorName, translationResult)
		translationResultJson, err := json.Marshal(translationResult)
		log.Printf("%s: translate success, json: %+v", translatorName, string(translationResultJson))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(translationResultJson)
	}
}

func getTranslation(transl Translator, translationOptions translator.TranslationOptions) (translator.TranslationResult, error) {
	return transl.GetTranslation(translationOptions)
}

func getTranslationSplitBySentences(transl Translator, translationOptions translator.TranslationOptions) (translator.TranslationResult, error) {
	sentences := util.SplitToSentencesJP(translationOptions.Sentence)

	translationOptionsSplitted := []translator.TranslationOptions{}
	for _, sentence := range sentences {
		option := translationOptions
		option.Sentence = sentence
		translationOptionsSplitted = append(translationOptionsSplitted, option)
	}

	translationResults := make([]translator.TranslationResult, len(translationOptionsSplitted))
	translationErrors := make([]error, len(translationOptionsSplitted))

	var wg sync.WaitGroup
	for i, option := range translationOptionsSplitted {
		wg.Add(1)
		i := i
		option := option
		go func() {
			defer wg.Done()

			res, err := transl.GetTranslation(option)
			translationResults[i] = res
			translationErrors[i] = err
		}()
	}
	wg.Wait()

	resultErrStr := ""
	for _, err := range translationErrors {
		if err != nil {
			resultErrStr += err.Error() + "; "
		}
	}
	if resultErrStr != "" {
		return translator.TranslationResult{}, fmt.Errorf(resultErrStr)
	}

	finalTranslationResultsStr := []string{}
	for _, transRes := range translationResults {
		finalTranslationResultsStr = append(finalTranslationResultsStr, transRes.Translation)
	}

	res := translator.TranslationResult{
		TranslationOptions: translationOptions,
		Translation:        strings.Join(finalTranslationResultsStr, ""),
	}

	return res, nil
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
	single := query.Get("single")

	translationOptions := translator.NewTranslationOptions(sentence)
	if from != "" {
		translationOptions.From = from
	}
	if to != "" {
		translationOptions.To = to
	}
	translationOptions.Single = single == "true"

	return translationOptions
}
