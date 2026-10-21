package model

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// RenderSVG renders a bar pattern as an SVG string.
func RenderSVG(pattern *BarPattern, width, height, padding int) string {
	if width <= 0 {
		width = 200
	}
	if height <= 0 {
		height = 60
	}
	if padding < 0 {
		padding = 10
	}

	// Calculate total modules
	totalModules := 0
	for _, w := range pattern.Widths {
		totalModules += w
	}

	// Available width for bars
	availWidth := width - 2*padding
	if availWidth <= 0 {
		availWidth = width
	}
	moduleWidth := float64(availWidth) / float64(totalModules)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, width, height, width, height))
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="white"/>`, width, height))

	x := float64(padding)
	for i, w := range pattern.Widths {
		if i%2 == 0 {
			// Bar (black)
			barW := float64(w) * moduleWidth
			sb.WriteString(fmt.Sprintf(`<rect x="%.2f" y="%d" width="%.2f" height="%d" fill="black"/>`, x, padding, barW, height-2*padding))
		}
		x += float64(w) * moduleWidth
	}

	sb.WriteString("</svg>")
	return sb.String()
}

// RenderPNG renders a bar pattern as a PNG image encoded as base64.
// Since we can't use external image libraries, we generate a minimal PNG
// using a simple bitmap approach.
func RenderPNG(pattern *BarPattern, width, height, padding int) (string, error) {
	if width <= 0 {
		width = 200
	}
	if height <= 0 {
		height = 60
	}
	if padding < 0 {
		padding = 10
	}

	// Create a bitmap (1 bit per pixel, packed into bytes)
	// We'll create a simple grayscale PNG
	imgWidth := width
	imgHeight := height

	// Create pixel data (RGBA, 4 bytes per pixel)
	pixels := make([]byte, imgWidth*imgHeight*4)

	// Fill with white
	for i := range pixels {
		if i%4 == 3 {
			pixels[i] = 255 // alpha
		} else {
			pixels[i] = 255 // white
		}
	}

	// Calculate total modules
	totalModules := 0
	for _, w := range pattern.Widths {
		totalModules += w
	}

	availWidth := imgWidth - 2*padding
	if availWidth <= 0 {
		availWidth = imgWidth
	}
	moduleWidth := float64(availWidth) / float64(totalModules)

	// Draw bars
	x := float64(padding)
	for i, w := range pattern.Widths {
		if i%2 == 0 {
			// Bar (black)
			barStart := int(x)
			barEnd := int(x + float64(w)*moduleWidth)
			if barEnd > imgWidth {
				barEnd = imgWidth
			}
			for py := padding; py < imgHeight-padding; py++ {
				for px := barStart; px < barEnd; px++ {
					idx := (py*imgWidth + px) * 4
					if idx >= 0 && idx+3 < len(pixels) {
						pixels[idx] = 0   // R
						pixels[idx+1] = 0 // G
						pixels[idx+2] = 0 // B
					}
				}
			}
		}
		x += float64(w) * moduleWidth
	}

	// Encode as PNG
	pngData, err := encodePNG(pixels, imgWidth, imgHeight)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(pngData), nil
}

// encodePNG creates a minimal PNG file from RGBA pixel data.
func encodePNG(pixels []byte, width, height int) ([]byte, error) {
	// PNG signature
	var buf []byte
	buf = append(buf, 0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A)

	// IHDR chunk
	ihdr := make([]byte, 13)
	ihdr[0] = byte(width >> 24)
	ihdr[1] = byte(width >> 16)
	ihdr[2] = byte(width >> 8)
	ihdr[3] = byte(width)
	ihdr[4] = byte(height >> 24)
	ihdr[5] = byte(height >> 16)
	ihdr[6] = byte(height >> 8)
	ihdr[7] = byte(height)
	ihdr[8] = 8  // bit depth
	ihdr[9] = 6  // color type (RGBA)
	ihdr[10] = 0 // compression
	ihdr[11] = 0 // filter
	ihdr[12] = 0 // interlace
	buf = append(buf, makeChunk("IHDR", ihdr)...)

	// IDAT chunk - raw pixel data with filter bytes
	var raw []byte
	for y := 0; y < height; y++ {
		raw = append(raw, 0) // filter: none
		start := y * width * 4
		end := start + width*4
		if end > len(pixels) {
			end = len(pixels)
		}
		raw = append(raw, pixels[start:end]...)
	}

	compressed, err := zlibCompress(raw)
	if err != nil {
		return nil, err
	}
	buf = append(buf, makeChunk("IDAT", compressed)...)

	// IEND chunk
	buf = append(buf, makeChunk("IEND", nil)...)

	return buf, nil
}

// makeChunk creates a PNG chunk with CRC32.
func makeChunk(chunkType string, data []byte) []byte {
	var chunk []byte
	length := len(data)
	chunk = append(chunk, byte(length>>24), byte(length>>16), byte(length>>8), byte(length))
	chunk = append(chunk, []byte(chunkType)...)
	chunk = append(chunk, data...)
	crc := crc32(chunk[4:])
	chunk = append(chunk, byte(crc>>24), byte(crc>>16), byte(crc>>8), byte(crc))
	return chunk
}

// crc32 computes CRC32 checksum (PNG uses IEEE 802.3 polynomial).
func crc32(data []byte) uint32 {
	crc := uint32(0xFFFFFFFF)
	for _, b := range data {
		crc ^= uint32(b)
		for i := 0; i < 8; i++ {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ 0xEDB88320
			} else {
				crc >>= 1
			}
		}
	}
	return crc ^ 0xFFFFFFFF
}

// zlibCompress compresses data using zlib (deflate with header).
// We use a simple stored block approach (no actual compression).
func zlibCompress(data []byte) ([]byte, error) {
	var buf []byte
	// zlib header: CMF=0x78 (deflate, window size 32K), FLG=0x01 (no dict, level 0)
	buf = append(buf, 0x78, 0x01)

	// Deflate stored blocks
	offset := 0
	for offset < len(data) {
		blockLen := len(data) - offset
		if blockLen > 65535 {
			blockLen = 65535
		}
		isLast := byte(0)
		if offset+blockLen >= len(data) {
			isLast = 1
		}
		buf = append(buf, isLast) // BFINAL + BTYPE=00 (stored)
		// LEN and NLEN (little-endian)
		buf = append(buf, byte(blockLen), byte(blockLen>>8))
		nlen := ^blockLen
		buf = append(buf, byte(nlen), byte(nlen>>8))
		buf = append(buf, data[offset:offset+blockLen]...)
		offset += blockLen
	}

	// If data is empty, add an empty final block
	if len(data) == 0 {
		buf = append(buf, 1, 0, 0, 0xFF, 0xFF)
	}

	// Adler-32 checksum
	checksum := adler32(data)
	buf = append(buf, byte(checksum>>24), byte(checksum>>16), byte(checksum>>8), byte(checksum))

	return buf, nil
}

// adler32 computes the Adler-32 checksum.
func adler32(data []byte) uint32 {
	const mod = 65521
	a := uint32(1)
	b := uint32(0)
	for _, c := range data {
		a = (a + uint32(c)) % mod
		b = (b + a) % mod
	}
	return (b << 16) | a
}
