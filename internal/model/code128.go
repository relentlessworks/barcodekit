package model

import (
	"fmt"
	"strconv"
)

// Code 128 encoding table — 107 entries (values 0-106).
// Each entry is an 11-module pattern: '1' = bar, '0' = space.
// Source: ISO/IEC 15417:2007 / Wikipedia Code 128 article.
var code128Patterns = [107]string{
	"11011001100", "11001101100", "11001100110", "10010011000",
	"10010001100", "10001001100", "10011001000", "10011000100",
	"10001100100", "11001001000", "11001000100", "11000100100",
	"10110011100", "10011011100", "10011001110", "10111001100",
	"10011101100", "10011100110", "11001110010", "11001011100",
	"11001001110", "11011100100", "11001110100", "11101101110",
	"11101001100", "11100101100", "11100100110", "11101100100",
	"11100110100", "11100110010", "11011011000", "11011000110",
	"11000110110", "10100011000", "10001011000", "10001000110",
	"10110001000", "10001101000", "10001100010", "11010001000",
	"11000101000", "11000100010", "10110111000", "10110001110",
	"10001101110", "10111011000", "10111000110", "10001110110",
	"11101110110", "11010001110", "11000101110", "11011101000",
	"11011100010", "11011101110", "11101011000", "11101000110",
	"11100010110", "11101101000", "11101100010", "11100011010",
	"11101111010", "11001000010", "11110001010", "10100110000",
	"10100001100", "10010110000", "10010000110", "10000101100",
	"10000100110", "10110010000", "10110000100", "10011010000",
	"10011000010", "10000110100", "10000110010", "11000010010",
	"11001010000", "11110111010", "11000010100", "10001111010",
	"10100111100", "10010111100", "10010011110", "10111100100",
	"10011110100", "10011110010", "11110100100", "11110010100",
	"11110010010", "11011011110", "11011110110", "11110110110",
	"10101111000", "10100011110", "10001011110", "10111101000",
	"10111100010", "11110101000", "11110100010", "10111011110",
	"10111101110", "11101011110", "11110101110", "11010000100",
	"11010010000", "11010011100", "11000111010",
}

// Code 128 start codes.
const (
	code128StartA = 103
	code128StartB = 104
	code128StartC = 105
	code128CodeA  = 101
	code128CodeB  = 100
	code128CodeC  = 99
)

// patternToWidths converts an 11-module pattern string to bar widths.
// The pattern uses '1' for a bar module and '0' for a space module.
// We convert to alternating bar/space widths by counting runs.
func patternToWidths(pattern string) []int {
	var widths []int
	i := 0
	for i < len(pattern) {
		bit := pattern[i]
		count := 0
		for i < len(pattern) && pattern[i] == bit {
			count++
			i++
		}
		widths = append(widths, count)
	}
	return widths
}

func encodeCode128(data string) (*BarPattern, string, error) {
	if len(data) == 0 {
		return nil, "", fmt.Errorf("Code 128 data cannot be empty | hint: provide 1-80 ASCII characters")
	}
	if len(data) > 80 {
		return nil, "", fmt.Errorf("Code 128 data too long (max 80 characters, got %d) | hint: shorten the data", len(data))
	}

	// Validate ASCII
	for i := 0; i < len(data); i++ {
		if data[i] > 127 {
			return nil, "", fmt.Errorf("Code 128 only supports ASCII characters | hint: remove non-ASCII characters")
		}
	}

	// Determine best encoding strategy.
	// If all digits and even count, use Code C.
	// Otherwise, use Code B (supports printable ASCII).
	var values []int
	var encoded string

	if len(data) >= 2 && isAllDigits(data) && len(data)%2 == 0 {
		// Use Code C (numeric pairs)
		values = append(values, code128StartC)
		encoded = ""
		for i := 0; i < len(data); i += 2 {
			pair, _ := strconv.Atoi(data[i : i+2])
			values = append(values, pair)
			encoded += data[i : i+2]
		}
	} else {
		// Use Code B (printable ASCII 32-126)
		values = append(values, code128StartB)
		encoded = ""
		for i := 0; i < len(data); i++ {
			c := data[i]
			if c < 32 || c > 126 {
				// Use Code A for control characters (0-31)
				if c < 32 {
					values = append(values, code128CodeA)
					values = append(values, int(c))
				} else {
					// DEL (127) - use Code A
					values = append(values, code128CodeA)
					values = append(values, int(c-64))
				}
			} else {
				// Code B: value = char - 32
				values = append(values, int(c-32))
			}
			encoded += string(c)
		}
	}

	// Calculate checksum
	checksum := values[0]
	for i := 1; i < len(values); i++ {
		checksum += values[i] * i
	}
	checksum = checksum % 103
	values = append(values, checksum)

	// Build bar pattern
	var widths []int
	for _, v := range values {
		widths = append(widths, patternToWidths(code128Patterns[v])...)
	}
	// Stop code (pattern 106 + termination bar)
	widths = append(widths, patternToWidths(code128Patterns[106])...)
	// Final bar (2 modules)
	widths = append(widths, 2)

	return &BarPattern{Widths: widths}, encoded, nil
}
