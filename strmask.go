// Package strmask provides two simple functions to validate and format a string based on a mask.
package strmask

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

var asciiLetters = &unicode.RangeTable{
	R16: []unicode.Range16{
		{'A', 'Z', 1},
		{'a', 'z', 1},
	},
}

// Reverses the provided string. Since this functions reverses the string rune by rune, it has problems with unicode combining characters.
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// This functions just reverses any backslashes in the string with the rune preceding it. It's used to preserve escaped characters in reversed strings.
func reverseBackslashs(s string) string {
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' && i > 0 {
			runes[i-1], runes[i] = runes[i], runes[i-1]
		}
	}
	return string(runes)
}

// ValidateAndFormatMask formats the string s based on the mask, reporting any invalid runes that doesn't fit the provided mask.
//
// The mask argument consists of three fields separated by semicolons: the mask, a pad character and a Right-to-Left specifier. Only the first field is required.
//
// The pad character is the character added to the result string when a required character class cannot be matched in the source string (note that the unmatched character in the source string is not skipped). If no pad character is specified, a space is used.
//
// If RTL processing is requested ('1'), the mask and the source string are reversed before processing, and the output is reversed afterwards. This causes extra characters in the source string to be at the left of the masked text and if the mask is bigger than the source string, the extra required characters in the mask are also padded to left.
//
// An information table about the mask symbols can be found at: https://github.com/frones/strmask
func ValidateAndFormatMask(mask string, str string) (string, error) {
	padChar := ' '
	rtl := false

	args := strings.Split(mask, ";")
	for _, s := range strings.Split(mask, ";")[1:] {
		if r, _ := utf8.DecodeLastRuneInString(args[0]); r == '\\' {
			args[0] += ";" + s
			args = append(args[:1], args[2:]...)
		}
	}
	mask = args[0]
	if len(args) >= 2 {
		padChar, _ = utf8.DecodeRuneInString(args[1])
	}
	if len(args) >= 3 {
		rtl = args[2] == "1"
	}

	if rtl {
		mask = reverseBackslashs(reverse(mask))
		str = reverse(str)
	}

	printNext := false
	output := ""
	strOffset := 0
	lastErrOffset := -1
	err := ""
	charCase := 'N'
	for _, r := range mask {
		if printNext {
			output += string(r)
			printNext = false

			if r2, w := utf8.DecodeRuneInString(str[strOffset:]); w > 0 && r2 == r {
				strOffset += w
			}

			continue
		}

		switch r {
		case '\\':
			printNext = true
		case '>':
			charCase = 'U'
		case '<':
			charCase = 'L'
		case '=':
			charCase = 'N'
		case '9', '0':
			if r2, w := utf8.DecodeRuneInString(str[strOffset:]); w > 0 && unicode.IsDigit(r2) {
				output += string(r2)
				strOffset += w
			} else if r == '0' {
				if lastErrOffset < strOffset {
					if w > 0 {
						err += fmt.Sprintf("invalid character \"%s\" (expected a digit) at position %d\n", string(r2), strOffset)
					} else {
						err += fmt.Sprintf("expected a digit at position %d, but end of string found\n", strOffset)
					}
					lastErrOffset = strOffset
				}
				output += string(padChar)
			}
		case 'L', 'l':
			if r2, w := utf8.DecodeRuneInString(str[strOffset:]); w > 0 && unicode.In(r2, asciiLetters) {
				if charCase == 'U' {
					output += string(unicode.ToUpper(r2))
				} else if charCase == 'L' {
					output += string(unicode.ToLower(r2))
				} else {
					output += string(r2)
				}
				strOffset += w
			} else if r == 'L' {
				if lastErrOffset < strOffset {
					if w > 0 {
						err += fmt.Sprintf("invalid character \"%s\" (expected an ascii letter) at position %d\n", string(r2), strOffset)
					} else {
						err += fmt.Sprintf("expected an ascii letter at position %d, but end of string found\n", strOffset)
					}
					lastErrOffset = strOffset
				}
				output += string(padChar)
			}
		case 'A', 'a':
			if r2, w := utf8.DecodeRuneInString(str[strOffset:]); w > 0 && (unicode.In(r2, asciiLetters) || unicode.IsDigit(r2)) {
				if charCase == 'U' {
					output += string(unicode.ToUpper(r2))
				} else if charCase == 'L' {
					output += string(unicode.ToLower(r2))
				} else {
					output += string(r2)
				}
				strOffset += w
			} else if r == 'A' {
				if lastErrOffset < strOffset {
					if w > 0 {
						err += fmt.Sprintf("invalid character \"%s\" (expected an ascii letter or a digit) at position %d\n", string(r2), strOffset)
					} else {
						err += fmt.Sprintf("expected an ascii letter or a digit at position %d, but end of string found\n", strOffset)
					}
					lastErrOffset = strOffset
				}
				output += string(padChar)
			}
		case 'W', 'w':
			if r2, w := utf8.DecodeRuneInString(str[strOffset:]); w > 0 && unicode.IsLetter(r2) {
				if charCase == 'U' {
					output += string(unicode.ToUpper(r2))
				} else if charCase == 'L' {
					output += string(unicode.ToLower(r2))
				} else {
					output += string(r2)
				}
				strOffset += w
			} else if r == 'W' {
				if lastErrOffset < strOffset {
					if w > 0 {
						err += fmt.Sprintf("invalid character \"%s\" (expected an unicode letter) at position %d\n", string(r2), strOffset)
					} else {
						err += fmt.Sprintf("expected an unicode letter at position %d, but end of string found\n", strOffset)
					}
					lastErrOffset = strOffset
				}
				output += string(padChar)
			}
		case 'C', 'c':
			if r2, w := utf8.DecodeRuneInString(str[strOffset:]); w > 0 {
				if charCase == 'U' {
					output += string(unicode.ToUpper(r2))
				} else if charCase == 'L' {
					output += string(unicode.ToLower(r2))
				} else {
					output += string(r2)
				}
				strOffset += w
			}
		default:
			output += string(r)
			r2, w := utf8.DecodeRuneInString(str[strOffset:])
			if w > 0 && r2 == r {
				strOffset += w
			}
		}
	}

	if len(str) > strOffset {
		output += str[strOffset:]
	}

	if rtl {
		output = reverse(output)
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
