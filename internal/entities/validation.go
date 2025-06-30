package entities

import "regexp"

func isEmail(str string) bool {
	m, err := regexp.Match("^[а-яА-Яa-zA-Z0-9._+-]+@([a-zа-я-]+\\.)+[a-zа-я-]{2,4}$", []byte(str))
	if err != nil {
		return false
	}
	return m
}

func stringLengthInRange(str string, min, max int) bool {
	return len(str) >= min && len(str) <= max
}

func isLatinLetter(symbol rune) bool {
	return (symbol >= 'a' && symbol <= 'z') || (symbol >= 'A' && symbol <= 'Z')
}

func isDigit(symbol rune) bool {
	return symbol >= '0' && symbol <= '9'
}
