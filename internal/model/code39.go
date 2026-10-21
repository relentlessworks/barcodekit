package model

import (
	"fmt"
	"strings"
)

// Code 39 encoding table.
// Each character is encoded as 9 elements (5 bars, 4 spaces).
// 3 of the 9 elements are wide (width 3), the rest are narrow (width 1).
var code39Table = map[byte][]int{
	'0': {1, 1, 1, 3, 3, 1, 3, 1, 1},
	'1': {3, 1, 1, 1, 1, 3, 1, 3, 1},
	'2': {1, 3, 1, 1, 1, 3, 1, 3, 1},
	'3': {3, 3, 1, 1, 1, 1, 1, 3, 1},
	'4': {1, 1, 3, 1, 1, 3, 1, 3, 1},
	'5': {3, 1, 3, 1, 1, 1, 1, 3, 1},
	'6': {1, 3, 3, 1, 1, 1, 1, 3, 1},
	'7': {1, 1, 1, 3, 1, 3, 1, 3, 1},
	'8': {3, 1, 1, 3, 1, 1, 1, 3, 1},
	'9': {1, 3, 1, 3, 1, 1, 1, 3, 1},
	'A': {3, 1, 1, 1, 1, 1, 3, 1, 3},
	'B': {1, 3, 1, 1, 1, 1, 3, 1, 3},
	'C': {3, 3, 1, 1, 1, 1, 1, 1, 3},
	'D': {1, 1, 3, 1, 1, 1, 3, 1, 3},
	'E': {3, 1, 3, 1, 1, 1, 1, 1, 3},
	'F': {1, 3, 3, 1, 1, 1, 1, 1, 3},
	'G': {1, 1, 1, 3, 1, 1, 3, 1, 3},
	'H': {3, 1, 1, 3, 1, 1, 1, 1, 3},
	'I': {1, 3, 1, 3, 1, 1, 1, 1, 3},
	'J': {1, 1, 3, 3, 1, 1, 1, 1, 3},
	'K': {3, 1, 1, 1, 1, 3, 3, 1, 1},
	'L': {1, 3, 1, 1, 1, 3, 3, 1, 1},
	'M': {3, 3, 1, 1, 1, 1, 3, 1, 1},
	'N': {1, 1, 3, 1, 1, 3, 3, 1, 1},
	'O': {3, 1, 3, 1, 1, 1, 3, 1, 1},
	'P': {1, 3, 3, 1, 1, 1, 3, 1, 1},
	'Q': {1, 1, 1, 3, 1, 3, 3, 1, 1},
	'R': {3, 1, 1, 3, 1, 3, 1, 1, 1},
	'S': {1, 3, 1, 3, 1, 3, 1, 1, 1},
	'T': {1, 1, 3, 3, 1, 3, 1, 1, 1},
	'U': {3, 1, 1, 1, 3, 1, 1, 3, 1},
	'V': {1, 3, 1, 1, 3, 1, 1, 3, 1},
	'W': {3, 3, 1, 1, 3, 1, 1, 1, 1},
	'X': {1, 1, 3, 1, 3, 1, 1, 3, 1},
	'Y': {3, 1, 3, 1, 3, 1, 1, 1, 1},
	'Z': {1, 3, 3, 1, 3, 1, 1, 1, 1},
	'-': {1, 1, 1, 3, 3, 1, 1, 3, 1},
	'.': {3, 1, 1, 3, 3, 1, 1, 1, 1},
	' ': {1, 3, 1, 3, 3, 1, 1, 1, 1},
	'$': {1, 3, 1, 1, 3, 3, 1, 1, 1},
	'/': {1, 3, 1, 1, 1, 3, 3, 1, 1},
	'+': {1, 1, 3, 1, 3, 3, 1, 1, 1},
	'%': {1, 1, 3, 1, 1, 3, 3, 1, 1},
	'*': {1, 1, 1, 3, 3, 3, 1, 1, 1}, // start/stop
}

func encodeCode39(data string) (*BarPattern, string, error) {
	if len(data) == 0 {
		return nil, "", fmt.Errorf("Code 39 data cannot be empty | hint: provide 1-80 characters from A-Z, 0-9, - . $ / + %% SPACE")
	}
	if len(data) > 80 {
		return nil, "", fmt.Errorf("Code 39 data too long (max 80 characters, got %d) | hint: shorten the data", len(data))
	}

	upper := strings.ToUpper(data)
	for i := 0; i < len(upper); i++ {
		c := upper[i]
		if c == '*' {
			return nil, "", fmt.Errorf("Code 39 data cannot contain '*' character | hint: '*' is reserved as the start/stop character and is added automatically")
		}
		if _, ok := code39Table[c]; !ok {
			return nil, "", fmt.Errorf("Code 39 does not support character '%c' | hint: use only A-Z, 0-9, and - . $ / + %% SPACE", c)
		}
	}

	// Build with start/stop *
	encoded := "*" + upper + "*"

	var widths []int
	for i := 0; i < len(encoded); i++ {
		c := encoded[i]
		widths = append(widths, code39Table[c]...)
		// Inter-character gap (narrow space) between characters
		if i < len(encoded)-1 {
			widths = append(widths, 1) // gap space
		}
	}

	return &BarPattern{Widths: widths}, encoded, nil
}
