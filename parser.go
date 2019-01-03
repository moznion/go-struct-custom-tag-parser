package tagparser

import (
	"errors"
	"unicode"
)

// Parse parses a custom tag string.
func Parse(tagString string) (map[string]string, error) {
	key := make([]byte, 0, 100)
	keyCursor := 0
	value := make([]byte, 0, 100)
	valueCursor := 0
	inKeyParsing := true
	isEscaping := false
	justSplited := false
	var valueTerminator rune

	tagKeyValue := make(map[string]string)

	bs := []byte(tagString)
	previousIndex := -1
	i := 0

	for pos, r := range tagString {
		if justSplited {
			valueTerminator = r

			i += pos - previousIndex
			if i >= len(bs) {
				return nil, errors.New("invalid custom tag syntax: value must not be empty, but it gets empty")
			}
			previousIndex = pos

			inKeyParsing = false
			justSplited = false

			continue
		}

		if inKeyParsing {
			if unicode.IsSpace(r) {
				if keyCursor > 0 {
					return nil, errors.New("invalid custom tag syntax: key must not contain any white space, but it contains")
				}
				i += pos - previousIndex
				previousIndex = pos
				continue
			}

			if r == ':' {
				if keyCursor <= 0 {
					return nil, errors.New("invalid custom tag syntax: key must not be empty, but it gets empty")
				}

				justSplited = true

				i += pos - previousIndex
				previousIndex = pos

				continue
			}

			for j := 0; j < (pos - previousIndex); i, j = i+1, j+1 {
				key = append(key, bs[i])
			}

			keyCursor++
			previousIndex = pos
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

			i += pos - previousIndex
			previousIndex = pos
			continue
		}

		if r == '\\' {
			isEscaping = true
			i += pos - previousIndex
			previousIndex = pos
			continue
		}

		for j := 0; j < (pos - previousIndex); i, j = i+1, j+1 {
			value = append(value, bs[i])
		}

		isEscaping = false
		valueCursor++
		previousIndex = pos
	}

	if inKeyParsing && keyCursor > 0 {
		return nil, errors.New("invalid custom tag syntax: a delimiter of key and value is missing")
	}

	if !inKeyParsing && valueCursor > 0 {
		return nil, errors.New("invalid custom tag syntax: a value is not terminated with quote")
	}

	return tagKeyValue, nil

}
