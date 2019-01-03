package tagparser

import (
	"errors"
	"unicode"
)

// Parse parses a custom tag string.
func Parse(tagString string) (map[string]string, error) {
	return parse(tagString, false)
}

func parse(tagString string, isStrict bool) (map[string]string, error) {
	key := make([]rune, 0, 100)
	keyCursor := 0
	value := make([]rune, 0, 100)
	valueCursor := 0

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
				if keyCursor > 0 {
					return nil, errors.New("invalid custom tag syntax: key must not contain any white space, but it contains")
				}
				continue
			}

			if r == ':' {
				if keyCursor <= 0 {
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
			keyCursor++
			continue
		}

		// value parsing
		if !isEscaping && r == valueTerminator {
			tagKeyValue[string(key[:keyCursor])] = string(value[:valueCursor])
			key = key[:0]
			keyCursor = 0
			value = value[:0]
			valueCursor = 0
			inKeyParsing = true
			continue
		}

		if r == '\\' {
			isEscaping = true
			continue
		}
		value = append(value, r)
		isEscaping = false
		valueCursor++
	}

	if inKeyParsing && keyCursor > 0 {
		return nil, errors.New("invalid custom tag syntax: a delimiter of key and value is missing")
	}

	if !inKeyParsing && valueCursor > 0 {
		return nil, errors.New("invalid custom tag syntax: a value is not terminated with quote")
	}

	return tagKeyValue, nil

}
