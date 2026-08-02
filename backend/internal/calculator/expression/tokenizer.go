package expression

import (
	"errors"
	"strings"
	"unicode"
)

func Tokenize(input string) ([]Token, error) {
	characters := []rune(input)
	tokens := make([]Token, 0)

	for position := 0; position < len(characters); {
		character := characters[position]

		switch {
		case unicode.IsSpace(character):
			position++
		case unicode.IsDigit(character) || character == '.':
			token, nextPosition, err := scanNumber(characters, position)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
			position = nextPosition
		case isIdentifierStart(character):
			token, nextPosition, err := scanIdentifier(characters, position)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
			position = nextPosition
		case isOperator(character):
			tokens = append(tokens, Token{
				Type:     TokenOperator,
				Value:    normalizeOperator(character),
				Position: position,
			})
			position++
		case character == '(':
			tokens = append(tokens, Token{
				Type:     TokenLeftParenthesis,
				Value:    string(character),
				Position: position,
			})
			position++
		case character == ')':
			tokens = append(tokens, Token{
				Type:     TokenRightParenthesis,
				Value:    string(character),
				Position: position,
			})
			position++
		case character == '%':
			tokens = append(tokens, Token{
				Type:     TokenPercentage,
				Value:    "%",
				Position: position,
			})
			position++
		default:
			return nil, errors.New("invalid character found in input")
		}
	}

	if len(tokens) == 0 {
		return nil, errors.New("empty expression")
	}

	return tokens, nil
}

// scanNumber scans a number from the input
func scanNumber(characters []rune, startPosition int) (Token, int, error) {
	position := startPosition
	decimalPointFound := false
	digitCount := 0

	for position < len(characters) {
		character := characters[position]

		switch {
		case isDigit(character):
			digitCount++
			position++
		case isDecimalPoint(character):
			if decimalPointFound {
				return Token{}, position, errors.New("multiple decimal points found")
			}
			decimalPointFound = true
			position++
		default:
			goto numberComplete
		}
	}

numberComplete:
	if digitCount == 0 {
		return Token{}, 0, errors.New("number must start with a digit")
	}

	value := string(characters[startPosition:position])

	return Token{
		Type:     TokenNumber,
		Value:    value,
		Position: startPosition,
	}, position, nil
}

// scanIdentifier scans an identifier from the input
func scanIdentifier(characters []rune, startPosition int) (Token, int, error) {
	position := startPosition
	for position < len(characters) &&
		isIdentifierPart(characters[position]) {
		position++
	}

	value := strings.ToLower(
		string(characters[startPosition:position]),
	)

	return Token{
		Type:     TokenFunction,
		Value:    value,
		Position: startPosition,
	}, position, nil
}

// isDigit checks if the character is a digit
func isDigit(character rune) bool {
	return character >= '0' && character <= '9'
}

// isDecimalPoint checks if the character is a decimal point
func isDecimalPoint(character rune) bool {
	return character == '.'
}

// isIdentifierStart checks if the character is the start of an identifier
func isIdentifierStart(character rune) bool {
	return unicode.IsLetter(character) || character == '_'
}

// isIdentifierPart checks if the character is a part of an identifier
func isIdentifierPart(character rune) bool {
	return unicode.IsLetter(character) || unicode.IsDigit(character) || character == '_'
}

// isOperator checks if the character is an operator
func isOperator(character rune) bool {
	switch character {
	case '+', '-', '*', '/', '^', '×', '÷', '−':
		return true

	default:
		return false
	}
}

// normalizeOperator normalizes the operator to a standard form
func normalizeOperator(character rune) string {
	switch character {
	case '×':
		return "*"

	case '÷':
		return "/"

	case '−':
		return "-"

	default:
		return string(character)
	}
}
