FROM golang:1.25.6 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o built-binary .

FROM alpine:latest

WORKDIR /app

# Copy your app
COPY --from=builder /app/built-binary .

EXPOSE 8080

CMD ["./built-binary"]