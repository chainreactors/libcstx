// Package cstx is the typed Go facade for the CSTX native runtime.
//
// The Rust core is the sole owner of graph, repository, validation,
// and ingest semantics (issue #6). This package only converts between Go
// domain values and the JSON transport used by the cstx-ffi C boundary; it
// never reimplements business behavior.
//
// The engine behind Open always uses the bundled cstx-ffi static library.
// Consumers therefore build this package with CGO enabled. Prebuilt libraries
// live under lib/<platform>/ and are synced from the cstx release pipeline.
package cstx
