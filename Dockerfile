# FROM golang:1.15-alpine AS build
# WORKDIR /src
# ENV CGO_ENABLED=0

# # Download dependencies
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy source and build
# COPY . .
# RUN GOOS=linux GOARCH=amd64 go build -o /out/backup-validator .

FROM alpine:3.21.8@sha256:ce64758a109eb420d874a118f87920e625e12d3634e03b4a5573fd9f6e5d3507 AS bin
EXPOSE 9178
ENTRYPOINT [ "/backup-validator" ]
WORKDIR /workdir

# Install packages
RUN apk add --no-cache ca-certificates docker restic=0.11.0-r0 && update-ca-certificates

USER 1001
COPY backup-validator /backup-validator
