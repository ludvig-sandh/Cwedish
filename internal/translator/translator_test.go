package translator

import (
	"cwedish/internal/dictionary"
	"testing"
)

func verifyTranslation(input string, expected string, dict dictionary.Dictionary, t *testing.T) {
	inBytes := []byte(input)
	outBytes := Translate(inBytes, dict)
	output := string(outBytes)

	if output != expected {
		t.Fatalf("expected translation '%s', got '%s'", expected, output)
	}
}

func TestTranslateForLoop(t *testing.T) {
	dict := dictionary.Dictionary{
		"för": "for",
	}
	verifyTranslation(
		"för(int i=0;i<100;i++){}",
		"for(int i=0;i<100;i++){}",
		dict,
		t,
	)
}
