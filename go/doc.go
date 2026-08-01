// Package cstx is the typed Go facade for the CSTX v0.3 native runtime.
//
// The Rust core is the sole owner of graph, repository, validation,
// and ingest semantics (issue #6). This package only converts between Go
// domain values and the JSON transport used by the cstx-ffi C boundary; it
// never reimplements business behavior.
//
// The engine behind Open requires the cstx-ffi static library and is enabled
// with the cstx_native build tag:
//
//	go build -tags cstx_native
//
// Without the tag the package still compiles and all types remain available,
// but Open returns ErrNativeUnavailable. Prebuilt libraries live under
// lib/<platform>/ and are synced from the cstx release pipeline.
package cstx
