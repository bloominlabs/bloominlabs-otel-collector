# https://github.com/hashicorp/vault/pull/12358
VERSION 0.6
FROM golang:1.20
WORKDIR /bloominlabs-otel-collector

deps:
  RUN GO111MODULE=on go install go.opentelemetry.io/collector/cmd/builder@v0.80.0
  RUN GO111MODULE=on go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/mdatagen@v0.80.0

certs:
  RUN curl -k https://vault.prod.stratos.host:8200/v1/internal/ca/pem > /etc/ssl/certs/internal.pem
  SAVE ARTIFACT /etc/ssl/certs/internal.pem

mc-monitor:
  FROM itzg/mc-monitor:0.11.2
  SAVE ARTIFACT mc-monitor

generate:
  FROM +deps

  COPY otelcol-builder.yaml . 
  COPY ./processor/ ./processor/
  COPY ./receiver/vaultkvreceiver+receiver/vaultkvreceiver ./receiver/vaultkvreceiver/
  RUN GO111MODULE=on CGO_ENABLED=0 builder --output-path . --config otelcol-builder.yaml --name bloominlabs-otel-collector --skip-compilation
  SAVE ARTIFACT *.go AS LOCAL ./
  SAVE ARTIFACT ./receiver/vaultkvreceiver/ AS LOCAL ./receiver/vaultkvreceiver/
  SAVE ARTIFACT go.mod AS LOCAL go.mod
  SAVE ARTIFACT go.sum AS LOCAL go.sum

build:
  FROM +generate
  RUN go mod download
  RUN CGO_ENABLED=0 go build -ldflags="-s -w -extldflags '-static'" -trimpath -o bloominlabs-otel-collector
  SAVE ARTIFACT ./bloominlabs-otel-collector AS LOCAL ./bin/bloominlabs-otel-collector

docker:
  FROM --platform=linux/amd64 gcr.io/distroless/base-debian11:debug-nonroot
  COPY +mc-monitor/mc-monitor /usr/bin/mc-monitor
  COPY +certs/internal.pem /etc/ssl/certs/internal.pem
  COPY +build/bloominlabs-otel-collector .
  ENTRYPOINT ["./bloominlabs-otel-collector"]
  SAVE IMAGE infrastructure-otel-collector:latest
  SAVE IMAGE --push ghcr.io/bloominlabs/otel-collector:latest

# if we need to break it out into multiple different distributions
# processor:
#   FROM +deps
#   COPY ./processor/ ./processor/
# 
#   SAVE ARTIFACT ./processor/
# 
# receiver:
#   FROM +deps 
#   COPY ./receiver/vaultkvreceiver+receiver/vaultkvreceiver ./receiver/vaultkvreceiver/
#   SAVE ARTIFACT ./receiver/vaultkvreceiver/
