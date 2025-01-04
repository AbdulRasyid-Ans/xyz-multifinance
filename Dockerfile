# Stage 1: Build
FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Stage 2: Production Image
FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main .
COPY .env ./

EXPOSE 8080

CMD ["./main"]
