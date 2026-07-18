//go:build cgo

package cstx

import "testing"

func TestGraphLifecycle(t *testing.T) {
	g := New()
	defer g.Close()

	n := g.LoadAllPlugins()
	if n == 0 {
		t.Fatal("expected at least 1 plugin loaded")
	}
	t.Logf("loaded %d plugins", n)

	if g.NodeCount() != 0 {
		t.Error("new graph should have 0 nodes")
	}
}

func TestGraphAddNode(t *testing.T) {
	g := New()
	defer g.Close()
	g.LoadAllPlugins()

	err := g.AddNode("ip", "ip:10.0.0.1", `{"ip":"10.0.0.1","country":"CN"}`, []string{"test"}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if g.NodeCount() != 1 {
		t.Errorf("expected 1 node, got %d", g.NodeCount())
	}
	if !g.ContainsNode("ip:10.0.0.1") {
		t.Error("node should exist")
	}
}

func TestGraphIngestNative(t *testing.T) {
	g := New()
	defer g.Close()
	g.LoadAllPlugins()

	data := []byte(`{"host":"www.example.com","a":["1.2.3.4"],"ttl":300}` + "\n")
	result, err := g.IngestNative("easm", "dnsx", data)
	if err != nil {
		t.Fatal(err)
	}
	if result.RecordsParsed == 0 {
		t.Error("expected records_parsed > 0")
	}
	t.Logf("ingest: %+v", result)
}

func TestTransform(t *testing.T) {
	data := []byte(`{"host":"www.example.com","a":["1.2.3.4"],"ttl":300}` + "\n")
	nodes, err := Parse("dnsx", data)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) == 0 {
		t.Fatal("expected nodes")
	}
	for _, n := range nodes {
		t.Logf("  %s %s", n.CstxType(), n.CstxID())
	}
}

func TestVersion(t *testing.T) {
	v := Version()
	if v == "" || v == "stub" {
		t.Error("expected non-empty version")
	}
	t.Logf("version: %s", v)
}

func TestSupportedArtifacts(t *testing.T) {
	arts := SupportedArtifacts()
	if len(arts) == 0 {
		t.Fatal("expected artifacts")
	}
	t.Logf("artifacts: %v", arts)
}
