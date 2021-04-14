package translator

type Lang struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Langs []Lang

func (l Langs) RemoveDuplicates() Langs {
	set := map[string]Lang{}

	res := Langs{}
	for _, lang := range l {
		if _, ok := set[lang.Name]; ok {
			continue
		}
		set[lang.Name] = lang
		res = append(res, lang)
	}
	return res
}

type TranslationOptions struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Sentence string `json:"sentence"`
}

func NewTranslationOptions(sentence string) TranslationOptions {
	return TranslationOptions{
		From:     "auto",
		To:       "auto",
		Sentence: sentence,
	}
}

type TranslationResult struct {
	TranslationOptions
	Translation string `json:"translation"`
}
