package utils

import (
	"fmt"
	"math/rand"
)

var wordsList = []string{
	"ipsum", "semper", "habeo", "duo", "ut", "vis", "aliquyam", "eu", "splendide", "Ut", "mei", "eteu", "nec", "antiopam", "corpora", "kasd", "pretium", "cetero", "qui", "arcu", "assentior", "ei", "his", "usu", "invidunt", "kasd", "justo", "ne", "eleifend", "per", "ut", "eam", "graeci", "tincidunt", "impedit", "temporibus", "duo", "et", "facilisis", "insolens", "consequat", "cursus", "partiendo", "ullamcorper", "Vulputate", "facilisi", "donec", "aliquam", "labore", "inimicus", "voluptua", "penatibus", "sea", "vel", "amet", "his", "ius", "audire", "in", "mea", "repudiandae", "nullam", "sed", "assentior", "takimata", "eos", "at", "odio", "consequat", "iusto", "imperdiet", "dicunt", "abhorreant", "adipisci", "officiis", "rhoncus", "leo", "dicta", "vitae", "clita", "elementum", "mauris", "definiebas", "uonsetetur", "te", "inimicus", "nec", "mus", "usu", "duo", "aenean", "corrumpit", "aliquyam", "est", "eum",
}

func getRandomWord() string {
	return wordsList[rand.Intn(len(wordsList))]
}

func GenerateWords(length int) string {
	result := ""
	for i := 0; i < length-1; i++ {
		result += getRandomWord() + " "
	}
	return result
}

func GenerateParagraphs(count, length int, separator string) string {
	result := ""
	if length == 0 {
		for i := 0; i < count; i++ {
			result += fmt.Sprintf("%s%s", GenerateWords(10), separator)
		}
		return result
	} else {
		for i := 0; i < count; i++ {
			result += fmt.Sprintf("%s%s", GenerateWords(length), separator)
		}
		return result
	}
}
