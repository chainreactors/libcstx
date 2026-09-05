package cstx

import (
	"context"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
	"google.golang.org/protobuf/types/known/structpb"
)

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

func (r *Repository) Checkout(ctx context.Context, revision string, force bool) (*cstxproto.Commit, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoCheckout(ctx, revision, force)
}

func (r *Repository) Commit(
	ctx context.Context,
	message string,
	refName string,
	expectedHead *string,
	metadata *structpb.Struct,
) (*cstxproto.Commit, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
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
	metadata *structpb.Struct,
	timestamp *int64,
) (*cstxproto.PublicationPlan, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
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
func (r *Repository) Synchronize(ctx context.Context, state *cstxproto.RepositoryState) error {
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

// Missing plans immutable object reads for one repository operation.
func (r *Repository) Missing(ctx context.Context, plan *cstxproto.RepositoryObjectPlan) (*cstxproto.ObjectSelection, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoMissing(ctx, plan)
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
	limit *uint64,
	detail cstxproto.DiffDetail,
) (*cstxproto.GraphDiff, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoDiff(ctx, base, head, limit, detail)
}

func (r *Repository) Log(
	ctx context.Context,
	revision string,
	limit int,
) (*cstxproto.CommitLog, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoLog(ctx, revision, limit)
}

func (r *Repository) History(
	ctx context.Context,
	entityID string,
	revision string,
	limit *int,
) (*cstxproto.EntityHistory, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoHistory(ctx, entityID, revision, limit)
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
) (*cstxproto.Commit, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoMerge(ctx, source, target, expectedHead, message)
}

func (r *Repository) Stat(
	ctx context.Context,
	revision string,
	excludeMask uint64,
	includeMask uint64,
) (*cstxproto.GraphStats, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoStat(ctx, revision, excludeMask, includeMask)
}

func (r *Repository) Delta(
	ctx context.Context,
	revision string,
	startTimestamp *int64,
	endTimestamp *int64,
) (*cstxproto.GraphChangeSummary, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoDelta(ctx, revision, startTimestamp, endTimestamp)
}
