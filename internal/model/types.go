package model

// BarcodeType represents a supported barcode format.
type BarcodeType string

const (
	TypeCode128    BarcodeType = "code128"
	TypeEAN13      BarcodeType = "ean13"
	TypeEAN8       BarcodeType = "ean8"
	TypeUPCA       BarcodeType = "upca"
	TypeCode39     BarcodeType = "code39"
	TypeCode93     BarcodeType = "code93"
	TypeITF        BarcodeType = "itf"
	TypeCodabar    BarcodeType = "codabar"
	TypeMSI        BarcodeType = "msi"
	TypePharmacode BarcodeType = "pharmacode"
)

// OutputFormat represents the output format for a barcode.
type OutputFormat string

const (
	FormatSVG OutputFormat = "svg"
	FormatPNG OutputFormat = "png"
)

// BarcodeRequest holds parameters for generating a barcode.
type BarcodeRequest struct {
	Type    BarcodeType  `json:"type"`
	Data    string       `json:"data"`
	Format  OutputFormat `json:"format"`
	Width   int          `json:"width"`
	Height  int          `json:"height"`
	Padding int          `json:"padding"`
}

// BarcodeResult holds the generated barcode output.
type BarcodeResult struct {
	Type   BarcodeType  `json:"type"`
	Data   string       `json:"data"`
	Format OutputFormat `json:"format"`
	Output string       `json:"output"`
	Width  int          `json:"width"`
	Height int          `json:"height"`
}

// SupportedTypes returns all supported barcode types.
func SupportedTypes() []BarcodeType {
	return []BarcodeType{
		TypeCode128, TypeEAN13, TypeEAN8, TypeUPCA,
		TypeCode39, TypeCode93, TypeITF, TypeCodabar,
		TypeMSI, TypePharmacode,
	}
}

// SupportedFormats returns all supported output formats.
func SupportedFormats() []OutputFormat {
	return []OutputFormat{FormatSVG, FormatPNG}
}

// IsValidType checks if a barcode type is supported.
func IsValidType(t BarcodeType) bool {
	for _, st := range SupportedTypes() {
		if st == t {
			return true
		}
	}
	return false
}

// IsValidFormat checks if an output format is supported.
func IsValidFormat(f OutputFormat) bool {
	for _, sf := range SupportedFormats() {
		if sf == f {
			return true
		}
	}
	return false
}

// TypeDescription returns a human-readable description of a barcode type.
func TypeDescription(t BarcodeType) string {
	switch t {
	case TypeCode128:
		return "Code 128 — high-density linear barcode supporting all 128 ASCII characters"
	case TypeEAN13:
		return "EAN-13 — 13-digit international article number (GTIN-13)"
	case TypeEAN8:
		return "EAN-8 — 8-digit short form of EAN for small packages"
	case TypeUPCA:
		return "UPC-A — 12-digit Universal Product Code"
	case TypeCode39:
		return "Code 39 — alphanumeric barcode with 43 characters (A-Z, 0-9, and special)"
	case TypeCode93:
		return "Code 93 — compact alphanumeric barcode, superset of Code 39"
	case TypeITF:
		return "ITF (Interleaved 2 of 5) — numeric-only barcode, pairs of digits"
	case TypeCodabar:
		return "Codabar — numeric barcode with 6 special characters, used in libraries"
	case TypeMSI:
		return "MSI (Modified Plessey) — numeric barcode used for inventory"
	case TypePharmacode:
		return "Pharmacode — simple barcode used in pharmaceutical packaging"
	default:
		return "unknown barcode type"
	}
}

// TypeDataPattern returns a description of what data the barcode type accepts.
func TypeDataPattern(t BarcodeType) string {
	switch t {
	case TypeCode128:
		return "any ASCII text (1-80 characters)"
	case TypeEAN13:
		return "12 or 13 digits (checksum auto-calculated if 12 provided)"
	case TypeEAN8:
		return "7 or 8 digits (checksum auto-calculated if 7 provided)"
	case TypeUPCA:
		return "11 or 12 digits (checksum auto-calculated if 11 provided)"
	case TypeCode39:
		return "uppercase A-Z, 0-9, and - . $ / + % SPACE (start/stop * added automatically)"
	case TypeCode93:
		return "uppercase A-Z, 0-9, and - . $ / + % SPACE"
	case TypeITF:
		return "even number of digits (0-9)"
	case TypeCodabar:
		return "digits 0-9 and - $ : / . + (start/stop A B C D added automatically)"
	case TypeMSI:
		return "digits 0-9"
	case TypePharmacode:
		return "integer 3-131070"
	default:
		return "unknown"
	}
}
