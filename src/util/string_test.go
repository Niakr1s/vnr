package util

import (
	"reflect"
	"testing"
)

func TestSplitToSentences(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty string",
			input: "",
			want:  []string{""},
		},
		{
			name:  "no punctuation marks",
			input: "宗一郎「そうそうそれで良いのじゃ」",
			want:  []string{"宗一郎「そうそうそれで良いのじゃ」"},
		},
		{
			name:  "single punctuation mark",
			input: "。",
			want:  []string{"。"},
		},
		{
			name:  "one punctuation mark",
			input: "宗一郎「そうそう。それで良いのじゃ」",
			want:  []string{"宗一郎「そうそう。", "それで良いのじゃ」"},
		},
		{
			name:  "three dots",
			input: "宗一郎「そうそう。。。それで良いのじゃ」",
			want:  []string{"宗一郎「そうそう。。。", "それで良いのじゃ」"},
		},
		{
			name:  "three dots and question",
			input: "宗一郎「そうそう。。。？それで良いのじゃ」",
			want:  []string{"宗一郎「そうそう。。。？", "それで良いのじゃ」"},
		},
		{
			name:  "three dots and question",
			input: "宗一郎「そうそう。。。？それ。で良い。のじゃ？」",
			want:  []string{"宗一郎「そうそう。。。？", "それ。", "で良い。", "のじゃ？", "」"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitToSentencesJP(tt.input); !reflect.DeepEqual(tt.want, got) {
				t.Errorf("SplitToSentences() = %v, want %v", got, tt.want)
			}
		})
	}
}
