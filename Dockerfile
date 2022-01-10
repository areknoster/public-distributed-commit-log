FROM golang:1.17-alpine AS builder
WORKDIR /src/
COPY . ./
RUN mkdir -p bin/linux
RUN go build -o bin/linux ./cmd/...


FROM alpine:3.15.0

WORKDIR /app

COPY --from=builder /src/bin/linux/acceptance-sentinel ./

EXPOSE 8000

RUN addgroup -S nonroot && adduser -S -G nonroot nonroot

USER nonroot:nonroot

ENTRYPOINT ["./acceptance-sentinel"]
