FROM gliderlabs/alpine

RUN apk add --no-cache ca-certificates

copy ./datadog-metrics-tuner /

VOLUME /metrics.d
WORKDIR /

ENTRYPOINT ["/datadog-metrics-tuner"]
