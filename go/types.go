package cstx

import "github.com/chainreactors/libcstx/go/proto/cstxproto"

// NodeFlags are engine-compatible scalar masks used by APIs that accept flag
// masks directly. Structured graph values use cstxproto.NodeFlag.
const (
	FlagNone               uint64 = 0
	FlagHoneypot           uint64 = 1 << 0
	FlagNoise              uint64 = 1 << 1
	FlagFalsePositive      uint64 = 1 << 2
	FlagManualIgnored      uint64 = 1 << 3
	FlagThreatPresent      uint64 = 1 << 4
	FlagHistoricVulnerable uint64 = 1 << 5
	FlagInternal           uint64 = 1 << 6
)

const FlagsAllMask uint64 = FlagHoneypot | FlagNoise | FlagFalsePositive |
	FlagManualIgnored | FlagThreatPresent | FlagHistoricVulnerable | FlagInternal

const FlagsDefaultExcludeMask uint64 = FlagHoneypot | FlagNoise | FlagFalsePositive | FlagManualIgnored

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
