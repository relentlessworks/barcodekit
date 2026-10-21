FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o barcodekit ./cmd/barcodekit

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/barcodekit .
EXPOSE 8080
CMD ["./barcodekit", "-no-auth"]
