# barcodekit

Agentic-first barcode generation service. Generate barcodes in 10 formats as SVG or PNG. Plain text API, agent-driven, single Go binary.

## Supported Barcode Types

| Type | Description | Data Pattern |
|------|-------------|--------------|
| code128 | Code 128 — high-density, all ASCII | any ASCII text (1-80 chars) |
| ean13 | EAN-13 / GTIN-13 | 12 or 13 digits |
| ean8 | EAN-8 | 7 or 8 digits |
| upca | UPC-A | 11 or 12 digits |
| code39 | Code 39 — alphanumeric | A-Z, 0-9, - . $ / + % SPACE |
| code93 | Code 93 — compact alphanumeric | A-Z, 0-9, - . $ / + % SPACE |
| itf | ITF (Interleaved 2 of 5) | even number of digits |
| codabar | Codabar | digits 0-9 and - $ : / . + |
| msi | MSI (Modified Plessey) | digits 0-9 |
| pharmacode | Pharmacode | integer 3-131070 |

## Quick Start

```bash
# Build
make build

# Run
./barcodekit -no-auth

# Generate a barcode
curl "http://localhost:8080/generate?type=code128&data=Hello&format=svg"
curl "http://localhost:8080/code39?data=HELLO&format=png"
curl -X POST http://localhost:8080/generate -d 'type=ean13&data=590123412345&format=svg'
```

## API

### POST /generate
Generate a barcode. Parameters: `type`, `data`, `format` (svg/png), `width`, `height`, `padding`.

### GET /{type}
Shortcut endpoint. e.g. `GET /code128?data=Hello&format=svg`

### GET /help
Self-documenting help text for agents.

### GET /types
List all supported barcode types.

### GET /formats
List all supported output formats.

### POST /mcp
MCP (Model Context Protocol) JSON-RPC 2.0 endpoint.

## Response Format

Plain text by default. JSON via `Accept: application/json` header or `?format=json` query param.

## Build

```bash
CGO_ENABLED=0 go build -o barcodekit ./cmd/barcodekit
```

## License

MIT
