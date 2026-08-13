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

func (r *Repository) Diff(
	ctx context.Context,
	base string,
	head string,
	limit *int,
) (GraphDiff, error) {
	if err := contextError(ctx); err != nil {
		return GraphDiff{}, err
	}
	return r.eng.repoDiff(ctx, base, head, limit)
}

func (r *Repository) DiffStat(ctx context.Context, base, head string) (Delta, error) {
	if err := contextError(ctx); err != nil {
		return Delta{}, err
	}
	return r.eng.repoDiffStat(ctx, base, head)
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
