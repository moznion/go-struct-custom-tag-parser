package tagparser

import (
	"errors"
	"unicode"
)

const valueQuote = '"'
const keyValueDelimiter = ':'

// Parse parses a custom tag string.
//
// Strict mode is a mode selector:
// It raises an error when given unacceptable custom tag string when the mode is true.
// On the other hand, if the mode is false, it immediately returns the processed results until just before the invalid custom tag syntax.
func Parse(tagString string, isStrict bool) (map[string]string, error) {
	key := make([]rune, 0, 100)
	keyCursor := 0
	value := make([]rune, 0, 100)
	valueCursor := 0

	inKeyParsing := true
	isEscaping := false

	tagKeyValue := make(map[string]string)

	tagRunes := []rune(tagString)
	tagRunesLen := len(tagRunes)
	for i := 0; i < tagRunesLen; i++ {
		r := tagRunes[i]

		if inKeyParsing {
			if unicode.IsSpace(r) {
				if keyCursor > 0 {
					if isStrict {
						return nil, errors.New("invalid custom tag syntax: key must not contain any white space, but it contains")
					}
					break
				}
				continue
			}

			if r == valueQuote {
				if isStrict {
					return nil, errors.New("invalid custom tag syntax: key must not contain any double quote, but it contains")
				}
				// give up when key contains any double quote
				break
			}

			if r == keyValueDelimiter {
				if keyCursor <= 0 {
					if isStrict {
						return nil, errors.New("invalid custom tag syntax: key must not be empty, but it gets empty")
					}

					// if empty key has come, it should give up
					break
				}

				inKeyParsing = false
				i++
				if i >= tagRunesLen {
					if isStrict {
						return nil, errors.New("invalid custom tag syntax: value must not be empty, but it gets empty")
					}
					break
				}
				if tagRunes[i] != valueQuote {
					if isStrict {
						return nil, errors.New("invalid custom tag syntax: quote for value is missing")
					}
					break
				}
				continue
			}
			key = append(key, r)
			keyCursor++
			continue
		}

		// value parsing
		if !isEscaping && r == valueQuote {
			tagKeyValue[string(key[:keyCursor])] = string(value[:valueCursor])
			key = key[:0]
			keyCursor = 0
			value = value[:0]
			valueCursor = 0
			inKeyParsing = true
			continue
		}

		if r == '\\' {
			if isEscaping {
				value = append(value, r)
				valueCursor++
				isEscaping = false
				continue
			}
			isEscaping = true
			continue
		}
		value = append(value, r)
		isEscaping = false
		valueCursor++
	}

	if inKeyParsing && keyCursor > 0 && isStrict {
		return nil, errors.New("invalid custom tag syntax: a delimiter of key and value is missing")
	}

	if !inKeyParsing && valueCursor > 0 && isStrict {
		return nil, errors.New("invalid custom tag syntax: a value is not terminated with quote")
	}

	return tagKeyValue, nil
}
