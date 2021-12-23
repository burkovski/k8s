FROM golang:1.17-alpine AS builder
WORKDIR /app/
# Get dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download
# Build the executable
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main main.go


FROM alpine:latest
WORKDIR /app/
COPY --from=builder /app/main .
CMD ["./main"]
