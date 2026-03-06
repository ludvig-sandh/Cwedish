package main

import (
	"cwedish/internal/dictionary"
	"cwedish/internal/translator"
	"flag"
	"log"
	"os"
	"path/filepath"
)

const dictionaryFileName = "dictionary.txt"

func removeExtension(path string) string {
	extension := filepath.Ext(path)
	return path[0 : len(path)-len(extension)]
}

func changeExtension(path string, newExtension string) string {
	return removeExtension(path) + "." + newExtension
}

func parseArgs() (inFile string, outFile string, dictionaryPath string) {
	flag.StringVar(&outFile, "o", "", "output file")
	flag.StringVar(&dictionaryPath, "d", "", "custom dictionary file")
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

func loadDictionary(customPath string) dictionary.Dictionary {
	if customPath != "" {
		return dictionary.ParseDictionaryFile(customPath)
	}

	executablePath, err := os.Executable()
	if err != nil {
		log.Fatal("Executable path could not be resolved: ", err)
	}

	executableDictionaryPath := filepath.Join(filepath.Dir(executablePath), dictionaryFileName)
	if _, err := os.Stat(executableDictionaryPath); err != nil {
		log.Fatalf("Dictionary file could not be found beside the executable: %s", executableDictionaryPath)
	}

	return dictionary.ParseDictionaryFile(executableDictionaryPath)
}

func main() {
	inFile, outFile, dictionaryPath := parseArgs()

	inBytes, err := os.ReadFile(inFile)
	if err != nil {
		log.Fatal("Input file could not be read: ", err)
	}

	dict := loadDictionary(dictionaryPath)
	outBytes := translator.Translate(inBytes, dict)
	if err := os.WriteFile(outFile, outBytes, 0644); err != nil {
		log.Fatal("Output file could not be written: ", err)
	}
}
