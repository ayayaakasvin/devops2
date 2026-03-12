FROM golang:1.25.6 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o built-binary .

FROM alpine:latest

WORKDIR /app

# Copy your app
COPY --from=builder /app/built-binary .

ENV CONFIG_PATH=/etc/config/config.yaml

EXPOSE 80

CMD ["./built-binary"]