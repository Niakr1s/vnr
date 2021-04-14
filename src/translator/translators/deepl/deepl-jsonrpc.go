package deepl

/*
                    GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.
*/

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"vnr/src/translator"
)

// getDeeplTranslationRpcRequest makes request.
// This part of code partly was copied from https://github.com/Artikash/Textractor/blob/master/extensions/deepltranslate.cpp
func getDeeplTranslationRpcRequest(translationOptions translator.TranslationOptions) (*http.Request, error) {
	r := time.Now().UnixNano() / (1000000)
	n := strings.Count(translationOptions.Sentence, "i") + 1
	id := 10000*rand.Intn(9999) + 1
	var timeStamp int64 = r + (int64(n) - r%int64(n))
	requestBody := fmt.Sprintf(`
{
    "jsonrpc": "2.0",
    "method": "LMT_handle_jobs",
    "params": {
        "jobs": [
            {
                "kind": "default",
                "raw_en_sentence": "%s",
                "raw_en_context_before": [],
                "raw_en_context_after": [],
                "preferred_num_beams": 4
            }
        ],
        "lang": {
            "source_lang_user_selected": "%s",
            "target_lang": "%s"
        },
        "timestamp": %d
    },
    "id": %d
}
    `, translationOptions.Sentence, strings.ToUpper(translationOptions.From), strings.ToUpper(translationOptions.To), timeStamp, id)

	req, err := http.NewRequest(http.MethodPost, "https://www2.deepl.com/jsonrpc", strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Origin", "https://www.deepl.com")
	req.Header.Add("Referer", "https://www.deepl.com/")
	req.Header.Add("TE", "Trailers")

	return req, nil
}
