package cstx

import (
	"encoding/json"
	"fmt"
)

// NodeFlags are engine-compatible bit constants. Graph APIs accept ordinary
// uint64 masks built from these values.
const (
	FlagNone               uint64 = 0
	FlagHoneypot           uint64 = 1 << 0
	FlagNoise              uint64 = 1 << 1
	FlagFalsePositive      uint64 = 1 << 2
	FlagManualIgnored      uint64 = 1 << 3
	FlagThreatPresent      uint64 = 1 << 4
	FlagHistoricVulnerable uint64 = 1 << 5
	FlagInternal           uint64 = 1 << 6
)

// FlagsAllMask contains every currently defined node flag.
const FlagsAllMask uint64 = FlagHoneypot | FlagNoise | FlagFalsePositive |
	FlagManualIgnored | FlagThreatPresent | FlagHistoricVulnerable | FlagInternal

// FlagsDefaultExcludeMask is the engine's standard default-exclusion mask.
const FlagsDefaultExcludeMask uint64 = FlagHoneypot | FlagNoise | FlagFalsePositive | FlagManualIgnored

// Order is the deterministic ordering applied to a collection cursor.
type Order string

const (
	OrderUnspecified Order = "unspecified"
	OrderIDAsc       Order = "id_asc"
	OrderIDDesc      Order = "id_desc"
)

// Node is the CSTX graph node exchanged with the Rust runtime. Model
// holds schema-typed fields plus the reserved keys "__node_type__" and
// "cstx_flags".
type Node struct {
	ID      string         `json:"id"`
	Type    string         `json:"type"`
	Value   any            `json:"value"`
	Model   map[string]any `json:"model"`
	Sources []string       `json:"sources"`
	Extras  map[string]any `json:"extras"`
}

// MarshalJSON keeps list fields as [] rather than null; the Rust contract
// requires every field to be present.
func (n Node) MarshalJSON() ([]byte, error) {
	type wire Node
	if n.Sources == nil {
		n.Sources = []string{}
	}
	if n.Model == nil {
		n.Model = map[string]any{}
	}
	if n.Extras == nil {
		n.Extras = map[string]any{}
	}
	return json.Marshal(wire(n))
}

// Edge is the CSTX graph relationship exchanged with the Rust runtime.
type Edge struct {
	ID           string         `json:"id"`
	SourceID     string         `json:"source_id"`
	TargetID     string         `json:"target_id"`
	RelationType string         `json:"relation_type"`
	Sources      []string       `json:"sources"`
	Attrs        map[string]any `json:"attrs"`
}

// MarshalJSON keeps list fields as [] rather than null.
func (e Edge) MarshalJSON() ([]byte, error) {
	type wire Edge
	if e.Sources == nil {
		e.Sources = []string{}
	}
	if e.Attrs == nil {
		e.Attrs = map[string]any{}
	}
	return json.Marshal(wire(e))
}

// GraphStats is a small aggregate count summary.
type GraphStats struct {
	Nodes   map[string]int64 `json:"nodes"`
	Edges   map[string]int64 `json:"edges"`
	Sources map[string]int64 `json:"sources"`
}

// Delta counts graph elements changed by commits in one time range.
type Delta struct {
	AddedNodes   uint64 `json:"added_nodes"`
	UpdatedNodes uint64 `json:"updated_nodes"`
	RemovedNodes uint64 `json:"removed_nodes"`
	AddedEdges   uint64 `json:"added_edges"`
	UpdatedEdges uint64 `json:"updated_edges"`
	RemovedEdges uint64 `json:"removed_edges"`
}

// ChangeSet lists the IDs changed by one successfully committed mutation.
type ChangeSet struct {
	AddedNodeIDs   []string `json:"added_node_ids"`
	UpdatedNodeIDs []string `json:"updated_node_ids"`
	RemovedNodeIDs []string `json:"removed_node_ids"`
	AddedEdgeIDs   []string `json:"added_edge_ids"`
	UpdatedEdgeIDs []string `json:"updated_edge_ids"`
	RemovedEdgeIDs []string `json:"removed_edge_ids"`
	Reset          bool     `json:"reset"`
}

// Affected returns the total number of changed graph elements.
func (c ChangeSet) Affected() int {
	return len(c.AddedNodeIDs) + len(c.UpdatedNodeIDs) + len(c.RemovedNodeIDs) +
		len(c.AddedEdgeIDs) + len(c.UpdatedEdgeIDs) + len(c.RemovedEdgeIDs)
}

// Commit describes one immutable repository revision.
type Commit struct {
	ID        string   `json:"id"`
	Parents   []string `json:"parents"`
	Message   string   `json:"message"`
	Metadata  any      `json:"metadata"`
	Stats     Delta    `json:"stats"`
	CreatedAt int64    `json:"created_at"`
}

// GraphDiff groups added, removed, and modified element IDs by element type.
type GraphDiff struct {
	Added    map[string][]string `json:"added"`
	Removed  map[string][]string `json:"removed"`
	Modified map[string][]string `json:"modified"`
	// Truncated reports whether a limit stopped the diff before the last
	// change. An empty group is otherwise ambiguous between "nothing of that
	// type changed" and "the limit ran out first".
	Truncated bool `json:"truncated"`
}

// JoinRuleSpec is the portable native-linker rule shared by all bindings.
type JoinRuleSpec struct {
	LeftType      string  `json:"left_type"`
	RightType     string  `json:"right_type"`
	Relation      string  `json:"relation"`
	LeftKey       string  `json:"left_key"`
	RightKey      string  `json:"right_key"`
	Predicted     bool    `json:"predicted"`
	LeftTargetID  *string `json:"left_target_id,omitempty"`
	RightSourceID *string `json:"right_source_id,omitempty"`
}

// SCOSchemaContract describes one portable node schema.
type SCOSchemaContract struct {
	Schema     any            `json:"schema"`
	ValueField *string        `json:"value_field"`
	Metadata   map[string]any `json:"metadata"`
}

// SROSchemaContract describes one portable relationship schema.
type SROSchemaContract struct {
	Schema   any            `json:"schema"`
	Metadata map[string]any `json:"metadata"`
}

// ParserSchemaContract describes one portable parser input contract.
type ParserSchemaContract struct {
	InputSchema any            `json:"input_schema"`
	Metadata    map[string]any `json:"metadata"`
}

// PluginSchemaContract groups schemas published by one CSTX plugin.
type PluginSchemaContract struct {
	Version string                          `json:"version"`
	SCO     map[string]SCOSchemaContract    `json:"sco"`
	SRO     map[string]SROSchemaContract    `json:"sro"`
	Parsers map[string]ParserSchemaContract `json:"parsers"`
}

// SchemaContract is the atomic schema exchange unit shared by all bindings.
type SchemaContract struct {
	Format  string                          `json:"format"`
	Plugins map[string]PluginSchemaContract `json:"plugins"`
}

// AnchorConcept names one native concept and its member node types.
type AnchorConcept struct {
	Name      string
	NodeTypes []string
}

// UnmarshalJSON decodes the Rust (name, node_types) tuple transport.
func (c *AnchorConcept) UnmarshalJSON(data []byte) error {
	var pair struct {
		Name      string
		NodeTypes []string
	}
	var wire []json.RawMessage
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if len(wire) != 2 {
		return fmt.Errorf("cstx: anchor concept must be a two-item tuple")
	}
	if err := json.Unmarshal(wire[0], &pair.Name); err != nil {
		return err
	}
	if err := json.Unmarshal(wire[1], &pair.NodeTypes); err != nil {
		return err
	}
	c.Name, c.NodeTypes = pair.Name, pair.NodeTypes
	return nil
}

// Ref is one named repository reference and its commit ID.
type Ref struct {
	Name string
	Head string
}

// UnmarshalJSON decodes the Rust (name, head) tuple transport.
func (r *Ref) UnmarshalJSON(data []byte) error {
	var pair [2]string
	if err := json.Unmarshal(data, &pair); err != nil {
		return err
	}
	r.Name, r.Head = pair[0], pair[1]
	return nil
}
