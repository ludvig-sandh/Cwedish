package translator

import (
	"cwedish/internal/dictionary"
	"cwedish/internal/scanner"
)

func translateTokens(untranslated []scanner.Token, dictionary *dictionary.Dictionary) (translated []scanner.Token) {
	for _, token := range untranslated {
		word := string(token[:])

		translatedWord, ok := (*dictionary)[word]
		if ok {
			translated = append(translated, scanner.Token(translatedWord))
		} else {
			translated = append(translated, token)
		}
	}
	return
}

func concatenateTokens(tokens []scanner.Token) []byte {
	total := 0
	for _, token := range tokens {
		total += len(token)
	}

	outBytes := make([]byte, 0, total)
	for _, token := range tokens {
		outBytes = append(outBytes, token...)
	}

	return outBytes
}

func Translate(in []byte, dictionary dictionary.Dictionary) (out []byte) {
	tokens := scanner.Tokenize(in)
	translatedTokens := translateTokens(tokens, &dictionary)
	out = concatenateTokens(translatedTokens)
	return
}
