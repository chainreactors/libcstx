package cstx

import (
	"encoding/json"
	"testing"
)

func TestParseSCONodeReport(t *testing.T) {
	node, err := ParseSCONode([]byte(`{"cstx_type":"report","cstx_id":"report:t1","title":"Task report","body":{"markdown":"# report"},"artifact_ids":["vuln:v1"]}`))
	if err != nil {
		t.Fatal(err)
	}
	report, ok := node.(*ReportNode)
	if !ok {
		t.Fatalf("ParseSCONode() = %T, want *ReportNode", node)
	}
	if report.CstxType() != "report" || report.CstxID() != "report:t1" || report.Title != "Task report" {
		t.Fatalf("unexpected report node: %+v", report)
	}
	var body map[string]string
	if err := json.Unmarshal(report.Body, &body); err != nil || body["markdown"] != "# report" {
		t.Fatalf("unexpected report body: %s (%v)", report.Body, err)
	}
}

func TestRegisterSCOParserRejectsDuplicate(t *testing.T) {
	parser := func([]byte) (SCONode, error) { return nil, nil }
	if err := RegisterSCOParser("test_extension", parser); err != nil {
		t.Fatal(err)
	}
	if err := RegisterSCOParser("test_extension", parser); err == nil {
		t.Fatal("duplicate parser registration succeeded")
	}
}
