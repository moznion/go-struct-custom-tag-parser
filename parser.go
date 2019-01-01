package tagparser

import (
	"errors"
	"unicode"
)

// Parse parses a custom tag string.
func Parse(tagString string) (map[string]string, error) {
	key := make([]rune, 0)
	value := make([]rune, 0)
	inKeyParsing := true
	isEscaping := false
	var valueTerminator rune

	tagKeyValue := make(map[string]string)

	tagRunes := []rune(tagString)
	tagRunesLen := len(tagRunes)
	for i := 0; i < tagRunesLen; i++ {
		r := tagRunes[i]

		if inKeyParsing {
			if unicode.IsSpace(r) {
				continue
			}

			if r == ':' {
				if len(key) <= 0 {
					return nil, errors.New("invalid custom tag syntax: key must not be empty, but it gets empty")
				}

				inKeyParsing = false
				i++
				if i >= tagRunesLen {
					return nil, errors.New("invalid custom tag syntax: value must not be empty, but it gets empty")
				}
				valueTerminator = tagRunes[i]
				continue
			}
			key = append(key, r)
			continue
		}

		// value parsing
		if !isEscaping && r == valueTerminator {
			tagKeyValue[string(key)] = string(value)
			key = make([]rune, 0)
			value = make([]rune, 0)
			inKeyParsing = true
			continue
		}

		if r == '\\' {
			isEscaping = true
			continue
		}
		value = append(value, r)
		isEscaping = false
	}

	if inKeyParsing && len(key) > 0 {
		return nil, errors.New("invalid custom tag syntax: a delimiter of key and value is missing")
	}

	if !inKeyParsing && len(value) > 0 {
		return nil, errors.New("invalid custom tag syntax: a value is not terminated with quote")
	}

	return tagKeyValue, nil
}
