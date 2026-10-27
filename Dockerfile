# Stage 1: Build
FROM golang:1.27.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Stage 2: Runtime
FROM alpine:3.22

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server ./server

EXPOSE 3000

CMD ["./server"]