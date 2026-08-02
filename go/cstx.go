package cstx

import (
	"context"
	"errors"
	"sync"
)

// Config configures an in-memory CSTX runtime.
type Config struct {
	// ProjectID namespaces the runtime; defaults to "default".
	ProjectID string
	// CursorPageSize bounds incremental cursor materialization; defaults to
	// DefaultCursorPageSize.
	CursorPageSize int
}

func (c Config) normalize() Config {
	if c.ProjectID == "" {
		c.ProjectID = "default"
	}
	if c.CursorPageSize <= 0 {
		c.CursorPageSize = DefaultCursorPageSize
	}
	return c
}

// CSTX is the single owner of shared schema, graph, and repository
// state. Create it with Open and release it with Close.
type CSTX struct {
	eng       engine
	projectID string

	// Schemas, Graph, and Repo are lightweight namespaces sharing
	// this runtime's state.
	Schemas *Schemas
	Graph   *Graph
	Repo    *Repository
	Raw     *Raw

	mu     sync.Mutex
	closed bool
}

// Open creates an in-memory native runtime.
func Open(ctx context.Context, config Config) (*CSTX, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	config = config.normalize()
	eng, err := newEngine(config)
	if err != nil {
		return nil, err
	}
	return wrapRuntime(eng, config.ProjectID), nil
}

func wrapRuntime(eng engine, projectID string) *CSTX {
	rt := &CSTX{eng: eng, projectID: projectID}
	rt.Schemas = &Schemas{eng: eng}
	rt.Graph = &Graph{eng: eng}
	rt.Repo = &Repository{eng: eng}
	if raw, ok := eng.(rawEngine); ok {
		rt.Raw = &Raw{eng: raw}
	}
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
func (c *CSTX) LastChange(ctx context.Context) (ChangeSet, error) {
	if err := contextError(ctx); err != nil {
		return ChangeSet{}, err
	}
	return c.eng.lastChange(ctx)
}
