package model

import "fmt"

// EAN-13 encoding tables.
// L-code (odd parity) and G-code (even parity) for digits 0-9.
var eanLTable = [10][]int{
	{3, 2, 1, 1}, // 0
	{2, 2, 2, 1}, // 1
	{2, 1, 2, 2}, // 2
	{1, 4, 1, 1}, // 3
	{1, 1, 3, 2}, // 4
	{1, 2, 3, 1}, // 5
	{1, 1, 1, 4}, // 6
	{1, 3, 1, 2}, // 7
	{1, 2, 1, 3}, // 8
	{3, 1, 1, 2}, // 9
}

var eanGTable = [10][]int{
	{1, 1, 2, 3}, // 0
	{1, 2, 2, 2}, // 1
	{2, 2, 1, 2}, // 2
	{1, 1, 4, 1}, // 3
	{2, 3, 1, 1}, // 4
	{1, 3, 2, 1}, // 5
	{4, 1, 1, 1}, // 6
	{2, 1, 3, 1}, // 7
	{3, 1, 2, 1}, // 8
	{2, 1, 1, 3}, // 9
}

// R-code (even parity, same as G mirrored) for the right half.
var eanRTable = [10][]int{
	{3, 2, 1, 1}, // 0 (mirrored: 1,1,2,3)
	{2, 2, 2, 1}, // 1
	{2, 1, 2, 2}, // 2
	{1, 4, 1, 1}, // 3
	{1, 1, 3, 2}, // 4
	{1, 2, 3, 1}, // 5
	{1, 1, 1, 4}, // 6
	{1, 3, 1, 2}, // 7
	{1, 2, 1, 3}, // 8
	{3, 1, 1, 2}, // 9
}

// First digit encoding pattern for EAN-13 (determines L/G pattern for digits 2-7).
var ean13Parity = [10]string{
	"LLLLLL", // 0
	"LLGLGG", // 1
	"LLGGLG", // 2
	"LLGGGL", // 3
	"LGLLGG", // 4
	"LGGLLG", // 5
	"LGGGLL", // 6
	"LGLGLG", // 7
	"LGLGGL", // 8
	"LGGLGL", // 9
}

// EAN-13: start guard + 6 L/G digits + center guard + 6 R digits + end guard
func encodeEAN13(data string) (*BarPattern, string, error) {
	if len(data) != 12 && len(data) != 13 {
		return nil, "", fmt.Errorf("EAN-13 requires 12 or 13 digits, got %d | hint: provide 12 digits and the checksum will be calculated, or 13 digits with a valid checksum", len(data))
	}
	if !isAllDigits(data) {
		return nil, "", fmt.Errorf("EAN-13 data must contain only digits 0-9 | hint: remove any non-numeric characters")
	}

	// Calculate or verify checksum
	var full string
	if len(data) == 12 {
		check := eanChecksum(data)
		full = data + fmt.Sprintf("%d", check)
	} else {
		// Verify checksum
		check := eanChecksum(data[:12])
		if int(data[12]-'0') != check {
			return nil, "", fmt.Errorf("EAN-13 checksum mismatch: expected %d, got %c | hint: provide 12 digits and let the service calculate the checksum", check, data[12])
		}
		full = data
	}

	firstDigit := int(full[0] - '0')
	parity := ean13Parity[firstDigit]

	var widths []int
	// Start guard: bar-space-bar (1,1,1)
	widths = append(widths, 1, 1, 1)

	// Left half: digits 2-7 (index 1-6)
	for i := 1; i <= 6; i++ {
		d := int(full[i] - '0')
		if parity[i-1] == 'G' {
			widths = append(widths, eanGTable[d]...)
		} else {
			widths = append(widths, eanLTable[d]...)
		}
	}

	// Center guard: space-bar-space-bar-space (1,1,1,1,1)
	widths = append(widths, 1, 1, 1, 1, 1)

	// Right half: digits 8-13 (index 7-12) using R-code
	for i := 7; i <= 12; i++ {
		d := int(full[i] - '0')
		widths = append(widths, eanRTable[d]...)
	}

	// End guard: bar-space-bar (1,1,1)
	widths = append(widths, 1, 1, 1)

	return &BarPattern{Widths: widths}, full, nil
}

// EAN-8: start guard + 4 L digits + center guard + 4 R digits + end guard
func encodeEAN8(data string) (*BarPattern, string, error) {
	if len(data) != 7 && len(data) != 8 {
		return nil, "", fmt.Errorf("EAN-8 requires 7 or 8 digits, got %d | hint: provide 7 digits and the checksum will be calculated", len(data))
	}
	if !isAllDigits(data) {
		return nil, "", fmt.Errorf("EAN-8 data must contain only digits 0-9 | hint: remove any non-numeric characters")
	}

	var full string
	if len(data) == 7 {
		check := eanChecksum(data)
		full = data + fmt.Sprintf("%d", check)
	} else {
		check := eanChecksum(data[:7])
		if int(data[7]-'0') != check {
			return nil, "", fmt.Errorf("EAN-8 checksum mismatch: expected %d, got %c | hint: provide 7 digits and let the service calculate the checksum", check, data[7])
		}
		full = data
	}

	var widths []int
	// Start guard
	widths = append(widths, 1, 1, 1)

	// Left half: digits 1-4 using L-code
	for i := 0; i < 4; i++ {
		d := int(full[i] - '0')
		widths = append(widths, eanLTable[d]...)
	}

	// Center guard
	widths = append(widths, 1, 1, 1, 1, 1)

	// Right half: digits 5-8 using R-code
	for i := 4; i < 8; i++ {
		d := int(full[i] - '0')
		widths = append(widths, eanRTable[d]...)
	}

	// End guard
	widths = append(widths, 1, 1, 1)

	return &BarPattern{Widths: widths}, full, nil
}

// UPC-A: start guard + 6 L digits + center guard + 6 R digits + end guard
// UPC-A is essentially EAN-13 with a leading 0.
func encodeUPCA(data string) (*BarPattern, string, error) {
	if len(data) != 11 && len(data) != 12 {
		return nil, "", fmt.Errorf("UPC-A requires 11 or 12 digits, got %d | hint: provide 11 digits and the checksum will be calculated", len(data))
	}
	if !isAllDigits(data) {
		return nil, "", fmt.Errorf("UPC-A data must contain only digits 0-9 | hint: remove any non-numeric characters")
	}

	var full string
	if len(data) == 11 {
		check := upcaChecksum(data)
		full = data + fmt.Sprintf("%d", check)
	} else {
		check := upcaChecksum(data[:11])
		if int(data[11]-'0') != check {
			return nil, "", fmt.Errorf("UPC-A checksum mismatch: expected %d, got %c | hint: provide 11 digits and let the service calculate the checksum", check, data[11])
		}
		full = data
	}

	var widths []int
	// Start guard
	widths = append(widths, 1, 1, 1)

	// Left half: digits 1-6 using L-code
	for i := 0; i < 6; i++ {
		d := int(full[i] - '0')
		widths = append(widths, eanLTable[d]...)
	}

	// Center guard
	widths = append(widths, 1, 1, 1, 1, 1)

	// Right half: digits 7-12 using R-code
	for i := 6; i < 12; i++ {
		d := int(full[i] - '0')
		widths = append(widths, eanRTable[d]...)
	}

	// End guard
	widths = append(widths, 1, 1, 1)

	return &BarPattern{Widths: widths}, full, nil
}
