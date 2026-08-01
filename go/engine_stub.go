//go:build !cstx_native

package cstx

// newEngine without the cstx_native build tag: the package compiles and all
// types remain usable, but there is no runtime behind it.
func newEngine(Config) (engine, error) {
	return nil, ErrNativeUnavailable
}
