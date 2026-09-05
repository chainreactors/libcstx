package cstx

import (
	"fmt"
	"testing"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
	"google.golang.org/protobuf/types/known/structpb"
)

// A per-entity history is cheap because the index pages in the postings for that
// one entity, not every object the range's snapshots contain. That property is
// only reachable from Go once the history plan exists: a host that keeps objects
// outside the runtime has no other way to learn which index pages to hand over,
// and would have to fall back to materializing each snapshot and comparing
// content hashes — the very cost the index exists to avoid.
//
// This exercises the whole external-storage round trip and then measures both
// plans against the same stored objects.
func TestRepositoryHistoryHydratesOnlyEntityPostings(t *testing.T) {
	const rounds = 12
	const width = 20
	const tracked = "domain:tracked.example"

	writer := openRuntime(t)
	objects := map[string]*cstxproto.RepositoryState_Object{}
	var head, indexRoot string
	var commits []string
	var commitObject *cstxproto.RepositoryState_Object

	for round := range rounds {
		// The tracked node changes every round...
		trackedNode := domainNode("tracked.example")
		trackedNode.Id = stringPtr(tracked)
		trackedNode.Annotations = &structpb.Struct{Fields: map[string]*structpb.Value{"round": structpb.NewNumberValue(float64(round))}}
		if _, err := writer.Graph.AddNodes(testContext, []*cstxproto.Node{trackedNode}); err != nil {
			t.Fatalf("round %d tracked node: %v", round, err)
		}
		// ...surrounded by nodes that do not, so the snapshot is wide while the
		// entity's own history stays short.
		filler := make([]*cstxproto.Node, 0, width)
		for i := range width {
			filler = append(filler, domainNode(fmt.Sprintf("filler-%d-%d.example", round, i)))
		}
		if _, err := writer.Graph.AddNodes(testContext, filler); err != nil {
			t.Fatalf("round %d filler: %v", round, err)
		}

		var expected *string
		if head != "" {
			expected = &head
		}
		prepared, err := writer.Repo.Prepare(
			testContext, fmt.Sprintf("round %d", round), "main", expected, nil, nil,
		)
		if err != nil {
			t.Fatalf("prepare round %d: %v", round, err)
		}
		for _, object := range prepared.Objects {
			stored := &cstxproto.RepositoryState_Object{Id: object.Id, Payload: append([]byte(nil), object.Payload...)}
			objects[object.Id] = stored
			if object.Kind == cstxproto.RepositoryObjectKind_REPOSITORY_OBJECT_KIND_COMMIT && object.Id == prepared.Commit.Id {
				commitObject = stored
			}
		}
		if err := writer.Repo.Accept(testContext, prepared.Commit.Id); err != nil {
			t.Fatalf("accept round %d: %v", round, err)
		}
		head = prepared.Commit.Id
		commits = append(commits, head)
		indexRoot = prepared.IndexRoot
	}

	rootObject, ok := objects[indexRoot]
	if !ok {
		t.Fatalf("index root %s was never published as an object", indexRoot)
	}

	// seed gives a fresh runtime only the commit frontier: it owns no tree and no
	// postings, so everything a plan needs has to arrive through synchronize.
	seed := func() *CSTX {
		reader := openRuntime(t)
		if err := reader.Repo.Synchronize(testContext, &cstxproto.RepositoryState{
			Objects: []*cstxproto.RepositoryState_Object{commitObject, rootObject},
		}); err != nil {
			t.Fatalf("synchronize frontier objects: %v", err)
		}
		if err := reader.Repo.Synchronize(testContext, &cstxproto.RepositoryState{
			Refs:    []*cstxproto.RepositoryState_Ref{{Name: "main", CommitId: &head}},
			Indexes: []*cstxproto.RepositoryState_Index{{CommitId: head, IndexRoot: indexRoot}},
		}); err != nil {
			t.Fatalf("synchronize frontier refs: %v", err)
		}
		return reader
	}

	hydrate := func(reader *CSTX, plan func() (*cstxproto.ObjectSelection, error)) int {
		read := 0
		for {
			missing, err := plan()
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
					t.Fatalf("planner requested an object that was never stored: %s", id)
				}
				batch = append(batch, object)
			}
			read += len(batch)
			if err := reader.Repo.Synchronize(testContext, &cstxproto.RepositoryState{Objects: batch}); err != nil {
				t.Fatalf("synchronize: %v", err)
			}
		}
	}

	historyReader := seed()
	historyObjects := hydrate(historyReader, func() (*cstxproto.ObjectSelection, error) {
		return historyReader.Repo.Missing(testContext, &cstxproto.RepositoryObjectPlan{Kind: cstxproto.RepositoryPlanKind_REPOSITORY_PLAN_HISTORY, CommitId: head, EntityId: stringPtr(tracked)})
	})
	entries, err := historyReader.Repo.History(testContext, tracked, head, nil)
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	if len(entries.Changes) != rounds {
		t.Fatalf("history returned %d entries, want %d (one per round)", len(entries.Changes), rounds)
	}

	// The fallback without an indexed history plan is to materialize every
	// snapshot in the range and compare the entity's content hash across them.
	// One snapshot is cheap; the range is not, and it grows with history depth
	// while the entity's own change count does not.
	walkObjects := 0
	for _, commit := range commits {
		reader := openRuntime(t)
		commitEnvelope, ok := objects[commit]
		if !ok {
			t.Fatalf("commit object %s was never published", commit)
		}
		if err := reader.Repo.Synchronize(testContext, &cstxproto.RepositoryState{
			Objects: []*cstxproto.RepositoryState_Object{commitEnvelope},
		}); err != nil {
			t.Fatalf("synchronize commit %s: %v", commit, err)
		}
		walkObjects += 1 + hydrate(reader, func() (*cstxproto.ObjectSelection, error) {
			return reader.Repo.Missing(testContext, &cstxproto.RepositoryObjectPlan{Kind: cstxproto.RepositoryPlanKind_REPOSITORY_PLAN_TREE, CommitId: commit})
		})
	}

	if historyObjects >= walkObjects {
		t.Fatalf(
			"per-entity history hydrated %d objects and the snapshot walk hydrated %d "+
				"across %d commits; the index exists so the walk is not needed",
			historyObjects, walkObjects, len(commits),
		)
	}
	t.Logf(
		"per-entity history: %d objects for %d changes; snapshot walk over %d commits: %d objects",
		historyObjects, len(entries.Changes), len(commits), walkObjects,
	)
}
