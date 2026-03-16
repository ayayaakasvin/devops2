FROM golang:1.25.6 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest

WORKDIR /app

# Copy your app
COPY --from=builder /app/app .

ENV CONFIG_PATH=/etc/config/config.yaml

EXPOSE 80

CMD ["./app"]