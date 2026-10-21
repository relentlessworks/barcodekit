package model

import "fmt"

// Codabar encoding table.
// Each character is 7 elements: 4 bars and 3 spaces.
// 2 or 3 of the 7 elements are wide (width 3), rest are narrow (width 1).
var codabarTable = map[byte][]int{
	'0': {3, 1, 1, 1, 1, 3, 1},
	'1': {1, 1, 3, 1, 1, 3, 1},
	'2': {1, 1, 1, 3, 1, 3, 1},
	'3': {3, 1, 3, 3, 1, 1, 1},
	'4': {1, 3, 1, 1, 1, 3, 1},
	'5': {3, 3, 1, 1, 1, 1, 1},
	'6': {1, 3, 3, 1, 1, 1, 1},
	'7': {1, 1, 1, 3, 3, 1, 1},
	'8': {3, 1, 1, 3, 3, 1, 1},
	'9': {3, 1, 3, 1, 3, 1, 1},
	'-': {1, 3, 1, 3, 1, 3, 1},
	'$': {1, 3, 3, 3, 1, 1, 1},
	':': {3, 3, 1, 3, 3, 3, 1},
	'/': {1, 1, 3, 3, 3, 1, 1},
	'.': {3, 1, 3, 3, 3, 1, 1},
	'+': {1, 1, 1, 3, 1, 3, 3},
	// Start/stop characters
	'A': {1, 1, 3, 1, 3, 1, 3},
	'B': {1, 1, 3, 1, 3, 1, 3}, // same as A
	'C': {1, 1, 1, 3, 3, 1, 3},
	'D': {1, 1, 3, 3, 1, 1, 3},
}

func encodeCodabar(data string) (*BarPattern, string, error) {
	if len(data) == 0 {
		return nil, "", fmt.Errorf("Codabar data cannot be empty | hint: provide digits 0-9 and - $ : / . +")
	}
	if len(data) > 80 {
		return nil, "", fmt.Errorf("Codabar data too long (max 80 characters, got %d) | hint: shorten the data", len(data))
	}

	// Validate characters (exclude start/stop A B C D from data)
	for i := 0; i < len(data); i++ {
		c := data[i]
		if c >= 'a' && c <= 'z' {
			c = c - 32 // uppercase
		}
		if c == 'A' || c == 'B' || c == 'C' || c == 'D' {
			return nil, "", fmt.Errorf("Codabar data cannot contain start/stop characters A B C D | hint: these are added automatically; use only digits 0-9 and - $ : / . +")
		}
		if _, ok := codabarTable[c]; !ok {
			return nil, "", fmt.Errorf("Codabar does not support character '%c' | hint: use only digits 0-9 and - $ : / . +", data[i])
		}
	}

	// Add start/stop characters
	encoded := "A" + data + "A"

	var widths []int
	for i := 0; i < len(encoded); i++ {
		c := encoded[i]
		widths = append(widths, codabarTable[c]...)
		// Inter-character gap (narrow space)
		if i < len(encoded)-1 {
			widths = append(widths, 1)
		}
	}

	return &BarPattern{Widths: widths}, encoded, nil
}
