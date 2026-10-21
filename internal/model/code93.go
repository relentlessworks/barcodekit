package model

import (
	"fmt"
	"strings"
)

// Code 93 encoding table.
// Each character is encoded as 9 elements (3 bars, 3 spaces, each 1-4 modules wide).
var code93Table = map[byte][]int{
	'0': {1, 1, 1, 3, 1, 1, 2, 3, 1},
	'1': {1, 1, 1, 3, 1, 2, 1, 3, 1},
	'2': {1, 1, 1, 3, 2, 1, 1, 3, 1},
	'3': {1, 1, 1, 3, 2, 2, 1, 3, 1},
	'4': {1, 1, 1, 3, 1, 1, 1, 3, 2},
	'5': {1, 1, 1, 3, 1, 2, 1, 3, 2},
	'6': {1, 1, 1, 3, 2, 1, 1, 3, 2},
	'7': {1, 1, 1, 3, 2, 2, 1, 3, 2},
	'8': {1, 1, 2, 3, 1, 1, 1, 3, 1},
	'9': {1, 1, 2, 3, 1, 2, 1, 3, 1},
	'A': {1, 1, 2, 3, 2, 1, 1, 3, 1},
	'B': {1, 1, 2, 3, 2, 2, 1, 3, 1},
	'C': {1, 1, 2, 3, 1, 1, 1, 3, 2},
	'D': {1, 1, 2, 3, 1, 2, 1, 3, 2},
	'E': {1, 1, 2, 3, 2, 1, 1, 3, 2},
	'F': {1, 1, 2, 3, 2, 2, 1, 3, 2},
	'G': {1, 2, 1, 3, 1, 1, 1, 3, 1},
	'H': {1, 2, 1, 3, 1, 2, 1, 3, 1},
	'I': {1, 2, 1, 3, 2, 1, 1, 3, 1},
	'J': {1, 2, 1, 3, 2, 2, 1, 3, 1},
	'K': {1, 2, 1, 3, 1, 1, 1, 3, 2},
	'L': {1, 2, 1, 3, 1, 2, 1, 3, 2},
	'M': {1, 2, 1, 3, 2, 1, 1, 3, 2},
	'N': {1, 2, 1, 3, 2, 2, 1, 3, 2},
	'O': {1, 2, 2, 3, 1, 1, 1, 3, 1},
	'P': {1, 2, 2, 3, 1, 2, 1, 3, 1},
	'Q': {1, 2, 2, 3, 2, 1, 1, 3, 1},
	'R': {1, 2, 2, 3, 2, 2, 1, 3, 1},
	'S': {1, 2, 2, 3, 1, 1, 1, 3, 2},
	'T': {1, 2, 2, 3, 1, 2, 1, 3, 2},
	'U': {1, 2, 2, 3, 2, 1, 1, 3, 2},
	'V': {1, 2, 2, 3, 2, 2, 1, 3, 2},
	'W': {1, 3, 1, 3, 1, 1, 1, 3, 1},
	'X': {1, 3, 1, 3, 1, 2, 1, 3, 1},
	'Y': {1, 3, 1, 3, 2, 1, 1, 3, 1},
	'Z': {1, 3, 1, 3, 2, 2, 1, 3, 1},
	'-': {1, 3, 1, 3, 1, 1, 1, 3, 2},
	'.': {1, 3, 1, 3, 1, 2, 1, 3, 2},
	' ': {1, 3, 1, 3, 2, 1, 1, 3, 2},
	'$': {1, 3, 2, 3, 1, 1, 1, 3, 1},
	'/': {1, 3, 2, 3, 1, 2, 1, 3, 1},
	'+': {1, 3, 2, 3, 2, 1, 1, 3, 1},
	'%': {1, 3, 2, 3, 1, 1, 1, 3, 2},
}

// Code 93 character values for checksum calculation.
var code93Values = map[byte]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17, 'I': 18, 'J': 19,
	'K': 20, 'L': 21, 'M': 22, 'N': 23, 'O': 24, 'P': 25, 'Q': 26, 'R': 27, 'S': 28, 'T': 29,
	'U': 30, 'V': 31, 'W': 32, 'X': 33, 'Y': 34, 'Z': 35, '-': 36, '.': 37, ' ': 38, '$': 39,
	'/': 40, '+': 41, '%': 42,
}

var code93Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ-. $/+%"

// Start/stop pattern for Code 93
var code93StartStop = []int{1, 1, 1, 1, 4, 1, 1, 3, 1}

func code93ChecksumC(data string) byte {
	sum := 0
	weight := 1
	for i := len(data) - 1; i >= 0; i-- {
		v, ok := code93Values[data[i]]
		if !ok {
			return '?'
		}
		sum += v * weight
		weight++
		if weight > 20 {
			weight = 1
		}
	}
	return code93Chars[sum%47]
}

func code93ChecksumK(data string) byte {
	// data includes the C checksum
	sum := 0
	weight := 1
	for i := len(data) - 1; i >= 0; i-- {
		v, ok := code93Values[data[i]]
		if !ok {
			return '?'
		}
		sum += v * weight
		weight++
		if weight > 15 {
			weight = 1
		}
	}
	return code93Chars[sum%47]
}

func encodeCode93(data string) (*BarPattern, string, error) {
	if len(data) == 0 {
		return nil, "", fmt.Errorf("Code 93 data cannot be empty | hint: provide 1-80 characters from A-Z, 0-9, - . $ / + %% SPACE")
	}
	if len(data) > 80 {
		return nil, "", fmt.Errorf("Code 93 data too long (max 80 characters, got %d) | hint: shorten the data", len(data))
	}

	upper := strings.ToUpper(data)
	for i := 0; i < len(upper); i++ {
		c := upper[i]
		if _, ok := code93Values[c]; !ok {
			return nil, "", fmt.Errorf("Code 93 does not support character '%c' | hint: use only A-Z, 0-9, and - . $ / + %% SPACE", c)
		}
	}

	// Calculate checksums
	cCheck := code93ChecksumC(upper)
	kCheck := code93ChecksumK(upper + string(cCheck))

	// Build encoded string: data + C + K
	encoded := upper + string(cCheck) + string(kCheck)

	var widths []int
	// Start pattern
	widths = append(widths, code93StartStop...)

	for i := 0; i < len(encoded); i++ {
		c := encoded[i]
		if w, ok := code93Table[c]; ok {
			widths = append(widths, w...)
		}
	}

	// Stop pattern
	widths = append(widths, code93StartStop...)
	// Termination bar
	widths = append(widths, 1)

	return &BarPattern{Widths: widths}, encoded, nil
}
