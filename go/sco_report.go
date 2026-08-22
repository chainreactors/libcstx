package cstx

import (
	"encoding/json"
)

// ReportNode is the libcstx SCO used by a report Artifact Harness. A report
// is a graph node, not a Cairn persistence model: references remain explicit
// IDs and the natural-language body is kept as a JSON document for renderers.
// The graph relation (typically includes → sarif_vuln) is stored separately by
// Graph.AddEdges.
type ReportNode struct {
	nodeHeader
	Value           string            `json:"value,omitempty"`
	Title           string            `json:"title,omitempty"`
	Summary         string            `json:"summary,omitempty"`
	Body            json.RawMessage   `json:"body,omitempty"`
	TaskID          string            `json:"task_id,omitempty"`
	ArtifactVersion int               `json:"artifact_version,omitempty"`
	ArtifactStatus  string            `json:"artifact_status,omitempty"`
	FactIDs         []string          `json:"fact_ids,omitempty"`
	EvidenceIDs     []string          `json:"evidence_ids,omitempty"`
	EvidenceClaims  map[string]string `json:"evidence_claims,omitempty"`
	ObservationIDs  []string          `json:"observation_ids,omitempty"`
	AssetIDs        []string          `json:"asset_ids,omitempty"`
	// VulnerabilityIDs are the report's semantic membership. The graph stores
	// the corresponding includes edges; keeping the IDs on the typed SCO makes
	// the extension self-describing without reintroducing an Artifact entity.
	VulnerabilityIDs []string       `json:"vulnerability_ids,omitempty"`
	Provenance       map[string]any `json:"provenance,omitempty"`
}

func parseReportNode(data []byte) (SCONode, error) {
	var node ReportNode
	if err := json.Unmarshal(data, &node); err != nil {
		return nil, err
	}
	return &node, nil
}

func init() {
	// Report is an extension SCO, so it is registered through the same parser
	// mechanism available to downstream libcstx extensions.
	_ = RegisterSCOParser("report", parseReportNode)
}
