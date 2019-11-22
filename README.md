# strmask
This package provides a simple string formatter function for go based on masks.

# Usage
```go
fmt.Printf("%s", strmask.FormatMask("00.000.000/0000-00", "12520501000188"))
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
`\` | Escapes next character in output
`0` | Requires a decimal digit at this position
`9` | Permits a decimal digit at this position, but doesn't require it
`L` | Requires an unicode letter at this position
`l` | Permits an unicode letter at this position, but doesn't require it
`A` | Requires a decimal digit or unicode letter at this position
`a` | Permits a decimal digit or unicode letter at this position, but doesn't require it
`U` | Requires an unicode letter at this position and maps it to upper case
`u` | Permits an unicode letter at this position, but doesn't require it, and maps it to upper case
`W` | Requires an unicode letter at this position and maps it to lower case
`w` | Permits an unicode letter at this position, but doesn't require it, and maps it to lower case

# Acknowledgement
I took some inspiration from https://github.com/the-darc/string-mask and Delphi's `FormatMaskText`, but my requirements were so simple that I didn't bother reading their codes and just wrote it from scratch. I don't know if the output is compatible with any of these (or any other mask formatting library out there). I guess not, but it's good enough for what I need.
