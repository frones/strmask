package strmask

import (
	"fmt"
)

func ExampleFormatMask() {
	fmt.Printf("%s\n", strmask.FormatMask("00.000.000/0000-00;0;1", "12520501000188")) // Pad with zeroes, process RTL
	// Output: 12.520.501/0001-88
}

func ExampleValidateAndFormatMask() {
	s, err := strmask.ValidateAndFormatMask("LLL-0000", "OL4508")
	fmt.Printf("%s, %v\n", s, err)
	// Output: OL -4508, invalid character "4" (expected a letter) at position 2
}
