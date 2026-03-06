package dictionary

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSplitLinesNormalizesWindowsNewlines(t *testing.T) {
	t.Parallel()

	lines := splitLines("one\r\ntwo\nthree")
	expected := []string{"one", "two", "three"}

	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(lines))
	}

	for i := range expected {
		if lines[i] != expected[i] {
			t.Fatalf("line %d: expected %q, got %q", i, expected[i], lines[i])
		}
	}
}

func TestParseDictionaryLineValid(t *testing.T) {
	t.Parallel()

	original, swedish, ok := parseDictionaryLine("for:för")
	if !ok {
		t.Fatal("expected line to parse")
	}

	if original != "for" {
		t.Fatalf("expected original %q, got %q", "for", original)
	}

	if swedish != "för" {
		t.Fatalf("expected swedish %q, got %q", "för", swedish)
	}
}

func TestParseDictionaryLineRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	testCases := []string{
		"",
		"missing-colon",
		"too:many:colons",
	}

	for _, input := range testCases {
		_, _, ok := parseDictionaryLine(input)
		if ok {
			t.Fatalf("expected %q to be rejected", input)
		}
	}
}

func TestParseDictionaryFileBuildsReverseLookup(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	dictionaryPath := filepath.Join(tempDir, "dictionary.txt")
	content := "for:för\nint:hel\ninvalid line\nwhile:medan\n"

	if err := os.WriteFile(dictionaryPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test dictionary: %v", err)
	}

	dict := ParseDictionaryFile(dictionaryPath)

	expected := map[string]string{
		"för":   "for",
		"hel":   "int",
		"medan": "while",
	}

	if len(dict) != len(expected) {
		t.Fatalf("expected %d entries, got %d", len(expected), len(dict))
	}

	for swedish, original := range expected {
		if dict[swedish] != original {
			t.Fatalf("expected %q to map to %q, got %q", swedish, original, dict[swedish])
		}
	}
}
