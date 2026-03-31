FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=1.1.0" -o ipdock-client .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/ipdock-client /ipdock-client
ENTRYPOINT ["/ipdock-client"]
