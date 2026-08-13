package cstx

import (
	"context"
	"encoding/json"
	"fmt"
)

// CursorKind identifies the row shape returned by a GraphCursor.
type CursorKind string

const (
	CursorKindNodes       CursorKind = "nodes"
	CursorKindEdges       CursorKind = "edges"
	CursorKindComponents  CursorKind = "components"
	CursorKindNodeScores  CursorKind = "node_scores"
	CursorKindNodePairs   CursorKind = "node_pairs"
	CursorKindCycles      CursorKind = "cycles"
	CursorKindPaths       CursorKind = "paths"
	CursorKindCommunities CursorKind = "communities"
)

// CursorPage is one bounded, one-based page from a native graph result.
// Items stay as raw JSON until the caller chooses the domain type, avoiding a
// second in-memory graph or eager decoding of rows outside the requested page.
type CursorPage struct {
	Items   []json.RawMessage `json:"items"`
	Page    int               `json:"page"`
	Limit   int               `json:"limit"`
	HasNext bool              `json:"has_next"`
	Total   *uint64           `json:"total,omitempty"`
	Summary json.RawMessage   `json:"summary,omitempty"`
}

// Nodes decodes this page as CSTX nodes.
func (p CursorPage) Nodes() ([]Node, error) {
	items := make([]Node, len(p.Items))
	for index, item := range p.Items {
		if err := json.Unmarshal(item, &items[index]); err != nil {
			return nil, fmt.Errorf("cstx: decode node at page index %d: %w", index, err)
		}
	}
	return items, nil
}

// Edges decodes this page as CSTX relationships.
func (p CursorPage) Edges() ([]Edge, error) {
	items := make([]Edge, len(p.Items))
	for index, item := range p.Items {
		if err := json.Unmarshal(item, &items[index]); err != nil {
			return nil, fmt.Errorf("cstx: decode edge at page index %d: %w", index, err)
		}
	}
	return items, nil
}

// GraphCursor is the single cursor type for nodes, edges, queries and graph
// analysis results. Page uses a one-based page number and never reruns the
// operation that created the cursor.
type GraphCursor struct {
	inner graphCursor
	kind  CursorKind
	done  bool
}

// Kind returns the logical row shape emitted by this cursor.
func (c *GraphCursor) Kind() CursorKind { return c.kind }

// Page materializes one bounded page.
func (c *GraphCursor) Page(ctx context.Context, limit, page int) (CursorPage, error) {
	if c.done || c.inner == nil {
		return CursorPage{}, &Error{Code: CodeInvalidArgument, Operation: "cursor.page", Message: "cursor is closed"}
	}
	if err := contextError(ctx); err != nil {
		return CursorPage{}, err
	}
	return c.inner.page(ctx, limit, page)
}

// Close releases cursor state; repeated calls are safe.
func (c *GraphCursor) Close() error {
	if !c.done {
		c.done = true
		c.inner.close()
	}
	return nil
}
