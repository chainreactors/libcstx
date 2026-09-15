package cstx

import (
	"errors"
	"slices"
	"testing"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
	"google.golang.org/protobuf/types/known/structpb"
)

// These cover the operations the Go SDK declared through cstx_ffi.h but never
// bound, so each one asserts the Rust semantics rather than only that the call
// returns.

func TestFindNodeResolvesIdentifier(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "findable.example")

	node, err := rt.Graph.FindNode(testContext, "findable.example")
	if err != nil {
		t.Fatalf("find node: %v", err)
	}
	if got := domainValue(t, node); got != "findable.example" {
		t.Fatalf("found %q; want findable.example", got)
	}

	if _, err := rt.Graph.FindNode(testContext, "absent.example"); err == nil {
		t.Fatal("find node resolved an identifier that is not in the graph")
	}
}

func TestNodeTypesReportsPresentTypesOnly(t *testing.T) {
	rt := openRuntime(t)

	empty, err := rt.Graph.NodeTypes(testContext)
	if err != nil {
		t.Fatalf("node types: %v", err)
	}
	// easm is enabled, so the schema declares many types. None is present yet.
	if len(empty) != 0 {
		t.Fatalf("empty graph reported node types %v", empty)
	}

	addDomain(t, rt, "typed.example")
	types, err := rt.Graph.NodeTypes(testContext)
	if err != nil {
		t.Fatalf("node types: %v", err)
	}
	if !slices.Contains(types, "domain") {
		t.Fatalf("node types = %v; want it to contain domain", types)
	}
}

func TestDegreeCountsByDirection(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "source.example")
	addDomain(t, rt, "target.example")
	if _, err := rt.Graph.AddRelationships(testContext, []*cstxproto.Relationship{
		usesRelationship("domain:source.example", "domain:target.example"),
	}); err != nil {
		t.Fatalf("add relationship: %v", err)
	}

	for _, testCase := range []struct {
		direction string
		want      uint64
	}{
		{"out", 1},
		{"in", 0},
		{"both", 1},
		// An empty direction is the documented "both".
		{"", 1},
	} {
		got, err := rt.Graph.Degree(testContext, "domain:source.example", testCase.direction)
		if err != nil {
			t.Fatalf("degree %q: %v", testCase.direction, err)
		}
		if got != testCase.want {
			t.Fatalf("degree %q = %d; want %d", testCase.direction, got, testCase.want)
		}
	}
}

func TestUpdateNodeFlagsMergesAndReplaces(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "flagged.example")
	selection := &cstxproto.GraphSelection{NodeIds: []string{"domain:flagged.example"}}

	affected, err := rt.Graph.UpdateNodeFlags(testContext, &cstxproto.NodeFlagChange{
		Selection: selection,
		Update: &cstxproto.NodeFlagUpdate{
			Mode:    cstxproto.NodeFlagUpdateMode_NODE_FLAG_UPDATE_MERGE,
			AddMask: 0b101,
		},
	})
	if err != nil {
		t.Fatalf("merge flags: %v", err)
	}
	if affected != 1 {
		t.Fatalf("merge affected = %d; want 1", affected)
	}
	node, err := rt.Graph.Node(testContext, "domain:flagged.example")
	if err != nil {
		t.Fatalf("node: %v", err)
	}
	if node.GetFlagsMask() != 0b101 {
		t.Fatalf("flags after merge = %b; want 101", node.GetFlagsMask())
	}

	// Replace overwrites the mask instead of OR-ing into it, which is the
	// distinction the mode field exists for.
	if _, err := rt.Graph.UpdateNodeFlags(testContext, &cstxproto.NodeFlagChange{
		Selection: selection,
		Update: &cstxproto.NodeFlagUpdate{
			Mode:        cstxproto.NodeFlagUpdateMode_NODE_FLAG_UPDATE_REPLACE,
			ReplaceMask: 0b010,
		},
	}); err != nil {
		t.Fatalf("replace flags: %v", err)
	}
	node, err = rt.Graph.Node(testContext, "domain:flagged.example")
	if err != nil {
		t.Fatalf("node: %v", err)
	}
	if node.GetFlagsMask() != 0b010 {
		t.Fatalf("flags after replace = %b; want 010", node.GetFlagsMask())
	}
}

func TestPatchNodeAnnotationsMergesIntoSelection(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "annotated.example")

	affected, err := rt.Graph.PatchNodeAnnotations(testContext, &cstxproto.NodeAnnotationUpdate{
		Selection:   &cstxproto.GraphSelection{NodeIds: []string{"domain:annotated.example"}},
		Annotations: &structpb.Struct{Fields: map[string]*structpb.Value{"owner": structpb.NewStringValue("go-test")}},
	})
	if err != nil {
		t.Fatalf("patch annotations: %v", err)
	}
	if affected != 1 {
		t.Fatalf("affected = %d; want 1", affected)
	}
	node, err := rt.Graph.Node(testContext, "domain:annotated.example")
	if err != nil {
		t.Fatalf("node: %v", err)
	}
	if got := node.GetAnnotations().GetFields()["owner"].GetStringValue(); got != "go-test" {
		t.Fatalf("annotation owner = %q; want go-test", got)
	}
}

func TestUnionAndDifferenceLeaveInputsAlone(t *testing.T) {
	left := openRuntime(t)
	addDomain(t, left, "shared.example")
	addDomain(t, left, "left-only.example")

	right := openRuntime(t)
	addDomain(t, right, "shared.example")

	union, err := left.Graph.Union(testContext, right.Graph)
	if err != nil {
		t.Fatalf("union: %v", err)
	}
	t.Cleanup(func() { _ = union.Close() })
	count, err := union.Graph.NodeCount(testContext)
	if err != nil {
		t.Fatalf("union node count: %v", err)
	}
	if count != 2 {
		t.Fatalf("union node count = %d; want 2", count)
	}

	difference, err := left.Graph.Difference(testContext, right.Graph, "")
	if err != nil {
		t.Fatalf("difference: %v", err)
	}
	t.Cleanup(func() { _ = difference.Close() })
	if _, err := difference.Graph.Node(testContext, "domain:left-only.example"); err != nil {
		t.Fatalf("difference is missing the left-only node: %v", err)
	}
	if contains, err := difference.Graph.Contains(testContext, "domain:shared.example"); err != nil {
		t.Fatalf("difference contains: %v", err)
	} else if contains {
		t.Fatal("difference kept a node present in both graphs")
	}

	// Neither input was modified by either derived graph.
	leftCount, err := left.Graph.NodeCount(testContext)
	if err != nil {
		t.Fatalf("left node count: %v", err)
	}
	if leftCount != 2 {
		t.Fatalf("left node count = %d; want 2", leftCount)
	}
}

func TestMergeWritesIntoTheTargetGraph(t *testing.T) {
	target := openRuntime(t)
	addDomain(t, target, "target-only.example")

	source := openRuntime(t)
	addDomain(t, source, "merged-in.example")

	affected, err := target.Graph.Merge(testContext, source.Graph)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	if affected == 0 {
		t.Fatal("merge reported no change")
	}
	if contains, err := target.Graph.Contains(testContext, "domain:merged-in.example"); err != nil {
		t.Fatalf("contains: %v", err)
	} else if !contains {
		t.Fatal("merge did not write the source node into the target")
	}
	// Unlike Union, the source is untouched and the target grew in place.
	sourceCount, err := source.Graph.NodeCount(testContext)
	if err != nil {
		t.Fatalf("source node count: %v", err)
	}
	if sourceCount != 1 {
		t.Fatalf("source node count = %d; want 1", sourceCount)
	}
}

func TestInducedSubgraphKeepsSelectedNodes(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "kept.example")
	addDomain(t, rt, "dropped.example")

	induced, err := rt.Graph.InducedSubgraph(testContext, []string{"domain:kept.example"}, nil)
	if err != nil {
		t.Fatalf("induced subgraph: %v", err)
	}
	t.Cleanup(func() { _ = induced.Close() })

	count, err := induced.Graph.NodeCount(testContext)
	if err != nil {
		t.Fatalf("induced node count: %v", err)
	}
	if count != 1 {
		t.Fatalf("induced node count = %d; want 1", count)
	}
	if contains, err := induced.Graph.Contains(testContext, "domain:dropped.example"); err != nil {
		t.Fatalf("contains: %v", err)
	} else if contains {
		t.Fatal("induced subgraph kept an unselected node")
	}
}

func TestFilterWithReasonsExplainsExclusions(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "included.example")
	addDomain(t, rt, "excluded.example")
	if _, err := rt.Graph.UpdateNodeFlags(testContext, &cstxproto.NodeFlagChange{
		Selection: &cstxproto.GraphSelection{NodeIds: []string{"domain:excluded.example"}},
		Update: &cstxproto.NodeFlagUpdate{
			Mode:    cstxproto.NodeFlagUpdateMode_NODE_FLAG_UPDATE_MERGE,
			AddMask: 0b1,
		},
	}); err != nil {
		t.Fatalf("flag node: %v", err)
	}

	filter := &cstxproto.NodeFilter{FlagsNoneMask: 0b1}
	filtered, err := rt.Graph.Filter(testContext, filter)
	if err != nil {
		t.Fatalf("filter: %v", err)
	}
	t.Cleanup(func() { _ = filtered.Close() })
	count, err := filtered.Graph.NodeCount(testContext)
	if err != nil {
		t.Fatalf("filtered node count: %v", err)
	}
	if count != 1 {
		t.Fatalf("filtered node count = %d; want 1", count)
	}

	// The reasons variant returns the same projection plus why each node went.
	explained, report, err := rt.Graph.FilterWithReasons(testContext, filter)
	if err != nil {
		t.Fatalf("filter with reasons: %v", err)
	}
	t.Cleanup(func() { _ = explained.Close() })
	if len(report.GetExcludedNodes()) != 1 {
		t.Fatalf("excluded %d node(s); want 1: %+v", len(report.GetExcludedNodes()), report)
	}
	excluded := report.GetExcludedNodes()[0]
	if excluded.GetNodeId() != "domain:excluded.example" {
		t.Fatalf("excluded node = %q; want domain:excluded.example", excluded.GetNodeId())
	}
	if excluded.GetReason() == "" {
		t.Fatal("exclusion carries no reason, which is the whole point of this variant")
	}
}

func TestQuerySubgraphMaterializesTracedPath(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "traced.example")

	subgraph, err := rt.Graph.QuerySubgraph(testContext, &cstxproto.GraphQuery{Expression: "domain"})
	if err != nil {
		t.Fatalf("query subgraph: %v", err)
	}
	t.Cleanup(func() { _ = subgraph.Close() })

	count, err := subgraph.Graph.NodeCount(testContext)
	if err != nil {
		t.Fatalf("subgraph node count: %v", err)
	}
	if count != 1 {
		t.Fatalf("subgraph node count = %d; want 1", count)
	}
}

func TestFindAnchorsAndElevateShareConceptNames(t *testing.T) {
	rt := openRuntime(t)
	addDomain(t, rt, "anchor.example")

	concepts, err := rt.Extensions.AnchorConcepts(testContext)
	if err != nil {
		t.Fatalf("anchor concepts: %v", err)
	}
	if len(concepts.GetConcepts()) == 0 {
		t.Skip("the enabled extensions declare no anchor concept")
	}
	concept := concepts.GetConcepts()[0].GetName()

	if _, err := rt.Graph.FindAnchors(testContext, concept); err != nil {
		t.Fatalf("find anchors %q: %v", concept, err)
	}
	elevated, err := rt.Graph.Elevate(testContext, concept)
	if err != nil {
		t.Fatalf("elevate %q: %v", concept, err)
	}
	t.Cleanup(func() { _ = elevated.Close() })

	// An unknown concept is an error, not an empty result.
	if _, err := rt.Graph.FindAnchors(testContext, "not-a-concept"); err == nil {
		t.Fatal("find anchors accepted an unknown concept")
	}
}

func TestBinaryGraphOperationsRejectNilPeer(t *testing.T) {
	rt := openRuntime(t)
	for name, call := range map[string]func() error{
		"union":      func() error { _, err := rt.Graph.Union(testContext, nil); return err },
		"difference": func() error { _, err := rt.Graph.Difference(testContext, nil, ""); return err },
		"merge":      func() error { _, err := rt.Graph.Merge(testContext, nil); return err },
	} {
		err := call()
		if err == nil {
			t.Fatalf("%s accepted a nil peer", name)
		}
		var cstxErr *Error
		if !errors.As(err, &cstxErr) || cstxErr.Code != CodeInvalidArgument {
			t.Fatalf("%s error = %v; want CodeInvalidArgument", name, err)
		}
	}
}
