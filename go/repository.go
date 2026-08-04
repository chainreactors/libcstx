package cstx

import (
	"context"
	"encoding/json"
)

// Repository is the snapshot/version namespace of a CSTX runtime.
type Repository struct{ eng engine }

// IndexGraph derives canonical content hashes and raw CAS objects from the
// current graph using the native repository implementation.
func (r *Repository) IndexGraph(ctx context.Context) (CASIndex, error) {
	if err := contextError(ctx); err != nil {
		return CASIndex{}, err
	}
	raw, ok := r.eng.(rawEngine)
	if !ok {
		return CASIndex{}, &Error{Code: CodeNotInitialized, Operation: "repo.index_graph", Message: "native repository transport is unavailable"}
	}
	data, err := raw.rawRepositoryIndexGraph(ctx)
	if err != nil {
		return CASIndex{}, err
	}
	var index CASIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return CASIndex{}, &Error{Code: CodeParse, Operation: "repo.index_graph", Message: err.Error()}
	}
	return index, nil
}

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

// CommitDelta applies pre-indexed changes and removals to a repository ref.
func (r *Repository) CommitDelta(ctx context.Context, request CommitDeltaRequest) (Commit, error) {
	if err := contextError(ctx); err != nil {
		return Commit{}, err
	}
	if request.DeltaEntries == nil {
		request.DeltaEntries = map[string]map[string]string{}
	}
	request.RemovedNodeIDs = orEmpty(request.RemovedNodeIDs)
	request.RemovedEdgeIDs = orEmpty(request.RemovedEdgeIDs)
	return r.eng.repoCommitDelta(ctx, request)
}

// SetRef points a ref at an imported commit after native integrity validation.
func (r *Repository) SetRef(ctx context.Context, refName, commitHash string) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return r.eng.repoSetRef(ctx, refName, commitHash)
}

// ImportCommit loads one canonical commit object under its expected hash.
func (r *Repository) ImportCommit(ctx context.Context, hash string, data json.RawMessage) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return r.eng.repoImportCommit(ctx, hash, data)
}

// ExportCommit returns one canonical commit object, or JSON null when absent.
func (r *Repository) ExportCommit(ctx context.Context, hash string) (json.RawMessage, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoExportCommit(ctx, hash)
}

// ExportTreeObjects returns all Merkle objects reachable from a tree root.
func (r *Repository) ExportTreeObjects(ctx context.Context, rootHash string) (json.RawMessage, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoExportTreeObjects(ctx, rootHash)
}

// ImportTreeObjects loads canonical Merkle objects into the repository.
func (r *Repository) ImportTreeObjects(ctx context.Context, objects json.RawMessage) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return r.eng.repoImportTreeObjects(ctx, objects)
}

// TreeEntries returns content hashes grouped by element type and ID.
func (r *Repository) TreeEntries(ctx context.Context, rootHash string, types []string) (map[string]map[string]string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	if len(types) == 0 {
		stats, err := r.eng.repoTreeRootStats(ctx, rootHash)
		if err != nil {
			return nil, err
		}
		for nodeType := range stats {
			types = append(types, nodeType)
		}
	}
	return r.eng.repoTreeEntries(ctx, rootHash, types)
}

// FindTreeEntry returns one element content hash, or nil when absent.
func (r *Repository) FindTreeEntry(ctx context.Context, rootHash, elementID string) (*string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoFindTreeEntry(ctx, rootHash, elementID)
}

// DiffTreeEntries returns raw base/head content hashes for changed elements.
func (r *Repository) DiffTreeEntries(ctx context.Context, baseRoot, headRoot string) (TreeEntryDiff, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoDiffTreeEntries(ctx, baseRoot, headRoot)
}

// TreeRootStats returns element counts by type for a loaded tree root.
func (r *Repository) TreeRootStats(ctx context.Context, rootHash string) (map[string]int64, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return r.eng.repoTreeRootStats(ctx, rootHash)
}
