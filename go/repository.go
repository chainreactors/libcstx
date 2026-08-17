package cstx

import "context"

// Repository is the Git-like version namespace of one CSTX working tree.
type Repository struct{ eng engine }

func (r *Repository) Resolve(ctx context.Context, revision string) (string, error) {
	if err := contextError(ctx); err != nil {
		return "", err
	}
	return r.eng.repoResolve(ctx, revision)
}

func (r *Repository) Head(ctx context.Context, refName string) (*string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoHead(ctx, refName)
}

func (r *Repository) Checkout(ctx context.Context, revision string, force bool) (Commit, error) {
	if err := contextError(ctx); err != nil {
		return Commit{}, err
	}
	return r.eng.repoCheckout(ctx, revision, force)
}

func (r *Repository) Commit(
	ctx context.Context,
	message string,
	refName string,
	expectedHead *string,
	metadata any,
) (Commit, error) {
	if err := contextError(ctx); err != nil {
		return Commit{}, err
	}
	return r.eng.repoCommit(ctx, message, refName, expectedHead, metadata)
}

// Prepare computes one commit without publishing its ref. The caller must
// durably persist the returned payload and atomically advance the external ref,
// then call Accept. Call Discard when the external transaction fails.
func (r *Repository) Prepare(
	ctx context.Context,
	message string,
	refName string,
	expectedHead *string,
	metadata any,
	timestamp *int64,
) (PreparedCommit, error) {
	if err := contextError(ctx); err != nil {
		return PreparedCommit{}, err
	}
	return r.eng.repoPrepare(ctx, message, refName, expectedHead, metadata, timestamp)
}

// Accept finalizes a prepared commit after its external transaction commits.
func (r *Repository) Accept(ctx context.Context, commit string) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return r.eng.repoAccept(ctx, commit)
}

// Discard abandons a prepared commit while preserving the working journal.
func (r *Repository) Discard(ctx context.Context) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return r.eng.repoDiscard(ctx)
}

// Synchronize loads externally persisted objects, refs, and index roots into
// this computation session.
func (r *Repository) Synchronize(ctx context.Context, state RepositorySync) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return r.eng.repoSynchronize(ctx, state)
}

func (r *Repository) Contains(ctx context.Context, object string) (bool, error) {
	if err := contextError(ctx); err != nil {
		return false, err
	}
	return r.eng.repoContains(ctx, object)
}

// MissingTree plans immutable object reads required to materialize a commit.
func (r *Repository) MissingTree(ctx context.Context, commit string) ([]string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoMissingTree(ctx, commit)
}

// ObjectClosure returns every object one commit and its ancestry are built
// from. Deleting whatever the union of this set over every ref does not name
// reclaims space without breaking any supported operation on those refs.
//
// It answers from stored bytes, so unlike the Missing* planners the result does
// not depend on what this process has already loaded.
func (r *Repository) ObjectClosure(ctx context.Context, commit string) ([]string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoObjectClosure(ctx, commit)
}

// MissingPrepare plans index reads required before preparing a child commit.
func (r *Repository) MissingPrepare(ctx context.Context, commit string) ([]string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoMissingPrepare(ctx, commit)
}

// MissingHistory plans the index reads required to answer History for one
// entity. A host that keeps objects outside the runtime resolves this to empty
// before calling History; the index only pages in the postings for that entity,
// so the walk costs what the entity changed, not what the range contains.
func (r *Repository) MissingHistory(ctx context.Context, commit, entity string) ([]string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoMissingHistory(ctx, commit, entity)
}

// MissingStat plans the reads required to summarize a commit.
func (r *Repository) MissingStat(ctx context.Context, commit string) ([]string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoMissingStat(ctx, commit)
}

// MissingCommits plans the reads required to walk a commit's ancestry.
func (r *Repository) MissingCommits(ctx context.Context, commit string, limit int) ([]string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoMissingCommits(ctx, commit, limit)
}

// MissingDiff plans the reads required to diff two revisions at one detail
// level. A limit never narrows the plan, so it is not part of the request.
func (r *Repository) MissingDiff(ctx context.Context, base, head string, detail DiffDetail) ([]string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if detail == "" {
		detail = DiffEntities
	}
	return r.eng.repoMissingDiff(ctx, base, head, detail)
}

// MissingDelta plans the reads required to count changes in a time range.
func (r *Repository) MissingDelta(ctx context.Context, commit string, start, end *int64) ([]string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoMissingDelta(ctx, commit, start, end)
}

// MissingMerge plans the reads required to merge source into target. An empty
// target means the current head.
func (r *Repository) MissingMerge(ctx context.Context, source, target string) ([]string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoMissingMerge(ctx, source, target)
}

// ReleaseTransientObjects drops objects hydrated for one external operation.
func (r *Repository) ReleaseTransientObjects(ctx context.Context) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return r.eng.repoReleaseTransientObjects(ctx)
}

// Diff compares two revisions. Options select a limit on the reported entity
// IDs and whether they are reported at all; the counts in GraphDiff.Stats are
// exact either way.
func (r *Repository) Diff(
	ctx context.Context,
	base string,
	head string,
	options DiffOptions,
) (GraphDiff, error) {
	if err := contextError(ctx); err != nil {
		return GraphDiff{}, err
	}
	return r.eng.repoDiff(ctx, base, head, options)
}

func (r *Repository) Log(
	ctx context.Context,
	revision string,
	limit int,
) ([]map[string]any, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoLog(ctx, revision, limit)
}

// History is a structured, replayable result for one entity at a revision.
type History struct {
	EntityID string
	Revision string
	Limit    *int
	Entries  []map[string]any
}

func (r *Repository) History(
	ctx context.Context,
	entityID string,
	revision string,
	limit *int,
) (History, error) {
	if err := contextError(ctx); err != nil {
		return History{}, err
	}
	entries, err := r.eng.repoHistory(ctx, entityID, revision, limit)
	return History{
		EntityID: entityID,
		Revision: revision,
		Limit:    limit,
		Entries:  entries,
	}, err
}

func (r *Repository) Branch(ctx context.Context, name, startPoint string) (string, error) {
	if err := contextError(ctx); err != nil {
		return "", err
	}
	return r.eng.repoBranch(ctx, name, startPoint)
}

func (r *Repository) Merge(
	ctx context.Context,
	source string,
	target string,
	expectedHead *string,
	message *string,
) (Commit, error) {
	if err := contextError(ctx); err != nil {
		return Commit{}, err
	}
	return r.eng.repoMerge(ctx, source, target, expectedHead, message)
}

func (r *Repository) Stat(
	ctx context.Context,
	revision string,
	excludeMask uint64,
	includeMask uint64,
) (GraphStats, error) {
	if err := contextError(ctx); err != nil {
		return GraphStats{}, err
	}
	return r.eng.repoStat(ctx, revision, excludeMask, includeMask)
}

func (r *Repository) Delta(
	ctx context.Context,
	revision string,
	startTimestamp *int64,
	endTimestamp *int64,
) (Delta, error) {
	if err := contextError(ctx); err != nil {
		return Delta{}, err
	}
	return r.eng.repoDelta(ctx, revision, startTimestamp, endTimestamp)
}
