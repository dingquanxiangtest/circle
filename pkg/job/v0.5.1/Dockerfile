FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM golang:1.16.6-alpine3.14 AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o jobfilter -mod=vendor -ldflags='-s -w'  -installsuffix cgo pkg/job/v0.5.1/main.go

FROM scratch
COPY --from=certs /etc/ssl/certs /etc/ssl/certs

WORKDIR /structorjob
COPY --from=builder ./build/jobfilter ./cmd/

EXPOSE 80

ENTRYPOINT ["./cmd/jobfilter","-config=/configs/config.yml"]
