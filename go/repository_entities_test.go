package cstx

import (
	"fmt"
	"testing"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

// Reading a few entities at a revision must not cost the revision. Before the
// entity plan existed the only way to get a payload back was Checkout, so a
// host that wanted the twenty nodes a diff named had to materialize the whole
// board — twice, for a two-sided diff.
//
// This runs the real external-storage loop (Missing -> fetch -> Synchronize)
// for both routes against the same stored objects and compares what each one
// had to fetch, then checks the point read agrees with the checkout field for
// field.
func TestRepositoryEntitiesHydratesLessThanACheckout(t *testing.T) {
	const width = 400
	const wanted = "domain:tracked.example"

	writer := openRuntime(t)
	nodes := make([]*cstxproto.Node, 0, width+1)
	tracked := domainNode("tracked.example")
	tracked.Id = stringPtr(wanted)
	nodes = append(nodes, tracked)
	for i := range width {
		nodes = append(nodes, domainNode(fmt.Sprintf("filler-%d.example", i)))
	}
	if _, err := writer.Graph.AddNodes(testContext, nodes); err != nil {
		t.Fatalf("add nodes: %v", err)
	}

	prepared, err := writer.Repo.Prepare(testContext, "baseline", "main", nil, nil, nil)
	if err != nil {
		t.Fatalf("prepare: %v", err)
	}
	objects := map[string]*cstxproto.RepositoryState_Object{}
	var commitObject *cstxproto.RepositoryState_Object
	for _, object := range prepared.Objects {
		stored := &cstxproto.RepositoryState_Object{
			Id:      object.Id,
			Payload: append([]byte(nil), object.Payload...),
		}
		objects[object.Id] = stored
		if object.Id == prepared.Commit.Id {
			commitObject = stored
		}
	}
	if err := writer.Repo.Accept(testContext, prepared.Commit.Id); err != nil {
		t.Fatalf("accept: %v", err)
	}
	head := prepared.Commit.Id

	// A fresh runtime holding only the commit envelope: everything a plan needs
	// has to arrive through Synchronize, so what it asks for is observable.
	seed := func() *CSTX {
		reader := openRuntime(t)
		if err := reader.Repo.Synchronize(testContext, &cstxproto.RepositoryState{
			Objects: []*cstxproto.RepositoryState_Object{commitObject},
			Refs:    []*cstxproto.RepositoryState_Ref{{Name: "main", CommitId: &head}},
		}); err != nil {
			t.Fatalf("synchronize frontier: %v", err)
		}
		return reader
	}

	hydrate := func(reader *CSTX, plan *cstxproto.RepositoryObjectPlan) int {
		read := 0
		for {
			missing, err := reader.Repo.Missing(testContext, plan)
			if err != nil {
				t.Fatalf("plan: %v", err)
			}
			if len(missing.ObjectIds) == 0 {
				return read
			}
			batch := make([]*cstxproto.RepositoryState_Object, 0, len(missing.ObjectIds))
			for _, id := range missing.ObjectIds {
				object, ok := objects[id]
				if !ok {
					t.Fatalf("planner asked for an object that was never stored: %s", id)
				}
				batch = append(batch, object)
			}
			read += len(batch)
			if err := reader.Repo.Synchronize(testContext, &cstxproto.RepositoryState{Objects: batch}); err != nil {
				t.Fatalf("synchronize: %v", err)
			}
		}
	}

	pointReader := seed()
	pointObjects := hydrate(pointReader, &cstxproto.RepositoryObjectPlan{
		Kind:      cstxproto.RepositoryPlanKind_REPOSITORY_PLAN_ENTITIES,
		CommitId:  head,
		EntityIds: []string{wanted},
	})
	read, err := pointReader.Repo.Entities(testContext, head, []string{wanted})
	if err != nil {
		t.Fatalf("entities: %v", err)
	}
	if len(read.Nodes) != 1 {
		t.Fatalf("entities returned %d nodes, want 1", len(read.Nodes))
	}

	checkoutReader := seed()
	checkoutObjects := hydrate(checkoutReader, &cstxproto.RepositoryObjectPlan{
		Kind:     cstxproto.RepositoryPlanKind_REPOSITORY_PLAN_TREE,
		CommitId: head,
	})
	if _, err := checkoutReader.Repo.Checkout(testContext, head, true); err != nil {
		t.Fatalf("checkout: %v", err)
	}
	resolved, err := checkoutReader.Graph.Node(testContext, wanted)
	if err != nil {
		t.Fatalf("graph node: %v", err)
	}

	if read.Nodes[0].String() != resolved.String() {
		t.Fatalf("point read %v disagrees with checkout %v", read.Nodes[0], resolved)
	}
	if pointObjects >= checkoutObjects {
		t.Fatalf(
			"entity read hydrated %d objects and the checkout hydrated %d over %d nodes; "+
				"the point read exists so the board does not have to be fetched",
			pointObjects, checkoutObjects, width+1,
		)
	}
	t.Logf("entity read: %d objects; checkout of the same revision: %d", pointObjects, checkoutObjects)
}
