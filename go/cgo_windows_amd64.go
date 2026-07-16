//go:build windows && amd64

package cstx

// #cgo LDFLAGS: -L${SRCDIR}/lib/windows_amd64 -lcstx_ffi -lm -lpthread -lws2_32 -luserenv -lbcrypt -lntdll
// #cgo CFLAGS: -I${SRCDIR}
import "C"
