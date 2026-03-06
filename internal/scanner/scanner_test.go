package scanner

import "testing"

type tokenizeTestCase struct {
	name     string
	input    string
	expected []string
}

func assertTokens(t *testing.T, input string, expected []string) {
	t.Helper()

	tokens := Tokenize([]byte(input))
	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d: %q", len(expected), len(tokens), stringifyTokens(tokens))
	}

	for i := range expected {
		actual := string(tokens[i])
		if actual != expected[i] {
			t.Fatalf("token %d: expected %q, got %q", i, expected[i], actual)
		}
	}
}

func stringifyTokens(tokens []Token) []string {
	out := make([]string, len(tokens))
	for i, token := range tokens {
		out[i] = string(token)
	}
	return out
}

func TestTokenizeSeparatorsAndWhitespace(t *testing.T) {
	t.Parallel()

	testCases := []tokenizeTestCase{
		{
			name:     "basic declaration",
			input:    "int x = 1;",
			expected: []string{"int", " ", "x", " ", "=", " ", "1", ";"},
		},
		{
			name:     "parentheses braces and colon",
			input:    "switch(x){case 1:return y;}",
			expected: []string{"switch", "(", "x", ")", "{", "case", " ", "1", ":", "return", " ", "y", ";", "}"},
		},
		{
			name:     "tabs and newlines preserved",
			input:    "if\t(x)\n{\n}",
			expected: []string{"if", "\t", "(", "x", ")", "\n", "{", "\n", "}"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertTokens(t, tc.input, tc.expected)
		})
	}
}

func TestTokenizeOperators(t *testing.T) {
	t.Parallel()

	testCases := []tokenizeTestCase{
		{
			name:     "increment and decrement",
			input:    "i++ + j--;",
			expected: []string{"i", "++", " ", "+", " ", "j", "--", ";"},
		},
		{
			name:     "compound assignment and shifts",
			input:    "value<<=2; mask>>=1;",
			expected: []string{"value", "<<=", "2", ";", " ", "mask", ">>=", "1", ";"},
		},
		{
			name:     "comparison and logical operators",
			input:    "a<=b && c>=d || e==f",
			expected: []string{"a", "<=", "b", " ", "&&", " ", "c", ">=", "d", " ", "||", " ", "e", "==", "f"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertTokens(t, tc.input, tc.expected)
		})
	}
}

func TestTokenizeStringsAndChars(t *testing.T) {
	t.Parallel()

	testCases := []tokenizeTestCase{
		{
			name:     "double quoted string",
			input:    `printf("for (;;) {}");`,
			expected: []string{"printf", "(", `"for (;;) {}"`, ")", ";"},
		},
		{
			name:     "escaped quote in string",
			input:    `printf("\"for\"");`,
			expected: []string{"printf", "(", `"\"for\""`, ")", ";"},
		},
		{
			name:     "char literal with escape",
			input:    `char quote = '\'';`,
			expected: []string{"char", " ", "quote", " ", "=", " ", `'\''`, ";"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertTokens(t, tc.input, tc.expected)
		})
	}
}

func TestTokenizeComments(t *testing.T) {
	t.Parallel()

	testCases := []tokenizeTestCase{
		{
			name:     "single line comment",
			input:    "int x; // comment\nreturn x;",
			expected: []string{"int", " ", "x", ";", " ", "// comment", "\n", "return", " ", "x", ";"},
		},
		{
			name:     "multi line comment",
			input:    "int x; /* line 1\nline 2 */ return x;",
			expected: []string{"int", " ", "x", ";", " ", "/* line 1\nline 2 */", " ", "return", " ", "x", ";"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertTokens(t, tc.input, tc.expected)
		})
	}
}

func TestTokenizePreservesTrailingTokenAtEOF(t *testing.T) {
	t.Parallel()

	assertTokens(t, "sizeof *ptr", []string{"sizeof", " ", "*", "ptr"})
}
