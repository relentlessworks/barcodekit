package model

import (
	"fmt"
	"strconv"
)

// Pharmacode (Laetus) encoding.
// The number is converted to binary, and each bit is encoded as:
// bit 1 = wide bar (3 modules) + narrow space (1 module)
// bit 0 = narrow bar (1 module) + wide space (3 modules)
// The barcode starts with a narrow space and ends with a narrow bar.
func encodePharmacode(data string) (*BarPattern, string, error) {
	if len(data) == 0 {
		return nil, "", fmt.Errorf("Pharmacode data cannot be empty | hint: provide an integer 3-131070")
	}

	num, err := strconv.Atoi(data)
	if err != nil {
		return nil, "", fmt.Errorf("Pharmacode data must be an integer | hint: provide a number between 3 and 131070")
	}
	if num < 3 || num > 131070 {
		return nil, "", fmt.Errorf("Pharmacode value must be between 3 and 131070, got %d | hint: use a value in the valid range", num)
	}

	var widths []int
	// Start with a narrow space
	widths = append(widths, 1)

	// Convert to binary and encode each bit
	// We process from LSB to MSB
	for num > 0 {
		bit := num & 1
		num >>= 1
		if bit == 1 {
			widths = append(widths, 3, 1) // wide bar, narrow space
		} else {
			widths = append(widths, 1, 3) // narrow bar, wide space
		}
	}

	// The last element should be a bar (narrow), not a space.
	// Remove the trailing space and add a narrow bar.
	if len(widths) > 0 {
		widths = widths[:len(widths)-1] // remove trailing space
	}
	widths = append(widths, 1) // narrow bar terminator

	return &BarPattern{Widths: widths}, data, nil
}
