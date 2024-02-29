# Start from the latest golang base image
FROM golang:alpine3.19 AS builder

# Install git
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy necessary files
COPY go.* ./
COPY . .

# Copy migrations directory
COPY shared/migrations /app/migrations

# Build the Go app with verbose output
RUN go build -o /app/main-out ./cmd/0xbase/

# Start a new stage from a minimal base image
FROM gcr.io/distroless/base-debian12 AS final

# Copy the binary and migrations from the builder stage
COPY --from=builder /app/main-out /app/
COPY --from=builder /app/migrations /app/migrations

# Copy the environment file
COPY --from=builder /app/app.env /

# Expose port
EXPOSE 5050

# Command to run the executable
CMD ["/app/main-out"]
