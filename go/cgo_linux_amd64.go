//go:build linux && amd64

package cstx

// #cgo LDFLAGS: -L${SRCDIR}/lib/linux_amd64 -lcstx_ffi -lm -ldl -lpthread
// #cgo CFLAGS: -I${SRCDIR}
import "C"
