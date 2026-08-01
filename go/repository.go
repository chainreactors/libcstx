package cstx

import "context"

// Repository is the snapshot/version namespace of a CSTX runtime.
type Repository struct{ eng engine }

// Commit writes the current graph content to a named repository ref.
func (r *Repository) Commit(ctx context.Context, refName, message string, metadata any) (Commit, error) {
	if err := contextError(ctx); err != nil {
		return Commit{}, err
	}
	return r.eng.repoCommit(ctx, refName, message, metadata)
}

// Diff returns added, removed, and modified IDs between two refs.
func (r *Repository) Diff(ctx context.Context, baseRef, headRef string) (GraphDiff, error) {
	if err := contextError(ctx); err != nil {
		return GraphDiff{}, err
	}
	return r.eng.repoDiff(ctx, baseRef, headRef)
}

// Dump returns a compact snapshot byte container. Compression is for example
// "zstd1"; an empty string selects the runtime default.
func (r *Repository) Dump(ctx context.Context, compression string) ([]byte, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoDump(ctx, compression)
}

// Load replaces graph state from a compact snapshot and reports the number
// of bytes consumed.
func (r *Repository) Load(ctx context.Context, data []byte, compression string) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	return r.eng.repoLoad(ctx, data, compression)
}

// DumpJSON returns canonical JSON snapshot bytes for interoperability.
func (r *Repository) DumpJSON(ctx context.Context) ([]byte, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoDumpJSON(ctx)
}

// LoadJSON loads canonical JSON snapshot bytes and reports bytes consumed.
func (r *Repository) LoadJSON(ctx context.Context, data []byte) (uint64, error) {
	if err := contextError(ctx); err != nil {
		return 0, err
	}
	return r.eng.repoLoadJSON(ctx, data)
}

// SnapshotFingerprint returns a stable hash of the canonical snapshot.
func (r *Repository) SnapshotFingerprint(ctx context.Context) (string, error) {
	if err := contextError(ctx); err != nil {
		return "", err
	}
	return r.eng.repoSnapshotFingerprint(ctx)
}

// Head returns the commit ID of one named ref, or nil when the ref does not
// exist.
func (r *Repository) Head(ctx context.Context, refName string) (*string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoHead(ctx, refName)
}

// Refs lists repository refs without exposing Merkle internals.
func (r *Repository) Refs(ctx context.Context) ([]Ref, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoRefs(ctx)
}
