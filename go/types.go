package cstx

import "github.com/chainreactors/libcstx/go/proto/cstxproto"

// FlagNone is the empty mask. Every other flag bit is *declared by an
// extension*, not by this SDK: a bit's meaning lives in
// `<extension>.schema.json`, and hard-coding one product's seven security
// words here was the same violation the `NodeFlag` enum was on the wire.
// Read the declarations with FlagRegistry (flags.go).
const FlagNone uint64 = 0

// Affected returns the number of graph entities changed by a generated
// protobuf change set. It is a function instead of a shadow SDK struct method.
func Affected(change *cstxproto.GraphChangeSet) int {
	if change == nil {
		return 0
	}
	return len(change.AddedNodeIds) + len(change.UpdatedNodeIds) +
		len(change.RemovedNodeIds) + len(change.AddedRelationshipIds) +
		len(change.UpdatedRelationshipIds) + len(change.RemovedRelationshipIds)
}

func algorithmCursorKind(algorithm *cstxproto.Algorithm) CursorKind {
	if algorithm == nil {
		return CursorKindNodes
	}
	switch kind := algorithm.Kind.(type) {
	case *cstxproto.Algorithm_Bfs:
		return CursorKindNodes
	case *cstxproto.Algorithm_Betweenness, *cstxproto.Algorithm_Closeness:
		return CursorKindNodeScores
	case *cstxproto.Algorithm_Leiden:
		return CursorKindCommunities
	case *cstxproto.Algorithm_ShortestPaths:
		return CursorKindPaths
	case *cstxproto.Algorithm_Parameterless:
		switch kind.Parameterless {
		case cstxproto.ParameterlessAlgorithm_PARAMETERLESS_WEAK_COMPONENTS,
			cstxproto.ParameterlessAlgorithm_PARAMETERLESS_STRONG_COMPONENTS:
			return CursorKindComponents
		case cstxproto.ParameterlessAlgorithm_PARAMETERLESS_CYCLE_BASIS:
			return CursorKindCycles
		case cstxproto.ParameterlessAlgorithm_PARAMETERLESS_BRIDGES:
			return CursorKindNodePairs
		case cstxproto.ParameterlessAlgorithm_PARAMETERLESS_CORE_NUMBERS:
			return CursorKindNodeScores
		default:
			return CursorKindNodes
		}
	default:
		return CursorKindNodes
	}
}
