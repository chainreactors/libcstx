package cstx

import "context"

// NodeCursor is a lazy node iterator that avoids materializing a full result
// set. Callers must Close it; Close is idempotent. A mutation that
// invalidates the cursor surfaces through Err as CodeCursorInvalidated.
type NodeCursor struct {
	iter nodeIter
	node Node
	err  error
	done bool
}

// Next advances the cursor. It returns false at exhaustion or on error;
// check Err after the loop.
func (c *NodeCursor) Next(ctx context.Context) bool {
	if c.done {
		return false
	}
	if err := contextError(ctx); err != nil {
		c.err = err
		c.done = true
		return false
	}
	node, ok, err := c.iter.next(ctx)
	if err != nil {
		c.err = err
		c.done = true
		return false
	}
	if !ok {
		c.done = true
		return false
	}
	c.node = node
	return true
}

// Node returns the current node. It is valid only after Next returned true.
func (c *NodeCursor) Node() Node { return c.node }

// Err returns the first error encountered during iteration.
func (c *NodeCursor) Err() error { return c.err }

// Close releases cursor state; repeated calls are safe.
func (c *NodeCursor) Close() error {
	if !c.done {
		c.done = true
		c.iter.close()
	}
	return nil
}

// EdgeCursor is a lazy relationship iterator.
type EdgeCursor struct {
	iter edgeIter
	edge Edge
	err  error
	done bool
}

func (c *EdgeCursor) Next(ctx context.Context) bool {
	if c.done {
		return false
	}
	if err := contextError(ctx); err != nil {
		c.err = err
		c.done = true
		return false
	}
	edge, ok, err := c.iter.next(ctx)
	if err != nil {
		c.err = err
		c.done = true
		return false
	}
	if !ok {
		c.done = true
		return false
	}
	c.edge = edge
	return true
}

// Edge returns the current relationship.
func (c *EdgeCursor) Edge() Edge { return c.edge }

// Err returns the first error encountered during iteration.
func (c *EdgeCursor) Err() error { return c.err }

// Close releases cursor state; repeated calls are safe.
func (c *EdgeCursor) Close() error {
	if !c.done {
		c.done = true
		c.iter.close()
	}
	return nil
}
