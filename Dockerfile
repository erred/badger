FROM golang:alpine AS build

WORKDIR /app
COPY . .

RUN apk add --update --no-cache ca-certificates
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o=badger

FROM scratch

WORKDIR /app
COPY --from=build /app/badger .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080/tcp
ENTRYPOINT ["/app/badger"]
