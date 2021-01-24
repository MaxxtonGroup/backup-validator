FROM golang:1.15-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /out/backup-validator .

FROM alpine:3.13 AS bin
EXPOSE 9178
ENTRYPOINT [ "/backup-validator" ]
WORKDIR /workdir

# Install packages
RUN apk add --no-cache ca-certificates restic=0.11.0-r0 && update-ca-certificates

USER 1001
COPY --from=build /out/backup-validator /backup-validator
