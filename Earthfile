# https://github.com/hashicorp/vault/pull/12358
VERSION 0.6
FROM golang:1.20
WORKDIR /bloominlabs-otel-collector

tools:
  RUN GO111MODULE=on go install go.opentelemetry.io/collector/cmd/builder@v0.80.0
  RUN GO111MODULE=on go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/mdatagen@v0.80.0
  SAVE ARTIFACT /go/bin/builder
  SAVE ARTIFACT /go/bin/mdatagen

files:
  COPY ./processor/ ./processor/
  COPY ./receiver/vaultkvreceiver+receiver/vaultkvreceiver ./receiver/vaultkvreceiver/
  COPY ./receiver/digitaloceanreceiver+receiver/digitaloceanreceiver ./receiver/digitaloceanreceiver/

deps:
  FROM +files
  COPY go.mod go.sum ./
  RUN go mod download
  # Output these back in case go mod download changes them.
  SAVE ARTIFACT go.mod AS LOCAL go.mod
  SAVE ARTIFACT go.sum AS LOCAL go.sum

certs:
  RUN curl -k https://vault.prod.stratos.host:8200/v1/internal/ca/pem > /etc/ssl/certs/internal.pem
  SAVE ARTIFACT /etc/ssl/certs/internal.pem

generate:
  FROM +deps
  COPY +tools/builder /go/bin/builder
  COPY otelcol-builder.yaml . 
  RUN GO111MODULE=on CGO_ENABLED=0 builder --output-path . --config otelcol-builder.yaml --name bloominlabs-otel-collector --skip-compilation
  SAVE ARTIFACT *.go AS LOCAL ./
  SAVE ARTIFACT go.mod AS LOCAL go.mod
  SAVE ARTIFACT go.sum AS LOCAL go.sum

build:
  FROM +generate
  RUN CGO_ENABLED=0 go build -ldflags="-s -w -extldflags '-static'" -trimpath -o bloominlabs-otel-collector
  SAVE ARTIFACT ./bloominlabs-otel-collector AS LOCAL ./bin/bloominlabs-otel-collector

release:
  FROM +build
  BUILD +docker
  RUN tar -cvzf bloominlabs-otel-collector_linux_amd64.tar.gz ./bloominlabs-otel-collector
  RUN sha256sum *.tar.gz > ./SHA256SUMS
  SAVE ARTIFACT ./SHA256SUMS AS LOCAL ./bin/SHA256SUMS
  SAVE ARTIFACT *.tar.gz AS LOCAL ./bin/

docker:
  FROM --platform=linux/amd64 gcr.io/distroless/base-debian11:debug-nonroot
  COPY +certs/internal.pem /etc/ssl/certs/internal.pem
  COPY +build/bloominlabs-otel-collector .
  ENTRYPOINT ["./bloominlabs-otel-collector"]
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
