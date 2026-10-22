# syntax=docker/dockerfile:1.27

ARG GO_VERSION=1.27.3
ARG ALPINE_VERSION=3.24

FROM golang:${GO_VERSION}-alpine AS build
RUN apk add --no-cache --upgrade ca-certificates=20261020-r0 ca-certificates-bundle=20261020-r0
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/songsee ./cmd/songsee

FROM alpine:${ALPINE_VERSION}
RUN apk add --no-cache --upgrade ca-certificates=20261020-r0 ca-certificates-bundle=20261020-r0
RUN apk add --no-cache ca-certificates ffmpeg tzdata \
    && adduser -D -u 10001 -h /home/songsee songsee \
    && mkdir -p /input /output \
    && chown -R songsee:songsee /input /output
VOLUME ["/input", "/output"]
WORKDIR /output
COPY --from=build /out/songsee /usr/local/bin/songsee
USER songsee
ENTRYPOINT ["songsee"]
CMD ["--help"]
