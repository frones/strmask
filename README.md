# strmask
This package provides a simple string formatter function for go based on masks.

The mask argument consists of three fields separated by semicolons: the mask, a pad character and a Right-to-Left specifier. Only the first field is required.

## Options Fields

### Pad character
The pad character is the character added to the result string when a required character class cannot be matched in the source string (notice that the unmatched character in the source string is not skipped). If no pad character is specified, a space is used.

### Right-to-Left specifier
If RTL processing is requested ('1'), the mask and the source string are reversed before processing, and the output is reversed afterwards. This causes extra characters in the source string to be at the left of the masked text and if the mask is bigger than the source string, the extra required characters in the mask are also padded to left.

# Usage
```go
fmt.Printf("%s", strmask.FormatMask("00.000.000/0000-00;0;1", "12520501000188")) // Pad with zeroes, process RTL
```
Outputs `12.520.501/0001-88`

```go
s, err := strmask.ValidateAndFormatMask("LLL-0000", "OL4508")
fmt.Printf("%s, %v\n", s, err)
```
Outputs `OL -4508, invalid character "4" (expected a letter) at position 2`

# Mask Symbols
Symbol | Meaning
--- | ---
`;` | Used as a field separator for the mask
`\` | Escapes next character
`>` | Maps all letters that follow the symbol to upper case until the end of the mask or until a resetting `=` appears
`<` | Maps all letters that follow the symbol to lower case until the end of the mask or until a resetting `=` appears
`=` | Resets the case modifiers (`>` and `<`)
`0` | Requires a decimal digit at this position
`9` | Permits a decimal digit at this position, but doesn't require it
`L` | Requires an ascii letter at this position
`l` | Permits an ascii letter at this position, but doesn't require it
`A` | Requires a decimal digit or ascii letter at this position
`a` | Permits a decimal digit or ascii letter at this position, but doesn't require it
`W` | Requires an unicode letter at this position
`w` | Permits an unicode letter at this position, but doesn't require it
`C` | Requires an arbitrary character at this position
`c` | Permits an arbitrary character at this position, but doesn't require it

# Acknowledgement
I took some inspiration from Delphi's `FormatMaskText` and from https://github.com/the-darc/string-mask, but didn't bother reading their codes and just wrote it from scratch. I don't know if the output is compatible with any of these (or any other mask formatting library out there). I guess not, but it's good enough for what I need.
