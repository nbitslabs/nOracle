# --- Build Stage ---
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN go build -o oracle cmd/main.go

# --- Final Stage: Minimal Container ---
FROM alpine:3.21.3

# Copy only the compiled binary from builder
COPY --from=builder /app/oracle /oracle

# Run the binary
ENTRYPOINT ["/oracle"]
