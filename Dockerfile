FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application with minimal image size optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -a -installsuffix cgo \
    -o /netkit .

FROM scratch

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /netkit /netkit
COPY --from=builder /app/templates /templates
COPY --from=builder /app/static /static

EXPOSE 8080

ENTRYPOINT ["/netkit"]