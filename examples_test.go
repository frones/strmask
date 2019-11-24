package strmask

import (
	"fmt"
)

func ExampleFormatMask() {
	fmt.Printf("%s\n", strmask.FormatMask("00.000.000/0000-00;0;1", "12520501000188")) // Pad with zeroes, process Right-to-Left
	// Output: 12.520.501/0001-88
}

func ExampleValidateAndFormatMask() {
	s, err := strmask.ValidateAndFormatMask("LLL-0000", "OL4508") // Pad with spaces, process Left-to-Right
	fmt.Printf("%s, %v\n", s, err)
	// Output: OL -4508, invalid character "4" (expected an ascii letter) at position 2
}

func ExampleValidateAndFormatMask_rtl() {
	s, err := strmask.ValidateAndFormatMask("00.000.000;0;1", "450554") // Pad with zeroes, process Right-to-Left
	fmt.Printf("%s, %v\n", s, err)
	// Output: 00.450.554, expected a digit at position 6, but end of string found
}
