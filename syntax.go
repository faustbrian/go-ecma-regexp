package ecmascript

import "unicode/utf8"

// TokenKind classifies lexical pattern input. Tokens retain source byte spans.
type TokenKind uint8

const (
	// TokenCharacter retains a non-punctuation pattern character.
	TokenCharacter TokenKind = iota + 1
	// TokenEscape retains a backslash and the next decoded character for contextual
	// parsing.
	TokenEscape
	// TokenAlternation marks the alternative separator.
	TokenAlternation
	// TokenDot marks dot syntax.
	TokenDot
	// TokenCaret marks start-assertion syntax outside a class.
	TokenCaret
	// TokenDollar marks end-assertion syntax outside a class.
	TokenDollar
	// TokenLeftParen opens a group or assertion.
	TokenLeftParen
	// TokenRightParen closes a group or assertion.
	TokenRightParen
	// TokenLeftBracket opens class syntax.
	TokenLeftBracket
	// TokenRightBracket closes class syntax or is a legacy literal by context.
	TokenRightBracket
	// TokenLeftBrace begins a braced quantifier or is interpreted by legacy grammar.
	TokenLeftBrace
	// TokenRightBrace ends a braced quantifier or is interpreted by legacy grammar.
	TokenRightBrace
	// TokenStar marks zero-or-more repetition.
	TokenStar
	// TokenPlus marks one-or-more repetition.
	TokenPlus
	// TokenQuestion marks optional repetition, a lazy modifier or a group extension by
	// context.
	TokenQuestion
	// TokenComma is interpreted in braced counts or as a literal by context.
	TokenComma
	// TokenEOF marks the empty end-of-source span.
	TokenEOF
)

// Token is an immutable lexical token.
type Token struct {
	kind TokenKind
	text string
	span Span
}

// Kind reports lexical punctuation/character classification, not a parsed grammar
// production.
func (t Token) Kind() TokenKind { return t.kind }

// Text returns retained original source bytes; parser decoding can interpret invalid
// UTF-8 as U+FFFD.
func (t Token) Text() string { return t.text }

// Span returns the token byte range; EOF has an empty span at source end.
func (t Token) Span() Span { return t.span }

// Tokenize lexes pattern while retaining exact byte spans. Escapes are kept as
// two-code-point lexical units so the parser can interpret them by edition and
// grammar context. Only edition and PatternBytes are checked here; structural
// AST allowances apply in Parse. The returned slice is caller-owned. No context,
// normalization or JavaScript literal-delimiter parsing is supplied.
func Tokenize(pattern string, options ParseOptions) ([]Token, error) {
	if err := options.validate(pattern); err != nil {
		return nil, err
	}

	tokens := make([]Token, 0, min(len(pattern)+1, 256))
	for offset := 0; offset < len(pattern); {
		start := offset
		char, size := utf8.DecodeRuneInString(pattern[offset:])
		offset += size
		if char == '\\' {
			if offset == len(pattern) {
				return nil, &SyntaxError{Code: SyntaxUnexpectedEOF, Span: Span{Start: start, End: offset}, Message: "escape has no escaped character"}
			}
			_, escapedSize := utf8.DecodeRuneInString(pattern[offset:])
			offset += escapedSize
			tokens = append(tokens, Token{kind: TokenEscape, text: pattern[start:offset], span: Span{Start: start, End: offset}})
			continue
		}

		tokens = append(tokens, Token{kind: tokenKind(char), text: pattern[start:offset], span: Span{Start: start, End: offset}})
	}
	tokens = append(tokens, Token{kind: TokenEOF, span: Span{Start: len(pattern), End: len(pattern)}})

	return tokens, nil
}

func tokenKind(char rune) TokenKind {
	switch char {
	case '|':
		return TokenAlternation
	case '.':
		return TokenDot
	case '^':
		return TokenCaret
	case '$':
		return TokenDollar
	case '(':
		return TokenLeftParen
	case ')':
		return TokenRightParen
	case '[':
		return TokenLeftBracket
	case ']':
		return TokenRightBracket
	case '{':
		return TokenLeftBrace
	case '}':
		return TokenRightBrace
	case '*':
		return TokenStar
	case '+':
		return TokenPlus
	case '?':
		return TokenQuestion
	case ',':
		return TokenComma
	default:
		return TokenCharacter
	}
}
