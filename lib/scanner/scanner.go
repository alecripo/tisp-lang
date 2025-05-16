package scanner

import "fmt"

type TokenType uint8

const (
	// Single character tokens
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	PIPE
	DOUBLE_QUOTES
	QUESTION_MARK
	EXCLAMATION_MARK
	DOUBLE_COLON

	// Literals
	STRING
	INTEGER

	IDENTIFIER

	// Keywords
	IF
	AND
	OR

	EOF
)

type Token struct {
	Type         TokenType
	Line, Column uint8
	Lexeme       string
}

func (t Token) String() string {
	typeNames := make(map[TokenType]string)
	typeNames[LEFT_PAREN] = "LEFT_PAREN"
	typeNames[RIGHT_PAREN] = "RIGHT_PAREN"
	typeNames[PIPE] = "PIPE"
	typeNames[DOUBLE_QUOTES] = "DOUBLE_QUOTES"
	typeNames[QUESTION_MARK] = "QUESTION_MARK"
	typeNames[EXCLAMATION_MARK] = "EXCLAMATION_MARK"
	typeNames[DOUBLE_COLON] = "DOUBLE_COLON"
	typeNames[STRING] = "STRING"
	typeNames[INTEGER] = "INTEGER"
	typeNames[IDENTIFIER] = "IDENTIFIER"
	typeNames[IF] = "IF"
	typeNames[AND] = "AND"
	typeNames[OR] = "OR"
	typeNames[EOF] = "EOF"

	return fmt.Sprintf("%s %s\n", typeNames[t.Type], t.Lexeme)
}

type ParseError struct {
	Line  uint8  `json:"line"`
	Cause string `json:"cause"`
}

func (e ParseError) Error() string {
	return fmt.Sprintf(
		"[line %d] %s",
		e.Line,
		e.Cause,
	)
}

type Scanner struct {
	Source       string
	line, column uint8
	lastIndex    uint16
	currIndex    uint16
}

func (s Scanner) Scan() ([]Token, error) {
	s.currIndex = 0
	s.line = 1
	tokens := make([]Token, 0)
	errors := make([]ParseError, 0)
	for !s.isEnd() {
		s.lastIndex = s.currIndex
		token, err := s.scanToken()
		if err != nil {
			errors = append(errors, ParseError{Line: s.line, Cause: err.Error()})
			continue
		}
		tokens = append(tokens, token)
	}
	tokens = append(tokens, Token{Type: EOF})
	return tokens, nil
}

func (s *Scanner) scanToken() (t Token, err error) {
	switch c := s.advance(); c {
	case '|':
		return s.newToken(PIPE), nil
	case '(':
		return s.newToken(LEFT_PAREN), nil
	case ')':
		return s.newToken(RIGHT_PAREN), nil
	case '?':
		return s.newToken(QUESTION_MARK), nil
	case '!':
		return s.newToken(EXCLAMATION_MARK), nil
	case ':':
		return s.newToken(DOUBLE_COLON), nil
	case '"':
		return s.scanString()
	case ' ', '\t', '\r':
		return s.scanToken()
	case '\n':
		s.line++
		s.column = 0
		return s.scanToken()
	default:
		if isInt(c) {
			return s.scanInt()
		}
		if isAlpha(c) {
			return s.scanIdentifier()
		}
		return Token{}, fmt.Errorf("Unexpected character %s", string(c))
	}
}

// Returns the current character and advances the index.
func (s *Scanner) advance() byte {
	s.currIndex++
	s.column++
	return s.Source[s.currIndex-1]
}

// Advances the current index and character under the new index with the given one.
func (s *Scanner) isNext(c byte) bool {
	if s.isEnd() {
		return false
	}

	return s.advance() == c
}

func (s Scanner) peek() byte {
	return s.Source[s.currIndex]
}

// Builds a new token of the given type from current scanner state
func (s Scanner) newToken(t TokenType) Token {
	return Token{
		Type:   t,
		Line:   s.line,
		Column: s.column,
		Lexeme: s.Source[s.lastIndex:s.currIndex],
	}
}

func (s *Scanner) scanInt() (t Token, err error) {
	for next := s.peek(); isInt(next); next = s.advance() {
	}
	return s.newToken(INTEGER), nil
}
func (s *Scanner) scanString() (t Token, err error) {
	for next := s.peek(); next != '"' && !s.isEnd(); next = s.peek() {
		if next == '\n' {
			s.line++
			s.column = 0
		}
		s.advance()
	}

	if s.isEnd() {
		return Token{}, fmt.Errorf("unterminated string literal.")
	}

	s.advance() // consume closing "

	return Token{
		Type:   STRING,
		Line:   s.line,
		Column: s.column,
		Lexeme: s.Source[s.lastIndex+1 : s.currIndex-1],
	}, nil
}

func (s *Scanner) scanIdentifier() (t Token, err error) {
	keywords := make(map[string]TokenType)
	keywords["if"] = IF
	keywords["or"] = OR
	keywords["and"] = AND
	for next := s.peek(); isAlpha(next); next = s.peek() {
		s.advance()
	}
	if ktype, ok := keywords[s.Source[s.lastIndex:s.currIndex]]; ok {
		return s.newToken(ktype), nil
	}
	return s.newToken(IDENTIFIER), nil
}

// Are we at the end of the source?
func (s Scanner) isEnd() bool {
	return int(s.currIndex) >= len(s.Source)
}

func isInt(c byte) bool {
	return (c >= '0' && c <= '9')
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
