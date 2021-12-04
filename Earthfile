# https://github.com/hashicorp/vault/pull/12358
VERSION 0.6
FROM golang:1.17
WORKDIR /bloominlabs-otel-collector

deps:
  RUN GO111MODULE=on go install go.opentelemetry.io/collector/cmd/builder@v0.40.0

build:
  FROM +deps
  COPY processor processor
  COPY otelcol-builder.yaml . 
  RUN builder --output-path . --config otelcol-builder.yaml --name bloominlabs-otel-collector  
  SAVE ARTIFACT ./bloominlabs-otel-collector AS LOCAL ./bin/bloominlabs-otelcol
  # Output these back in case go mod download changes them.
  SAVE ARTIFACT *.go AS LOCAL ./
	SAVE ARTIFACT go.mod AS LOCAL go.mod
	SAVE ARTIFACT go.sum AS LOCAL go.sum

docker:
  COPY +build/bloominlabs-otel-collector .
  ENTRYPOINT ["./bloominlabs-otel-collector"]
  SAVE IMAGE otel-collector:latest
  SAVE IMAGE --push ghcr.io/bloominlabs/otel-collector:latest
