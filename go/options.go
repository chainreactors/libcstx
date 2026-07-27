package cstx

import "encoding/json"

// DefaultCursorPageSize matches the Rust runtime default.
const DefaultCursorPageSize = 1024

// NodeFilter pushes node selection into Rust before values are materialized.
type NodeFilter struct {
	Types     []string `json:"types"`
	IDs       []string `json:"ids"`
	Sources   []string `json:"sources"`
	FlagsAll  uint64   `json:"flags_all"`
	FlagsAny  uint64   `json:"flags_any"`
	FlagsNone uint64   `json:"flags_none"`
}

func (f NodeFilter) MarshalJSON() ([]byte, error) {
	type wire NodeFilter
	f.Types = orEmpty(f.Types)
	f.IDs = orEmpty(f.IDs)
	f.Sources = orEmpty(f.Sources)
	return json.Marshal(wire(f))
}

// EdgeFilter pushes relationship selection into Rust.
type EdgeFilter struct {
	SourceID  string   `json:"source_id"`
	TargetID  string   `json:"target_id"`
	Relations []string `json:"relations"`
	Sources   []string `json:"sources"`
}

func (f EdgeFilter) MarshalJSON() ([]byte, error) {
	// source_id/target_id are Option<String> in Rust: absent means null,
	// and Some("") would filter on an empty ID instead of disabling it.
	type wire struct {
		SourceID  *string  `json:"source_id"`
		TargetID  *string  `json:"target_id"`
		Relations []string `json:"relations"`
		Sources   []string `json:"sources"`
	}
	return json.Marshal(wire{
		SourceID:  nonEmpty(f.SourceID),
		TargetID:  nonEmpty(f.TargetID),
		Relations: orEmpty(f.Relations),
		Sources:   orEmpty(f.Sources),
	})
}

// CollectionOptions bounds cursor materialization and ordering.
// The zero value uses the runtime defaults (no limit, page size 1024,
// unspecified order).
type CollectionOptions struct {
	Limit    *int  `json:"limit"`
	PageSize int   `json:"page_size"`
	Order    Order `json:"order"`
}

func (o CollectionOptions) normalize() CollectionOptions {
	if o.PageSize <= 0 {
		o.PageSize = DefaultCursorPageSize
	}
	if o.Order == "" {
		o.Order = OrderUnspecified
	}
	return o
}

// QueryOptions applies collection semantics to graph DSL queries.
type QueryOptions struct {
	Collection CollectionOptions `json:"collection"`
}

func orEmpty(value []string) []string {
	if value == nil {
		return []string{}
	}
	return value
}

func nonEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
