package main

type stateFn func(*scanner) stateFn

type Token []byte

type scanner struct {
	sequence []byte
	tokens   []Token

	tokStart int
	i        int // Index into sequence

	state stateFn
}

func newscanner(bytes []byte) *scanner {
	s := &scanner{sequence: bytes}
	s.tokStart = 0
	s.i = 0
	s.state = stateRegularCode
	return s
}

func (s *scanner) finishToken() {
	if s.tokStart == s.i {
		return
	}
	token := s.sequence[s.tokStart:s.i]
	s.tokens = append(s.tokens, token)
	s.tokStart = s.i
}

func (s *scanner) extendToken() {
	s.i++
}

func (s *scanner) getCurrentByte() byte {
	return s.sequence[s.i]
}

func (s *scanner) getPreviousByte() byte {
	if s.i == 0 {
		panic("Programming error: cannot call scanner.getPreviousByte() when i == 0")
	}
	return s.sequence[s.i-1]
}

func (s *scanner) isDone() bool {
	return s.i == len(s.sequence)
}

func Tokenize(bytes []byte) []Token {
	scanner := newscanner(bytes)
	for !scanner.isDone() {
		scanner.state = scanner.state(scanner)
	}
	return scanner.tokens
}

func stateRegularCode(s *scanner) stateFn {
	c := s.getCurrentByte()
	switch c {
	case '\n', '\r', '\t', ' ':
		s.finishToken()
		s.extendToken()
		s.finishToken()
		return stateRegularCode

	case '\'':
		s.finishToken()
		s.extendToken()
		return stateSingleQuoteString

	case '"':
		s.finishToken()
		s.extendToken()
		return stateDoubleQuoteString

	case '{', '}', '(', ')', '=', ',', ';':
		s.finishToken()
		s.extendToken()
		s.finishToken()
		return stateRegularCode

	case '+', '-', '/', '*', '|', '&', '^', '~', '<', '>':
		s.finishToken()
		s.extendToken()
		return statePossiblyMultiCharOperator

	default:
		s.extendToken()
		return stateRegularCode
	}
}

func stateEscapedSingleQuoteString(s *scanner) stateFn {
	s.extendToken()
	return stateSingleQuoteString
}

func stateEscapedDoubleQuoteString(s *scanner) stateFn {
	s.extendToken()
	return stateDoubleQuoteString
}

func stateSingleQuoteString(s *scanner) stateFn {
	c := s.getCurrentByte()
	switch c {
	case '\\':
		s.extendToken()
		return stateEscapedSingleQuoteString
	case '\'':
		s.extendToken()
		s.finishToken()
		return stateRegularCode
	default:
		s.extendToken()
		return stateSingleQuoteString
	}
}

func stateDoubleQuoteString(s *scanner) stateFn {
	c := s.getCurrentByte()
	switch c {
	case '\\':
		s.extendToken()
		return stateEscapedDoubleQuoteString
	case '"':
		s.extendToken()
		s.finishToken()
		return stateRegularCode
	default:
		s.extendToken()
		return stateDoubleQuoteString
	}
}

func stateSingleLineComment(s *scanner) stateFn {
	c := s.getCurrentByte()
	if c == '\n' || c == '\r' {
		s.finishToken()
		s.extendToken()
		s.finishToken()
		return stateRegularCode
	}
	s.extendToken()
	return stateSingleLineComment
}

func stateMultiLineComment(s *scanner) stateFn {
	c := s.getCurrentByte()
	prev := s.getPreviousByte()
	s.extendToken()

	if prev == '*' && c == '/' {
		s.finishToken()
		return stateRegularCode
	}
	return stateMultiLineComment
}

func statePossiblyMultiCharOperator(s *scanner) stateFn {
	c := s.getCurrentByte()
	prev := s.getPreviousByte()

	switch {
	case c == '=':
		// +=, |=, <<=, etc.
		s.extendToken()
		s.finishToken()
		return stateRegularCode

	case prev == '/' && c == '/':
		s.extendToken()
		return stateSingleLineComment

	case prev == '/' && c == '*':
		s.extendToken()
		return stateMultiLineComment

	case prev == '+' && c == '+':
		s.extendToken()
		s.finishToken()
		return stateRegularCode

	case prev == '-' && c == '-':
		s.extendToken()
		s.finishToken()
		return stateRegularCode

	case prev == '<' && c == '<':
		s.extendToken()
		return statePossiblyMultiCharOperator

	case prev == '>' && c == '>':
		s.extendToken()
		return statePossiblyMultiCharOperator

	default:
		s.finishToken()
		return stateRegularCode
	}
}
