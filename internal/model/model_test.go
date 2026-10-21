package model

import (
	"strings"
	"testing"
)

func TestEAN13(t *testing.T) {
	tests := []struct {
		data    string
		wantErr  bool
		wantFull string
	}{
		{"590123412345", false, "5901234123457"}, // 12 digits, checksum auto
		{"5901234123457", false, "5901234123457"}, // 13 digits with valid checksum
		{"5901234123458", true, ""},               // 13 digits with invalid checksum
		{"123", true, ""},                          // too short
		{"abc", true, ""},                          // non-numeric
	}

	for _, tt := range tests {
		pattern, full, err := encodeEAN13(tt.data)
		if tt.wantErr {
			if err == nil {
				t.Errorf("encodeEAN13(%q) expected error, got none (full=%s)", tt.data, full)
			}
			continue
		}
		if err != nil {
			t.Errorf("encodeEAN13(%q) unexpected error: %v", tt.data, err)
			continue
		}
		if full != tt.wantFull {
			t.Errorf("encodeEAN13(%q) full = %q, want %q", tt.data, full, tt.wantFull)
		}
		if pattern == nil || len(pattern.Widths) == 0 {
			t.Errorf("encodeEAN13(%q) got nil/empty pattern", tt.data)
		}
	}
}

func TestEAN8(t *testing.T) {
	tests := []struct {
		data    string
		wantErr  bool
		wantFull string
	}{
		{"7351353", false, "73513537"}, // 7 digits, checksum auto
		{"73513537", false, "73513537"}, // 8 digits with valid checksum
		{"73513538", true, ""},          // 8 digits with invalid checksum
		{"123", true, ""},                 // too short
	}

	for _, tt := range tests {
		pattern, full, err := encodeEAN8(tt.data)
		if tt.wantErr {
			if err == nil {
				t.Errorf("encodeEAN8(%q) expected error, got none", tt.data)
			}
			continue
		}
		if err != nil {
			t.Errorf("encodeEAN8(%q) unexpected error: %v", tt.data, err)
			continue
		}
		if full != tt.wantFull {
			t.Errorf("encodeEAN8(%q) full = %q, want %q", tt.data, full, tt.wantFull)
		}
		if pattern == nil || len(pattern.Widths) == 0 {
			t.Errorf("encodeEAN8(%q) got nil/empty pattern", tt.data)
		}
	}
}

func TestUPCA(t *testing.T) {
	tests := []struct {
		data    string
		wantErr  bool
		wantFull string
	}{
		{"03600029145", false, "036000291452"}, // 11 digits, checksum auto
		{"036000291452", false, "036000291452"}, // 12 digits with valid checksum
		{"036000291453", true, ""},               // 12 digits with invalid checksum
		{"123", true, ""},                          // too short
	}

	for _, tt := range tests {
		pattern, full, err := encodeUPCA(tt.data)
		if tt.wantErr {
			if err == nil {
				t.Errorf("encodeUPCA(%q) expected error, got none", tt.data)
			}
			continue
		}
		if err != nil {
			t.Errorf("encodeUPCA(%q) unexpected error: %v", tt.data, err)
			continue
		}
		if full != tt.wantFull {
			t.Errorf("encodeUPCA(%q) full = %q, want %q", tt.data, full, tt.wantFull)
		}
		if pattern == nil || len(pattern.Widths) == 0 {
			t.Errorf("encodeUPCA(%q) got nil/empty pattern", tt.data)
		}
	}
}

func TestCode39(t *testing.T) {
	tests := []struct {
		data    string
		wantErr bool
	}{
		{"HELLO", false},
		{"ABC123", false},
		{"", true},
		{"hello", false}, // lowercase should be uppercased
		{"HELLO*", true}, // * is reserved
		{"hello!", true}, // ! not supported
	}

	for _, tt := range tests {
		pattern, _, err := encodeCode39(tt.data)
		if tt.wantErr {
			if err == nil {
				t.Errorf("encodeCode39(%q) expected error, got none", tt.data)
			}
			continue
		}
		if err != nil {
			t.Errorf("encodeCode39(%q) unexpected error: %v", tt.data, err)
			continue
		}
		if pattern == nil || len(pattern.Widths) == 0 {
			t.Errorf("encodeCode39(%q) got nil/empty pattern", tt.data)
		}
	}
}

func TestCode93(t *testing.T) {
	tests := []struct {
		data    string
		wantErr bool
	}{
		{"HELLO", false},
		{"ABC123", false},
		{"", true},
		{"hello", false}, // lowercase should be uppercased
		{"HELLO!", true}, // ! not supported
	}

	for _, tt := range tests {
		pattern, _, err := encodeCode93(tt.data)
		if tt.wantErr {
			if err == nil {
				t.Errorf("encodeCode93(%q) expected error, got none", tt.data)
			}
			continue
		}
		if err != nil {
			t.Errorf("encodeCode93(%q) unexpected error: %v", tt.data, err)
			continue
		}
		if pattern == nil || len(pattern.Widths) == 0 {
			t.Errorf("encodeCode93(%q) got nil/empty pattern", tt.data)
		}
	}
}

func TestITF(t *testing.T) {
	tests := []struct {
		data    string
		wantErr bool
	}{
		{"1234567890", false}, // even digits
		{"12345", true},        // odd digits
		{"", true},
		{"1234abc", true}, // non-numeric
	}

	for _, tt := range tests {
		pattern, _, err := encodeITF(tt.data)
		if tt.wantErr {
			if err == nil {
				t.Errorf("encodeITF(%q) expected error, got none", tt.data)
			}
			continue
		}
		if err != nil {
			t.Errorf("encodeITF(%q) unexpected error: %v", tt.data, err)
			continue
		}
		if pattern == nil || len(pattern.Widths) == 0 {
			t.Errorf("encodeITF(%q) got nil/empty pattern", tt.data)
		}
	}
}

func TestCodabar(t *testing.T) {
	tests := []struct {
		data    string
		wantErr bool
	}{
		{"123456", false},
		{"A123456B", true}, // start/stop chars not allowed in data
		{"", true},
		{"12345X", true}, // X not supported
	}

	for _, tt := range tests {
		pattern, _, err := encodeCodabar(tt.data)
		if tt.wantErr {
			if err == nil {
				t.Errorf("encodeCodabar(%q) expected error, got none", tt.data)
			}
			continue
		}
		if err != nil {
			t.Errorf("encodeCodabar(%q) unexpected error: %v", tt.data, err)
			continue
		}
		if pattern == nil || len(pattern.Widths) == 0 {
			t.Errorf("encodeCodabar(%q) got nil/empty pattern", tt.data)
		}
	}
}

func TestMSI(t *testing.T) {
	tests := []struct {
		data    string
		wantErr bool
	}{
		{"123456", false},
		{"", true},
		{"1234abc", true},
	}

	for _, tt := range tests {
		pattern, _, err := encodeMSI(tt.data)
		if tt.wantErr {
			if err == nil {
				t.Errorf("encodeMSI(%q) expected error, got none", tt.data)
			}
			continue
		}
		if err != nil {
			t.Errorf("encodeMSI(%q) unexpected error: %v", tt.data, err)
			continue
		}
		if pattern == nil || len(pattern.Widths) == 0 {
			t.Errorf("encodeMSI(%q) got nil/empty pattern", tt.data)
		}
	}
}

func TestPharmacode(t *testing.T) {
	tests := []struct {
		data    string
		wantErr bool
	}{
		{"3", false},
		{"131070", false},
		{"2", true},   // too small
		{"131071", true}, // too large
		{"", true},
		{"abc", true},
	}

	for _, tt := range tests {
		pattern, _, err := encodePharmacode(tt.data)
		if tt.wantErr {
			if err == nil {
				t.Errorf("encodePharmacode(%q) expected error, got none", tt.data)
			}
			continue
		}
		if err != nil {
			t.Errorf("encodePharmacode(%q) unexpected error: %v", tt.data, err)
			continue
		}
		if pattern == nil || len(pattern.Widths) == 0 {
			t.Errorf("encodePharmacode(%q) got nil/empty pattern", tt.data)
		}
	}
}

func TestCode128(t *testing.T) {
	tests := []struct {
		data    string
		wantErr bool
	}{
		{"Hello", false},
		{"12345678", false}, // even digits -> Code C
		{"", true},
		{"ABC", false},
	}

	for _, tt := range tests {
		pattern, _, err := encodeCode128(tt.data)
		if tt.wantErr {
			if err == nil {
				t.Errorf("encodeCode128(%q) expected error, got none", tt.data)
			}
			continue
		}
		if err != nil {
			t.Errorf("encodeCode128(%q) unexpected error: %v", tt.data, err)
			continue
		}
		if pattern == nil || len(pattern.Widths) == 0 {
			t.Errorf("encodeCode128(%q) got nil/empty pattern", tt.data)
		}
	}
}

func TestRenderSVG(t *testing.T) {
	pattern, _, err := encodeCode39("HELLO")
	if err != nil {
		t.Fatalf("encodeCode39 failed: %v", err)
	}

	svg := RenderSVG(pattern, 200, 60, 10)
	if !strings.Contains(svg, "<svg") {
		t.Error("RenderSVG output doesn't contain <svg tag")
	}
	if !strings.Contains(svg, "</svg>") {
		t.Error("RenderSVG output doesn't contain closing </svg> tag")
	}
}

func TestRenderPNG(t *testing.T) {
	pattern, _, err := encodeCode39("HELLO")
	if err != nil {
		t.Fatalf("encodeCode39 failed: %v", err)
	}

	b64, err := RenderPNG(pattern, 200, 60, 10)
	if err != nil {
		t.Fatalf("RenderPNG failed: %v", err)
	}
	if len(b64) == 0 {
		t.Error("RenderPNG returned empty output")
	}
}

func TestIsValidType(t *testing.T) {
	if !IsValidType(TypeCode128) {
		t.Error("TypeCode128 should be valid")
	}
	if IsValidType("invalid") {
		t.Error("invalid type should not be valid")
	}
}

func TestIsValidFormat(t *testing.T) {
	if !IsValidFormat(FormatSVG) {
		t.Error("FormatSVG should be valid")
	}
	if !IsValidFormat(FormatPNG) {
		t.Error("FormatPNG should be valid")
	}
	if IsValidFormat("invalid") {
		t.Error("invalid format should not be valid")
	}
}

func TestEANChecksum(t *testing.T) {
	// EAN-13: 590123412345 -> checksum 7
	check := eanChecksum("590123412345")
	if check != 7 {
		t.Errorf("eanChecksum(590123412345) = %d, want 7", check)
	}
}

func TestUPCAChecksum(t *testing.T) {
	// UPC-A: 03600029145 -> checksum 2
	check := upcaChecksum("03600029145")
	if check != 2 {
		t.Errorf("upcaChecksum(03600029145) = %d, want 2", check)
	}
}
