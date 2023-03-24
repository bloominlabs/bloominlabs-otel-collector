VERSION=0.0.9

pushd ./infrastructure 

earthly +build && \
  cp ./bin/bloominlabs-otel-collector ./bin/bloominlabs-otel-collector-linux-amd64 && \
  sha256sum ./bin/bloominlabs-otel-collector-linux-amd64 > ./bin/SHA256SUMS && \
  gh release -R bloominlabs/bloominlabs-otel-collector create --title v$VERSION v$VERSION ./bin/*

popd 
