package internal

import (
	"log"
	"strings"
)

// ParseAndApply parses the passed format string and applies styles, returning a styled string.
func ParseAndApply(text string, disable bool) string {
	var resParts = make([]string, 0, 50)

	var err error

	var tempTokenStartIndex int
	var tempTokenEndIndex int

	var lastToken string

	var inText bool
	var inFormat bool
	var inFormatGroup bool

	var lastIsFormatGroup bool
	var needCheckFormatGroupStyle bool

	var startBracketFound bool
	var endBracketFound bool
	var colonFound bool

	inText = true
	formatLen := len(text)

	for index := 0; index < formatLen; index++ {
		s := text[index]

		// if just text
		if inText {
			// if found {
			if startBracketFound {
				if s == '{' {
					// index before {{
					tempTokenEndIndex = index - 1

					// last part before {{
					resParts = append(resParts, text[tempTokenStartIndex:tempTokenEndIndex])

					// index after {{
					tempTokenStartIndex = index + 1

					// set current state
					inFormatGroup = true
					inFormat = false
					inText = false
				}
				startBracketFound = false
				continue
			}

			if s == '{' {
				startBracketFound = true
				continue
			}

			// if after {{}} no ::style
			if needCheckFormatGroupStyle {
				needCheckFormatGroupStyle = false

				if !(s == ':' && index+1 < formatLen && text[index+1] == ':' && index+2 < formatLen && text[index+2] != ' ') {
					lastIsFormatGroup = false

					if lastToken != "" {
						resParts = append(resParts, "{{"+lastToken+"}}")
						tempTokenStartIndex = index
					}
				}
			}

			// if found :
			if colonFound {
				if s == ':' {
					// index after ::
					tempTokenStartIndex = index + 1

					// set current state
					inFormatGroup = false
					inFormat = true
					inText = false
				}
				colonFound = false
				continue
			}

			if s == ':' && lastIsFormatGroup {
				colonFound = true
				continue
			}

			continue
		}

		// if in {{}}
		if inFormatGroup {
			// if found }
			if endBracketFound {
				if s == '}' {
					// index before }}
					tempTokenEndIndex = index - 1

					// part in {{}}
					lastToken = text[tempTokenStartIndex:tempTokenEndIndex]

					// set current state
					inFormatGroup = false
					inFormat = false
					inText = true

					lastIsFormatGroup = true
					needCheckFormatGroupStyle = true
				}
				endBracketFound = false
				continue
			}

			if s == '}' {
				countBracketInRow := 0
				for i := 0; ; i++ {
					if index+i < formatLen && text[index+i] != '}' {
						break
					}
					countBracketInRow++
				}

				if countBracketInRow > 2 {
					index += countBracketInRow - 2 - 1
					continue
				}

				endBracketFound = true
				continue
			}
			continue
		}

		// if after {{}}
		if inFormat {
			if s == '|' {
				// index before |
				tempTokenEndIndex = index

				// text from :: or previous | to current |
				singleFormat := text[tempTokenStartIndex:tempTokenEndIndex]
				lastToken, err = applyStyle(lastToken, singleFormat, disable)
				if err != nil {
					log.Fatalf("Error parse style string in '%s' text string: %v", strings.ReplaceAll(text, "\n", "\\n"), err)
				}

				// index after |
				tempTokenStartIndex = index + 1
				continue
			}

			if !(s >= 'a' && s <= 'z' || s >= 'A' && s <= 'Z') && s != '|' && !(s >= '0' && s <= '9') && s != '#' {
				tempTokenEndIndex = index

				// last format
				singleFormat := text[tempTokenStartIndex:tempTokenEndIndex]
				lastToken, err = applyStyle(lastToken, singleFormat, disable)
				if err != nil {
					log.Fatalf("Error parse style string in '%s' text string: %v", strings.ReplaceAll(text, "\n", "\\n"), err)
				}

				tempTokenStartIndex = index

				resParts = append(resParts, lastToken)

				// set current state
				inFormatGroup = false
				inFormat = false
				inText = true

				lastIsFormatGroup = false

				if s == '{' {
					index--
				}

				continue
			}
			continue
		}
	}

	if lastIsFormatGroup {
		singleFormat := text[tempTokenStartIndex:]
		lastToken, err = applyStyle(lastToken, singleFormat, disable)
		if err != nil {
			log.Fatalf("Error parse style string in '%s' text string: %v", strings.ReplaceAll(text, "\n", "\\n"), err)
		}
		resParts = append(resParts, lastToken)
	} else {
		resParts = append(resParts, text[tempTokenStartIndex:])
	}

	return strings.Join(resParts, "")
}
