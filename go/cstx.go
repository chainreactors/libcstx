package cstx

import (
	"context"
	"errors"
	"sync"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

type runtimeConfig struct {
	projectID      string
	cursorPageSize int
	// Which spelling of a node payload reads return. A program holding no
	// generated type for a node needs the value spelling, so this has to
	// survive normalization rather than being dropped here.
	payloadFormat cstxproto.PayloadFormat
}

func normalizeRuntimeConfig(value *cstxproto.RuntimeConfig) runtimeConfig {
	config := runtimeConfig{projectID: "default", cursorPageSize: DefaultCursorPageSize}
	if value == nil {
		return config
	}
	if value.ProjectId != "" {
		config.projectID = value.ProjectId
	}
	if value.CursorPageSize > 0 {
		config.cursorPageSize = int(value.CursorPageSize)
	}
	config.payloadFormat = value.PayloadFormat
	return config
}

// CSTX is the single owner of shared schema, graph, and repository
// state. Create it with Open and release it with Close.
type CSTX struct {
	eng       engine
	projectID string

	// Extensions, Graph, and Repo are lightweight namespaces sharing
	// this runtime's state.
	Extensions *Extensions
	Graph      *Graph
	Repo       *Repository

	mu     sync.Mutex
	closed bool
}

// Open creates an in-memory native runtime.
func Open(ctx context.Context, value *cstxproto.RuntimeConfig) (*CSTX, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	config := normalizeRuntimeConfig(value)
	eng, err := newEngine(config)
	if err != nil {
		return nil, err
	}
	return wrapRuntime(eng, config.projectID), nil
}

func wrapRuntime(eng engine, projectID string) *CSTX {
	rt := &CSTX{eng: eng, projectID: projectID}
	rt.Extensions = &Extensions{eng: eng}
	rt.Graph = &Graph{eng: eng}
	rt.Repo = &Repository{eng: eng}
	return rt
}

func contextError(ctx context.Context) error {
	if ctx == nil {
		return errors.New("cstx: nil context")
	}
	return ctx.Err()
}

// ProjectID returns the immutable project namespace supplied at Open.
func (c *CSTX) ProjectID() string { return c.projectID }

// Close releases shared state and invalidates retained services and cursors.
// Repeated calls are safe.
func (c *CSTX) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	return c.eng.close()
}

// Closed reports whether Close has been called.
func (c *CSTX) Closed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

// LastChange returns the IDs changed by the most recent committed mutation.
func (c *CSTX) LastChange(ctx context.Context) (*cstxproto.GraphChangeSet, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return c.eng.lastChange(ctx)
}
