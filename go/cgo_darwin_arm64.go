//go:build cstx_native && darwin && arm64

package cstx

// #cgo LDFLAGS: -L${SRCDIR}/lib/darwin_arm64 -lcstx_ffi -lm -framework Security -framework CoreFoundation
// #cgo CFLAGS: -I${SRCDIR}
import "C"
