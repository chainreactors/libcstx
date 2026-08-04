package cstx

import "encoding/json"

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

// Node is the canonical graph node exchanged with the Rust runtime. Model
// holds schema-typed fields plus the reserved keys "__node_type__",
// "cstx_flags", "created_at", and "updated_at".
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

// Edge is the canonical graph relationship exchanged with the Rust runtime.
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

// Commit describes one immutable graph tree on a repository ref.
type Commit struct {
	ID        string `json:"id"`
	Tree      string `json:"tree"`
	Parent    string `json:"parent"`
	Message   string `json:"message"`
	Metadata  any    `json:"metadata"`
	Stats     any    `json:"stats"`
	CreatedAt int64  `json:"created_at"`
}

// GraphDiff groups added, removed, and modified element IDs by element type.
type GraphDiff struct {
	Added    map[string][]string `json:"added"`
	Removed  map[string][]string `json:"removed"`
	Modified map[string][]string `json:"modified"`
}

// CASIndex contains canonical graph entries and their serialized immutable
// objects. Object bytes remain raw so storage adapters do not deserialize and
// reserialize data at the language boundary.
type CASIndex struct {
	Entries map[string]map[string]string `json:"entries"`
	Objects map[string]json.RawMessage   `json:"objects"`
}

// CommitDeltaRequest commits pre-indexed element hashes and removals onto a ref.
// The referenced objects must already exist in the repository object store.
type CommitDeltaRequest struct {
	DeltaEntries   map[string]map[string]string `json:"delta_entries"`
	RemovedNodeIDs []string                     `json:"removed_node_ids"`
	RemovedEdgeIDs []string                     `json:"removed_edge_ids"`
	RefName        string                       `json:"ref_name"`
	Message        string                       `json:"message"`
	Metadata       any                          `json:"metadata"`
	CreatedAt      int64                        `json:"created_at"`
}

// TreeEntryChange contains the base and head content hashes for one element.
// A nil side means the element is absent from that tree.
type TreeEntryChange [2]*string

// TreeEntryDiff groups raw content-hash changes by element type and ID.
type TreeEntryDiff map[string]map[string]TreeEntryChange

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
