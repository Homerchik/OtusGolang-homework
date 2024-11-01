package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Word struct {
	Word string
	Qty  int
}

func (w *Word) Greater(other *Word) bool {
	if w.Qty > other.Qty {
		return true
	}
	if w.Qty == other.Qty {
		return strings.Compare(w.Word, other.Word) < 0
	}
	return false
}

func Top10(s string) []string {
	// Place your code here.
	words := make(map[string]int)
	for _, word := range strings.Fields(s) {
		words[word]++
	}
	wordsList := make([]Word, len(words))
	i := 0
	for k, v := range words {
		wordsList[i] = Word{k, v}
		i++
	}
	sort.SliceStable(wordsList, func(i int, j int) bool { return wordsList[i].Greater(&wordsList[j]) })
	wordsSlice := make([]string, min(10, i))
	for i := range wordsSlice {
		wordsSlice[i] = wordsList[i].Word
	}
	return wordsSlice
}
