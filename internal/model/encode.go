package model

import (
	"fmt"
	"strings"
)

// BarPattern represents a barcode as a sequence of bar widths.
// Each element is the width of a bar or space, alternating (bar, space, bar, space, ...).
type BarPattern struct {
	Widths []int // alternating bar/space widths in modules
}

// Encode generates the bar pattern for the given barcode type and data.
// Returns the bar widths pattern and the encoded data (with checksums etc).
func Encode(req BarcodeRequest) (*BarPattern, string, error) {
	switch req.Type {
	case TypeCode128:
		return encodeCode128(req.Data)
	case TypeEAN13:
		return encodeEAN13(req.Data)
	case TypeEAN8:
		return encodeEAN8(req.Data)
	case TypeUPCA:
		return encodeUPCA(req.Data)
	case TypeCode39:
		return encodeCode39(req.Data)
	case TypeCode93:
		return encodeCode93(req.Data)
	case TypeITF:
		return encodeITF(req.Data)
	case TypeCodabar:
		return encodeCodabar(req.Data)
	case TypeMSI:
		return encodeMSI(req.Data)
	case TypePharmacode:
		return encodePharmacode(req.Data)
	default:
		return nil, "", fmt.Errorf("unsupported barcode type: %s", req.Type)
	}
}

// isAllDigits checks if a string contains only digits.
func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// eanChecksum calculates the EAN/UPC checksum digit.
// The checksum is weighted from the right: the rightmost digit gets weight 3,
// the next gets weight 1, alternating. This works for both EAN-13 (12 digits)
// and EAN-8 (7 digits).
func eanChecksum(digits string) int {
	sum := 0
	n := len(digits)
	for i, c := range digits {
		d := int(c - '0')
		// Position from right (0-indexed)
		fromRight := n - 1 - i
		if fromRight%2 == 0 {
			sum += d * 3
		} else {
			sum += d
		}
	}
	return (10 - sum%10) % 10
}

// upcaChecksum calculates the UPC-A checksum digit.
// Same algorithm as EAN: weighted from the right.
func upcaChecksum(digits string) int {
	sum := 0
	n := len(digits)
	for i, c := range digits {
		d := int(c - '0')
		fromRight := n - 1 - i
		if fromRight%2 == 0 {
			sum += d * 3
		} else {
			sum += d
		}
	}
	return (10 - sum%10) % 10
}

// widthsToBinary converts bar widths to a binary string (1=bar, 0=space).
func widthsToBinary(widths []int) string {
	var sb strings.Builder
	for i, w := range widths {
		bit := byte('1')
		if i%2 == 1 {
			bit = '0'
		}
		for j := 0; j < w; j++ {
			sb.WriteByte(bit)
		}
	}
	return sb.String()
}
