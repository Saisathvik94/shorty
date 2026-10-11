FROM golang:1.26-alpine AS builder

WORKDIR /api

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN addgroup -S -g 10001 shorty && \
    adduser -S -D -H -u 10001 -G shorty shorty

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /shorty-api \
    ./cmd/server


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /shorty-api /shorty-api

USER shorty

EXPOSE 3000

ENTRYPOINT ["/shorty-api"]