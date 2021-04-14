package main

import (
	"os"
	"vnr/src/server"
	"vnr/src/translator/translators/deepl"
	"vnr/src/translator/translators/google"
)

func main() {
	server.StartServer(server.ServerOptions{
		Port: env("PORT", ":5322"),
		Translators: map[string]server.Translator{
			"deepl":  deepl.NewDeeplTranslator(),
			"google": google.NewGoogleTranslator(),
		},
	})
}

func env(k string, defaultV string) string {
	v, ok := os.LookupEnv(k)
	if !ok {
		return defaultV
	}
	return v
}
