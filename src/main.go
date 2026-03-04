package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Not enough arguments: missing path to translation file")
		os.Exit(1)
	}

	fileToTranslate := os.Args[1]
	bytes, err := os.ReadFile(fileToTranslate)
	if err != nil {
		fmt.Println("Input file could not be read: ", err)
		os.Exit(1)
	}

	tokens := Tokenize(bytes)
	// for i, token := range tokens {
	// 	fmt.Printf("Token %d: \"%s\"\n", i, string(token[:]))
	// }

	dictionary := ParseMappingsFile("keywords-table.txt")
	// fmt.Println(dictionary)

	translatedTokens := []Token{}
	for _, token := range tokens {
		word := string(token[:])

		translated, ok := dictionary[word]
		if ok {
			translatedTokens = append(translatedTokens, Token(translated))
		} else {
			translatedTokens = append(translatedTokens, token)
		}
	}

	outBytes := []byte{}
	for _, token := range translatedTokens {
		for _, b := range token {
			outBytes = append(outBytes, b)
		}
	}

	os.WriteFile("example.c", outBytes, 0644)
}
