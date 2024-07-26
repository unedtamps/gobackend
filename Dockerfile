FROM golang:1.22.1-alpine3.19 as builder
WORKDIR /app
RUN apk add --no-cache curl
COPY . .
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN go build -o be main.go

FROM alpine:3.19 as final
WORKDIR /app
COPY --from=builder /app/migrate /bin/migrate
COPY --from=builder /app/be /bin/app
COPY tools/start.sh .
COPY internal/migration ./internal/migration
COPY public ./public
RUN chmod +x /app/start.sh
USER root
CMD ["./start.sh" ]
