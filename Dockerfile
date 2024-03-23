FROM golang:1.22.1-alpine3.19 as builder
WORKDIR /apps
COPY . .
COPY internal/migration /migration
RUN go build -o /todo main.go

FROM alpine:3.19 as final
WORKDIR /app
RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN mv ./migrate /bin
COPY --from=builder /todo /bin
COPY tools/start.sh .
COPY internal/migration ./migration
COPY public ./public
RUN chmod +x  /app/start.sh
RUN rm README.md
USER root
CMD ["./start.sh" ]
