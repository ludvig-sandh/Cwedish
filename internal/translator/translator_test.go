package translator

import (
	"cwedish/internal/dictionary"
	"strings"
	"testing"
)

func fullDictionary() dictionary.Dictionary {
	return dictionary.Dictionary{
		"auto":      "auto",
		"bryt":      "break",
		"fall":      "case",
		"kar":       "char",
		"konst":     "const",
		"fortsätt":  "continue",
		"standard":  "default",
		"gör":       "do",
		"dubbel":    "double",
		"annars":    "else",
		"uppr":      "enum",
		"extern":    "extern",
		"flyt":      "float",
		"för":       "for",
		"gåtill":    "goto",
		"om":        "if",
		"hel":       "int",
		"lång":      "long",
		"register":  "register",
		"returnera": "return",
		"kort":      "short",
		"tecknad":   "signed",
		"storlekav": "sizeof",
		"statisk":   "static",
		"strukt":    "struct",
		"byt":       "switch",
		"typdef":    "typedef",
		"union":     "union",
		"otecknad":  "unsigned",
		"tom":       "void",
		"volatil":   "volatile",
		"medan":     "while",
	}
}

func verifyTranslation(t *testing.T, input string, expected string, dict dictionary.Dictionary) {
	t.Helper()

	inBytes := []byte(input)
	outBytes := Translate(inBytes, dict)
	output := string(outBytes)

	if output != expected {
		t.Fatalf("expected translation %q, got %q", expected, output)
	}
}

func TestTranslateKeywordSeparators(t *testing.T) {
	t.Parallel()

	dict := fullDictionary()
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "for before opening paren",
			input:    "för(i=0;i<10;i++){}",
			expected: "for(i=0;i<10;i++){}",
		},
		{
			name:     "if around parentheses and braces",
			input:    "om(x){returnera;}",
			expected: "if(x){return;}",
		},
		{
			name:     "while after closing brace",
			input:    "}medan(x>0);",
			expected: "}while(x>0);",
		},
		{
			name:     "switch case default punctuation",
			input:    "byt(x){fall 1:returnera;standard:returnera;}",
			expected: "switch(x){case 1:return;default:return;}",
		},
		{
			name:     "struct and typedef with semicolons",
			input:    "typdef strukt Node{hel värde;};",
			expected: "typedef struct Node{int värde;};",
		},
		{
			name:     "comma separated declarations",
			input:    "hel a,b;kort c,d;",
			expected: "int a,b;short c,d;",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			verifyTranslation(t, tc.input, tc.expected, dict)
		})
	}
}

func TestTranslateKeywordsAroundOperators(t *testing.T) {
	t.Parallel()

	dict := fullDictionary()
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "for with increment operator",
			input:    "för(hel i=0;i<10;i++){}",
			expected: "for(int i=0;i<10;i++){}",
		},
		{
			name:     "while with decrement operator",
			input:    "medan(i-->0){}",
			expected: "while(i-->0){}",
		},
		{
			name:     "return after shift assign",
			input:    "x<<=1;returnera x;",
			expected: "x<<=1;return x;",
		},
		{
			name:     "continue after plus assign",
			input:    "x+=1;fortsätt;",
			expected: "x+=1;continue;",
		},
		{
			name:     "bitwise operators near if",
			input:    "om(a&&b||c){returnera;}",
			expected: "if(a&&b||c){return;}",
		},
		{
			name:     "sizeof before identifier",
			input:    "storlekav(värde)+storlekav *ptr",
			expected: "sizeof(värde)+sizeof *ptr",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			verifyTranslation(t, tc.input, tc.expected, dict)
		})
	}
}

func TestTranslateLeavesKeywordsInsideIdentifiersUntouched(t *testing.T) {
	t.Parallel()

	dict := fullDictionary()
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "prefix match in identifier",
			input:    "förify();",
			expected: "förify();",
		},
		{
			name:     "suffix match in identifier",
			input:    "myom=1;",
			expected: "myom=1;",
		},
		{
			name:     "keyword between letters and digits",
			input:    "xreturnera2=0;",
			expected: "xreturnera2=0;",
		},
		{
			name:     "typedef like identifier",
			input:    "typdefad=3;",
			expected: "typdefad=3;",
		},
		{
			name:     "struct name containing keyword",
			input:    "struktNode n;",
			expected: "struktNode n;",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			verifyTranslation(t, tc.input, tc.expected, dict)
		})
	}
}

func TestTranslateLeavesKeywordsInsideStringsAndCharsUntouched(t *testing.T) {
	t.Parallel()

	dict := fullDictionary()
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "double quoted string",
			input:    "printf(\"för om medan returnera\");",
			expected: "printf(\"för om medan returnera\");",
		},
		{
			name:     "escaped quote inside string",
			input:    "printf(\"ordet \\\"för\\\" ska stanna\");",
			expected: "printf(\"ordet \\\"för\\\" ska stanna\");",
		},
		{
			name:     "single quoted char",
			input:    "kar c='f';returnera c;",
			expected: "char c='f';return c;",
		},
		{
			name:     "escaped single quote char literal",
			input:    "kar q='\\'';returnera q;",
			expected: "char q='\\'';return q;",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			verifyTranslation(t, tc.input, tc.expected, dict)
		})
	}
}

func TestTranslateLeavesKeywordsInsideCommentsUntouched(t *testing.T) {
	t.Parallel()

	dict := fullDictionary()
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single line comment",
			input:    "hel x; // för om returnera\nreturnera x;",
			expected: "int x; // för om returnera\nreturn x;",
		},
		{
			name:     "single line comment with slash slash in code before it",
			input:    "x/=2; // fortsätt senare\nfortsätt;",
			expected: "x/=2; // fortsätt senare\ncontinue;",
		},
		{
			name:     "multi line comment",
			input:    "hel x; /* för\nom\nreturnera */ returnera x;",
			expected: "int x; /* för\nom\nreturnera */ return x;",
		},
		{
			name:     "comment between translated keywords",
			input:    "om(x)/* annars */annars returnera;",
			expected: "if(x)/* annars */else return;",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			verifyTranslation(t, tc.input, tc.expected, dict)
		})
	}
}

func TestTranslatePreservesWhitespaceBoundaries(t *testing.T) {
	t.Parallel()

	dict := fullDictionary()
	verifyTranslation(
		t,
		"hel\tmain(\tom\t)\n{\n\treturnera 0;\n}\n",
		"int\tmain(\tif\t)\n{\n\treturn 0;\n}\n",
		dict,
	)
}

func TestTranslateCompoundDeclarationsAndControlFlow(t *testing.T) {
	t.Parallel()

	dict := fullDictionary()
	verifyTranslation(
		t,
		"statisk konst otecknad lång hel räknare=0; medan(räknare<3){räknare++; om(räknare==2){fortsätt;} annars {returnera;}}",
		"static const unsigned long int räknare=0; while(räknare<3){räknare++; if(räknare==2){continue;} else {return;}}",
		dict,
	)
}

func TestTranslateSwitchEnumStructUnionKeywords(t *testing.T) {
	t.Parallel()

	dict := fullDictionary()
	verifyTranslation(
		t,
		"uppr Färg{RÖD,GRÖN}; union Data{hel i; dubbel d;}; byt(v){fall 0:returnera;standard:returnera;}",
		"enum Färg{RÖD,GRÖN}; union Data{int i; double d;}; switch(v){case 0:return;default:return;}",
		dict,
	)
}

func TestTranslateAllDictionaryKeywords(t *testing.T) {
	t.Parallel()

	dict := fullDictionary()

	inputWords := []string{
		"auto",
		"bryt",
		"fall",
		"kar",
		"konst",
		"fortsätt",
		"standard",
		"gör",
		"dubbel",
		"annars",
		"uppr",
		"extern",
		"flyt",
		"för",
		"gåtill",
		"om",
		"hel",
		"lång",
		"register",
		"returnera",
		"kort",
		"tecknad",
		"storlekav",
		"statisk",
		"strukt",
		"byt",
		"typdef",
		"union",
		"otecknad",
		"tom",
		"volatil",
		"medan",
	}

	expectedWords := []string{
		"auto",
		"break",
		"case",
		"char",
		"const",
		"continue",
		"default",
		"do",
		"double",
		"else",
		"enum",
		"extern",
		"float",
		"for",
		"goto",
		"if",
		"int",
		"long",
		"register",
		"return",
		"short",
		"signed",
		"sizeof",
		"static",
		"struct",
		"switch",
		"typedef",
		"union",
		"unsigned",
		"void",
		"volatile",
		"while",
	}

	verifyTranslation(
		t,
		strings.Join(inputWords, " "),
		strings.Join(expectedWords, " "),
		dict,
	)
}
