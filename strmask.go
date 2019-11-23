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

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func reverseBackslashs(s string) string {
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' && i > 0 {
			runes[i-1], runes[i] = runes[i], runes[i-1]
		}
	}
	return string(runes)
}

// ValidateAndFormatMask formats the string s based on the mask, reporting any invalid runes that doesn't fit the provided mask. Information about the mask symbols can be found at: https://github.com/frones/strmask
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

			if r2, w := utf8.DecodeRuneInString(str[strOffset:]); r2 == r {
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
			r2, w := utf8.DecodeRuneInString(str[strOffset:])
			if unicode.IsDigit(r2) {
				output += string(r2)
				strOffset += w
			} else if r == '0' {
				if lastErrOffset < strOffset {
					err += fmt.Sprintf("invalid character \"%s\" (expected a digit) at position %d\n", string(r2), strOffset)
					lastErrOffset = strOffset
				}
				output += string(padChar)
			}
		case 'L', 'l':
			r2, w := utf8.DecodeRuneInString(str[strOffset:])
			if unicode.In(r2, asciiLetters) {
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
					err += fmt.Sprintf("invalid character \"%s\" (expected an ascii letter) at position %d\n", string(r2), strOffset)
					lastErrOffset = strOffset
				}
				output += string(padChar)
			}
		case 'A', 'a':
			r2, w := utf8.DecodeRuneInString(str[strOffset:])
			if unicode.In(r2, asciiLetters) || unicode.IsDigit(r2) {
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
					err += fmt.Sprintf("invalid character \"%s\" (expected an ascii letter or a digit) at position %d\n", string(r2), strOffset)
					lastErrOffset = strOffset
				}
				output += string(padChar)
			}
		case 'W', 'w':
			r2, w := utf8.DecodeRuneInString(str[strOffset:])
			if unicode.IsLetter(r2) {
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
					err += fmt.Sprintf("invalid character \"%s\" (expected an unicode letter) at position %d\n", string(r2), strOffset)
					lastErrOffset = strOffset
				}
				output += string(padChar)
			}
		case 'C', 'c':
			r2, w := utf8.DecodeRuneInString(str[strOffset:])
			if charCase == 'U' {
				output += string(unicode.ToUpper(r2))
			} else if charCase == 'L' {
				output += string(unicode.ToLower(r2))
			} else {
				output += string(r2)
			}
			strOffset += w
		default:
			output += string(r)
			r2, w := utf8.DecodeRuneInString(str[strOffset:])
			if r2 == r {
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
