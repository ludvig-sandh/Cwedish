package main

import (
	"cwedish/internal/dictionary"
	"cwedish/internal/scanner"
	"flag"
	"log"
	"os"
	"path/filepath"
)

func removeExtension(path string) string {
	extension := filepath.Ext(path)
	return path[0 : len(path)-len(extension)]
}

func changeExtension(path string, newExtension string) string {
	return removeExtension(path) + "." + newExtension
}

func parseArgs() (inFile string, outFile string) {
	flag.StringVar(&outFile, "o", "", "output file")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("Not enough arguments: missing path to translation file")
	}

	inFile = flag.Arg(0)

	if outFile == "" {
		outFile = changeExtension(inFile, "c")
	}

	return
}

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

func main() {
	inFile, outFile := parseArgs()

	bytes, err := os.ReadFile(inFile)
	if err != nil {
		log.Fatal("Input file could not be read: ", err)
	}

	dictionary := dictionary.ParseDictionaryFile("keywords-table.txt")
	tokens := scanner.Tokenize(bytes)
	translatedTokens := translateTokens(tokens, &dictionary)
	outBytes := concatenateTokens(translatedTokens)
	os.WriteFile(outFile, outBytes, 0644)
}
