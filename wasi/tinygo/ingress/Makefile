IMAGE ?= ghcr.io/suffiks/examples/wasi/ingress-tinygo:v0.0.2

build/release.wasm: $(wildcard *.go) go.mod go.sum
	mkdir -p build
	CGO_ENABLED=0 tinygo build -o $@ -scheduler=none -target=wasi .

.PHONY: build-release
build-release: build/release.wasm

.PHONY: test
test: build-release
	extgen wasi test --ext ./k8s-resources.yaml -v build/release.wasm

.PHONY: publish
publish: test
	extgen wasi publish --tag ${IMAGE} --docs ./docs build/release.wasm
