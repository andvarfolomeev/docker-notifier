FROM golang:1.23-alpine AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=$(echo "$TARGETPLATFORM" | cut -d'/' -f1) && \
    GOARCH=$(echo "$TARGETPLATFORM" | cut -d'/' -f2) && \
    echo "GOOS=$GOOS GOARCH=$GOARCH" && \
    go build -o /app/bin/docker-notifier ./cmd/docker-notifier

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/bin/docker-notifier /app/docker-notifier

CMD ["/app/docker-notifier"]
