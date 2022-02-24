FROM golang:alpine3.15 as build
COPY . .
RUN unset GOPATH \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' .

FROM alpine:3.15
COPY --from=build go/npm-spdx /app/npm-spdx
COPY data /app/data
RUN apk add --no-cache ca-certificates
WORKDIR /app
ENTRYPOINT ["./npm-spdx"]
