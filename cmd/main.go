package main

import (
	"cwedish/internal/translator"
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

func main() {
	inFile, outFile := parseArgs()

	inBytes, err := os.ReadFile(inFile)
	if err != nil {
		log.Fatal("Input file could not be read: ", err)
	}

	outBytes := translator.Translate(inBytes)
	os.WriteFile(outFile, outBytes, 0644)
}
