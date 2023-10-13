IMAGE ?= ghcr.io/suffiks/examples/wasi/ingress-as:v0.0.6

build/release.wasm: $(wildcard assembly/*.ts)
	npm run asbuild:release

build-release: build/release.wasm

.PHONY: build-debug
test: build-release
	extgen wasi test --ext ./k8s-resources.yaml -v build/release.wasm

.PHONY: publish
publish: test
	extgen wasi publish --tag ${IMAGE} --docs ./docs build/release.wasm
