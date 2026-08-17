package cstx

import (
	"fmt"
	"testing"
)

// A per-entity history is cheap because the index pages in the postings for that
// one entity, not every object the range's snapshots contain. That property is
// only reachable from Go once MissingHistory exists: a host that keeps objects
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
	objects := map[string]RepositoryObject{}
	var head, indexRoot string
	var commits []string
	var commitObject RepositoryObject

	for round := range rounds {
		// The tracked node changes every round...
		if _, err := writer.Graph.AddNodes(testContext, []Node{{
			ID: tracked, Type: "domain", Value: "tracked.example",
			Model:   map[string]any{"domain": "tracked.example", "cstx_flags": 0},
			Sources: []string{"test"},
			Extras:  map[string]any{"round": round},
		}}); err != nil {
			t.Fatalf("round %d tracked node: %v", round, err)
		}
		// ...surrounded by nodes that do not, so the snapshot is wide while the
		// entity's own history stays short.
		filler := make([]Node, 0, width)
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
			stored := RepositoryObject{ID: object.ID, Envelope: append([]byte(nil), object.Envelope...)}
			objects[object.ID] = stored
			if object.Kind == "commit" && object.ID == prepared.Commit.ID {
				commitObject = stored
			}
		}
		if err := writer.Repo.Accept(testContext, prepared.Commit.ID); err != nil {
			t.Fatalf("accept round %d: %v", round, err)
		}
		head = prepared.Commit.ID
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
		if err := reader.Repo.Synchronize(testContext, RepositorySync{
			Objects: []RepositoryObject{commitObject, rootObject},
		}); err != nil {
			t.Fatalf("synchronize frontier objects: %v", err)
		}
		if err := reader.Repo.Synchronize(testContext, RepositorySync{
			Refs:    []RepositoryRef{{Name: "main", Commit: &head}},
			Indexes: []RepositoryIndex{{Commit: head, IndexRoot: indexRoot}},
		}); err != nil {
			t.Fatalf("synchronize frontier refs: %v", err)
		}
		return reader
	}

	hydrate := func(reader *CSTX, plan func() ([]string, error)) int {
		read := 0
		for {
			missing, err := plan()
			if err != nil {
				t.Fatalf("plan: %v", err)
			}
			if len(missing) == 0 {
				return read
			}
			batch := make([]RepositoryObject, 0, len(missing))
			for _, id := range missing {
				object, ok := objects[id]
				if !ok {
					t.Fatalf("planner requested an object that was never stored: %s", id)
				}
				batch = append(batch, object)
			}
			read += len(batch)
			if err := reader.Repo.Synchronize(testContext, RepositorySync{Objects: batch}); err != nil {
				t.Fatalf("synchronize: %v", err)
			}
		}
	}

	historyReader := seed()
	historyObjects := hydrate(historyReader, func() ([]string, error) {
		return historyReader.Repo.MissingHistory(testContext, head, tracked)
	})
	entries, err := historyReader.Repo.History(testContext, tracked, head, nil)
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	if len(entries.Entries) != rounds {
		t.Fatalf("history returned %d entries, want %d (one per round)", len(entries.Entries), rounds)
	}

	// The fallback a host without MissingHistory is stuck with: materialize every
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
		if err := reader.Repo.Synchronize(testContext, RepositorySync{
			Objects: []RepositoryObject{commitEnvelope},
		}); err != nil {
			t.Fatalf("synchronize commit %s: %v", commit, err)
		}
		walkObjects += 1 + hydrate(reader, func() ([]string, error) {
			return reader.Repo.MissingTree(testContext, commit)
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
		historyObjects, len(entries.Entries), len(commits), walkObjects,
	)
}
