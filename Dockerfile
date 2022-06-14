FROM golang:alpine as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o simple-http main.go
RUN adduser -u 10001 -D -H user

FROM scratch
WORKDIR /app
COPY --from=builder /app/simple-http /usr/bin/
COPY --from=0 /etc/passwd /etc/passwd
USER user
ENTRYPOINT ["simple-http"]
