FROM golang:1.14.2-alpine3.11

# Install FFmpeg dependencies
RUN apk add --no-cache ffmpeg

# Set working directory
WORKDIR /usr/src/app

# Copy application code into container
COPY . .

# Build pyx binaries
RUN go build -o pxy ./cmd/pxy/main.go

# Expose port
EXPOSE 8080

# Start App
CMD ["./pxy"]