package model

import "fmt"

// MSI (Modified Plessey) encoding.
// Each digit is encoded as 4 bits. Each bit is a bar-space pair:
// bit 1 = wide bar (3) + narrow space (1)
// bit 0 = narrow bar (1) + narrow space (1)
var msiDigitTable = [10]string{
	"0000", // 0
	"0001", // 1
	"0010", // 2
	"0011", // 3
	"0100", // 4
	"0101", // 5
	"0110", // 6
	"0111", // 7
	"1000", // 8
	"1001", // 9
}

// MSI start guard: narrow bar, narrow space, narrow bar, narrow space
var msiStart = []int{1, 1, 1, 1}

// MSI end guard: narrow bar, narrow space, wide bar
var msiEnd = []int{1, 1, 3}

func encodeMSI(data string) (*BarPattern, string, error) {
	if len(data) == 0 {
		return nil, "", fmt.Errorf("MSI data cannot be empty | hint: provide digits 0-9")
	}
	if len(data) > 20 {
		return nil, "", fmt.Errorf("MSI data too long (max 20 digits, got %d) | hint: shorten the data", len(data))
	}
	if !isAllDigits(data) {
		return nil, "", fmt.Errorf("MSI data must contain only digits 0-9 | hint: remove any non-numeric characters")
	}

	var widths []int
	// Start guard
	widths = append(widths, msiStart...)

	// Encode each digit
	for i := 0; i < len(data); i++ {
		d := int(data[i] - '0')
		bits := msiDigitTable[d]
		for j := 0; j < 4; j++ {
			if bits[j] == '1' {
				widths = append(widths, 3, 1) // wide bar, narrow space
			} else {
				widths = append(widths, 1, 1) // narrow bar, narrow space
			}
		}
	}

	// End guard
	widths = append(widths, msiEnd...)

	return &BarPattern{Widths: widths}, data, nil
}
