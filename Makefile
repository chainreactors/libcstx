.PHONY: go-test go-build ffi-local

# Build FFI library from local cstx-core source (dev mode)
CSTX_CORE ?= $(HOME)/cstx/cstx-core

ffi-local:
	cargo build --release -p cstx-ffi --manifest-path $(CSTX_CORE)/Cargo.toml
	mkdir -p go/lib/linux_amd64
	cp $(CSTX_CORE)/target/release/libcstx_ffi.a go/lib/linux_amd64/

go-build:
	cd go && CGO_ENABLED=1 go build -tags cstx_native ./...

go-test: go-build
	cd go && CGO_ENABLED=1 go test -tags cstx_native -v ./...
