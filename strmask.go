// Package strmask provides two simple functions to validate and format a string based on a mask.
package strmask

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// ValidateAndFormatMask formats the string s based on the mask, reporting any invalid runes that doesn't fit the provided mask. Information about the mask symbols can be found at: https://github.com/frones/strmask
func ValidateAndFormatMask(mask string, s string) (string, error) {
	printNext := false
	output := ""
	strOffset := 0
	lastErrOffset := -1
	err := ""
	for _, r := range mask {
		if printNext {
			output += string(r)
			printNext = false

			r2, w := utf8.DecodeRuneInString(s[strOffset:])
			if r2 == r {
				strOffset += w
			}

			continue
		}

		switch r {
		case '\\':
			printNext = true
		case '9', '0':
			r2, w := utf8.DecodeRuneInString(s[strOffset:])
			if unicode.IsDigit(r2) {
				output += string(r2)
				strOffset += w
			} else if r == '0' {
				if lastErrOffset < strOffset {
					err += fmt.Sprintf("invalid character \"%s\" (expected a digit) at position %d\n", string(r2), strOffset)
					lastErrOffset = strOffset
				}
				output += " "
			}
		case 'L', 'l':
			r2, w := utf8.DecodeRuneInString(s[strOffset:])
			if unicode.IsLetter(r2) {
				output += string(r2)
				strOffset += w
			} else if r == 'L' {
				if lastErrOffset < strOffset {
					err += fmt.Sprintf("invalid character \"%s\" (expected a letter) at position %d\n", string(r2), strOffset)
					lastErrOffset = strOffset
				}
				output += " "
			}
		case 'A', 'a':
			r2, w := utf8.DecodeRuneInString(s[strOffset:])
			if unicode.IsLetter(r2) || unicode.IsDigit(r2) {
				output += string(r2)
				strOffset += w
			} else if r == 'A' {
				if lastErrOffset < strOffset {
					err += fmt.Sprintf("invalid character \"%s\" (expected a letter or digit) at position %d\n", string(r2), strOffset)
					lastErrOffset = strOffset
				}
				output += " "
			}
		case 'U', 'u':
			r2, w := utf8.DecodeRuneInString(s[strOffset:])
			if unicode.IsLetter(r2) {
				output += string(unicode.ToUpper(r2))
				strOffset += w
			} else if r == 'U' {
				if lastErrOffset < strOffset {
					err += fmt.Sprintf("invalid character \"%s\" (expected a letter) at position %d\n", string(r2), strOffset)
					lastErrOffset = strOffset
				}
				output += " "
			}
		case 'W', 'w':
			r2, w := utf8.DecodeRuneInString(s[strOffset:])
			if unicode.IsLetter(r2) {
				output += string(unicode.ToLower(r2))
				strOffset += w
			} else if r == 'W' {
				if lastErrOffset < strOffset {
					err += fmt.Sprintf("invalid character \"%s\" (expected a letter) at position %d\n", string(r2), strOffset)
					lastErrOffset = strOffset
				}
				output += " "
			}
		default:
			output += string(r)
			r2, w := utf8.DecodeRuneInString(s[strOffset:])
			if r2 == r {
				strOffset += w
			}
		}
	}

	if len(s) > strOffset {
		output += s[strOffset:]
	}

	if err != "" {
		return output, fmt.Errorf(err)
	} else {
		return output, nil
	}
}

// FormatMask is a single-return value helper function to the ValidateAndFormatMask function. It just ignores any errors and outputs the formatted string returned.
func FormatMask(mask string, s string) string {
	output, _ := ValidateAndFormatMask(mask, s)
	return output
}
