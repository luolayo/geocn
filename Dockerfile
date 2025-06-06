# Stage 1: Build Go binary
FROM golang:1.24 as builder

WORKDIR /app

# Copy go mod and vendor
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o geocn main.go

# Stage 2: Create runtime image
FROM alpine:latest
# Add certificates for HTTPS
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update \
    && apk add --no-cache ca-certificates

# Set workdir
WORKDIR /app

# Copy binary
COPY --from=builder /app/geocn .

# Create data directory
RUN mkdir /app/data

# Expose port
EXPOSE 8080

# Run binary
ENTRYPOINT ["./geocn"]