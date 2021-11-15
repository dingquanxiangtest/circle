FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM golang:1.16.6-alpine3.14 AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o jobmenu -mod=vendor -ldflags='-s -w'  -installsuffix cgo pkg/job/v0.6_menu/main.go

FROM scratch
COPY --from=certs /etc/ssl/certs /etc/ssl/certs

WORKDIR /jobmenu
COPY --from=builder ./build/jobmenu ./cmd/

EXPOSE 80

ENTRYPOINT ["./cmd/jobmenu","-config=/configs/config.yml"]
