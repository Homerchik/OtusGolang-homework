package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func processShieldCase(s string, startIdx int, result *strings.Builder) (int, error) {
	nextRune := rune(s[startIdx+1])
	if startIdx+2 < len(s) {
		afterNextRune := rune(s[startIdx+2])
		switch {
		case unicode.IsDigit(afterNextRune):
			result.WriteString(strings.Repeat(string(nextRune), int(afterNextRune-'0')))
			return 3, nil
		case unicode.IsLetter(afterNextRune) || afterNextRune == '\\':
			result.WriteString(string(nextRune))
			return 2, nil
		default:
			return 0, ErrInvalidString
		}
	} else {
		result.WriteString(string(nextRune))
		return 2, nil
	}
}

func Unpack(s string) (string, error) {
	fstCharIdx := 0
	result := strings.Builder{}
	for {
		if fstCharIdx > len(s)-2 {
			break
		}
		fstRune, nextRune := rune(s[fstCharIdx]), rune(s[fstCharIdx+1])
		increment := 0
		switch {
		case fstRune == '\\' && (nextRune == '\\' || unicode.IsDigit(nextRune)):
			if inc, err := processShieldCase(s, fstCharIdx, &result); err != nil {
				return "", err
			} else {
				increment = inc
			}
		case fstRune == '\\' && !(nextRune == '\\' || unicode.IsDigit(nextRune)):
			return "", ErrInvalidString
		case unicode.IsLetter(fstRune) && !unicode.IsDigit(nextRune):
			result.WriteRune(fstRune)
			increment = 1
		case unicode.IsLetter(fstRune) && unicode.IsDigit(nextRune):
			increment = 2
			result.WriteString(strings.Repeat(string(fstRune), int(nextRune-'0')))
		default:
			return "", ErrInvalidString
		}
		fstCharIdx += increment
	}
	if fstCharIdx < len(s) {
		lastRune := rune(s[len(s)-1])
		if lastRune == '\\' || !unicode.IsLetter(lastRune) {
			return "", ErrInvalidString
		}
		result.WriteByte(s[fstCharIdx])
	}
	return result.String(), nil
}
