//go:build !cstx_native

package cstx

import (
	"context"
	"errors"
	"testing"
)

func TestOpenWithoutNativeTag(t *testing.T) {
	rt, err := Open(context.Background(), Config{})
	if !errors.Is(err, ErrNativeUnavailable) {
		t.Fatalf("expected ErrNativeUnavailable, got rt=%v err=%v", rt, err)
	}
}

func TestOpenHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	rt, err := Open(ctx, Config{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got rt=%v err=%v", rt, err)
	}
}
