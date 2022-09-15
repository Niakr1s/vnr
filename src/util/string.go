package util

var jpPunctuationMarks []rune = []rune{'。', '…', '？', '！'}

func SplitToSentencesJP(input string) []string {
	res := []string{}

	wasPunctiationMark := false
	sentence := ""
	for _, r := range input {
		if SliceContains(jpPunctuationMarks, r) {
			wasPunctiationMark = true
		} else {
			if wasPunctiationMark {
				res = append(res, sentence)
				sentence = ""
				wasPunctiationMark = false
			}
		}
		sentence += string(r)
	}
	res = append(res, sentence)

	return res
}
