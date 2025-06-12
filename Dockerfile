FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/docker-notifier ./cmd/docker-notifier

FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata shadow docker-cli

WORKDIR /app

COPY --from=builder /app/bin/docker-notifier /app/docker-notifier

CMD ["/app/docker-notifier"]
