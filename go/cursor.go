package cstx

import (
	"context"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

// CursorKind identifies the row shape returned by a GraphCursor.
type CursorKind string

const (
	CursorKindNodes         CursorKind = "nodes"
	CursorKindRelationships CursorKind = "relationships"
	CursorKindComponents    CursorKind = "components"
	CursorKindNodeScores    CursorKind = "node_scores"
	CursorKindNodePairs     CursorKind = "node_pairs"
	CursorKindCycles        CursorKind = "cycles"
	CursorKindPaths         CursorKind = "paths"
	CursorKindCommunities   CursorKind = "communities"
)

// GraphCursor is the single cursor type for nodes, relationships, queries and graph
// analysis results. Page uses a one-based page number and never reruns the
// operation that created the cursor.
type GraphCursor struct {
	inner graphCursor
	kind  CursorKind
	done  bool
}

// Kind returns the logical row shape emitted by this cursor.
func (c *GraphCursor) Kind() CursorKind { return c.kind }

// Page returns the generated protobuf page from the native boundary. The
// caller can inspect the result oneof directly, avoiding an intermediate DTO
// and any JSON conversion.
func (c *GraphCursor) Page(ctx context.Context, limit, page int) (*cstxproto.GraphResultPage, error) {
	if c.done || c.inner == nil {
		return nil, &Error{Code: CodeInvalidArgument, Operation: "cursor.page", Message: "cursor is closed"}
	}
	if err := contextError(ctx); err != nil {
		return nil, err
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
