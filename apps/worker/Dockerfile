FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy worker go mod files
COPY apps/worker/go.mod apps/worker/go.sum ./
RUN go mod download

# Copy worker source code
COPY apps/worker/ ./

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./cmd/worker

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/worker .

CMD ["./worker"]

