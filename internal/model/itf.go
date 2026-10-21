package model

import "fmt"

// ITF (Interleaved 2 of 5) encoding table.
// Each digit is encoded as 5 elements: 2 wide (3 modules) and 3 narrow (1 module).
// Odd-position digits are bars, even-position digits are spaces (interleaved).
var itfWidths = [10][]int{
	{1, 1, 3, 3, 1}, // 0: NNWWN
	{3, 1, 1, 1, 3}, // 1: WNNNW
	{1, 3, 1, 1, 3}, // 2: NWNNW
	{3, 3, 1, 1, 1}, // 3: WWNNN
	{1, 1, 3, 1, 3}, // 4: NNWNW
	{3, 1, 3, 1, 1}, // 5: WNWNN
	{1, 3, 3, 1, 1}, // 6: NWWNN
	{1, 1, 1, 3, 3}, // 7: NNNWW
	{3, 1, 1, 3, 1}, // 8: WNNWN
	{1, 3, 1, 3, 1}, // 9: NWNWN
}

// ITF start guard: narrow bar, narrow space, narrow bar, narrow space
var itfStart = []int{1, 1, 1, 1}

// ITF end guard: wide bar, narrow space, narrow bar
var itfEnd = []int{3, 1, 1}

func encodeITF(data string) (*BarPattern, string, error) {
	if len(data) == 0 {
		return nil, "", fmt.Errorf("ITF data cannot be empty | hint: provide an even number of digits (0-9)")
	}
	if len(data)%2 != 0 {
		return nil, "", fmt.Errorf("ITF requires an even number of digits, got %d | hint: add a leading zero or pad to an even count", len(data))
	}
	if !isAllDigits(data) {
		return nil, "", fmt.Errorf("ITF data must contain only digits 0-9 | hint: remove any non-numeric characters")
	}

	var widths []int
	// Start guard
	widths = append(widths, itfStart...)

	// Interleave: odd-position digits are bars, even-position digits are spaces
	for i := 0; i < len(data); i += 2 {
		d1 := int(data[i] - '0')
		d2 := int(data[i+1] - '0')

		bars := itfWidths[d1]
		spaces := itfWidths[d2]

		// Interleave: bar, space, bar, space, bar, space, bar, space, bar, space
		for j := 0; j < 5; j++ {
			widths = append(widths, bars[j])   // bar
			widths = append(widths, spaces[j]) // space
		}
	}

	// End guard
	widths = append(widths, itfEnd...)

	return &BarPattern{Widths: widths}, data, nil
}
