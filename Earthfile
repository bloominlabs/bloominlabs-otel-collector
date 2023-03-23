# https://github.com/hashicorp/vault/pull/12358
VERSION 0.6
FROM golang:1.20
WORKDIR /bloominlabs-otel-collector

deps:
  RUN GO111MODULE=on go install go.opentelemetry.io/collector/cmd/builder@v0.74.0
  RUN GO111MODULE=on go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/mdatagen@v0.74.0

processor:
  FROM +deps
  COPY ./processor/ ./processor/

  SAVE ARTIFACT ./processor/

receiver:
  FROM +deps 
  COPY ./receiver/vaultkvreceiver+receiver/vaultkvreceiver ./receiver/vaultkvreceiver/
  SAVE ARTIFACT ./receiver/vaultkvreceiver/
