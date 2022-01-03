FROM alpine:3.15.0

WORKDIR /app

COPY linux/acceptance-sentinel .

EXPOSE 8080
EXPOSE 8000
EXPOSE 5001

RUN addgroup -S nonroot && adduser -S -G nonroot nonroot

USER nonroot:nonroot

ENTRYPOINT ["./acceptance-sentinel"]
