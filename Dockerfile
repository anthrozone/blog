FROM golang:alpine as builder
RUN adduser -D -g '' appuser
WORKDIR /opt/build
RUN apk update && apk add --no-cache git
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o anthro-api

FROM alpine
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /opt/build/anthro-api /anthro-api
USER appuser
ENTRYPOINT ["/anthro-api"]