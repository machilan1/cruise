# Build the Go Binary.
FROM golang:1.24-alpine3.21 AS build
ARG BUILD_REF

WORKDIR /service

COPY go.* ./
COPY cmd ./cmd
COPY internal ./internal
COPY vendor ./vendor

RUN go build -ldflags "-X main.build=${BUILD_REF} -extldflags=-static" -o main ./cmd/api

# Run the Go Binary in Alpine.
FROM alpine:3.21
ARG BUILD_DATE
ARG BUILD_REF

WORKDIR /service

RUN addgroup -g 1000 -S app && \
    adduser -u 1000 -h /service -G app -S app

USER app

COPY --from=build --chown=app:app /service/main ./main

CMD ["./main"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="api" \
      org.opencontainers.image.authors="Mayainfo <mayainfo.co.ltd@gmail.com>" \
      org.opencontainers.image.source="https://github.com/mayainfo" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="Mayainfo"
