# https://github.com/hashicorp/vault/pull/12358
VERSION 0.6
FROM golang:1.18
WORKDIR /bloominlabs-otel-collector

deps:
  RUN GO111MODULE=on go install go.opentelemetry.io/collector/cmd/builder@v0.54.0

certs:
  RUN curl -k https://vault.prod.stratos.host:8200/v1/internal/ca/pem > /etc/ssl/certs/internal.pem
  SAVE ARTIFACT /etc/ssl/certs/internal.pem

mc-monitor:
  FROM itzg/mc-monitor:0.10.6
  SAVE ARTIFACT mc-monitor

build:
  FROM +deps
  COPY processor processor
  COPY receiver receiver
  COPY otelcol-builder.yaml . 
  RUN GO111MODULE=on CGO_ENABLED=0 builder --output-path . --config otelcol-builder.yaml --name bloominlabs-otel-collector --skip-compilation
  SAVE ARTIFACT *.go AS LOCAL ./
	SAVE ARTIFACT go.mod AS LOCAL go.mod
	SAVE ARTIFACT go.sum AS LOCAL go.sum
  RUN go mod download
  RUN CGO_ENABLED=0 go build -ldflags="-s -w -extldflags '-static'" -trimpath -o bloominlabs-otel-collector
  SAVE ARTIFACT ./bloominlabs-otel-collector AS LOCAL ./bin/bloominlabs-otel-collector

docker:
  FROM alpine:latest
  RUN apk add --update ca-certificates
  COPY +mc-monitor/mc-monitor /usr/bin/mc-monitor
  COPY +certs/internal.pem /etc/ssl/certs/internal.pem
  COPY +build/bloominlabs-otel-collector .
  ENTRYPOINT ["./bloominlabs-otel-collector"]
  SAVE IMAGE otel-collector:latest
  SAVE IMAGE --push ghcr.io/bloominlabs/otel-collector:latest
