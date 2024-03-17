FROM alpine:3.19 as root-certs
RUN apk add --no-cache ca-certificates
RUN addgroup -g 1000 app
RUN adduser app -u 1000 -D -G app /home/app

FROM golang:1.22.1-alpine3.19 as builder
WORKDIR /apps
COPY --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY . .
RUN go build -o /todo main.go

FROM scratch as final
COPY --from=root-certs /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1000:1000 --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=1000:1000 --from=builder /todo /todo
COPY .env .
USER app
CMD ["/todo" ]
