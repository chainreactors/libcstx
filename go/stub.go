//go:build !cstx_native

package cstx

import "fmt"

var errStub = fmt.Errorf("cstx: native FFI not available (build with -tags cstx_native, run: go generate ./...)")

type Graph struct{}

func New() *Graph                                            { return nil }
func (g *Graph) Close()                                     {}
func (g *Graph) LoadPlugin(string) error                    { return errStub }
func (g *Graph) LoadAllPlugins() int                        { return 0 }
func AvailablePlugins() []string                            { return nil }
func (g *Graph) RegisterSchema(string, string, string) error { return errStub }
func (g *Graph) AddJoinRule(string, string, string, string, string, bool) error { return errStub }
func (g *Graph) AddNode(string, string, string, []string, uint64) error { return errStub }
func (g *Graph) AddEdge(string, string, string, string) error { return errStub }
func (g *Graph) AddNodesBatch(string) (*BatchResult, error) { return nil, errStub }
func (g *Graph) AddEdgesBatch(string) (int, error)          { return 0, errStub }
func (g *Graph) Link(string, []byte) (*IngestResult, error) { return nil, errStub }
func (g *Graph) IngestNative(string, string, []byte) (*IngestResult, error) { return nil, errStub }
func (g *Graph) HasNativeArtifact(string) bool              { return false }
func (g *Graph) IngestJSONL(string, []byte, string, string) (*IngestResult, error) { return nil, errStub }
func (g *Graph) LinkNodes(string) (string, error)           { return "", errStub }
func (g *Graph) Node(string) (string, bool)                 { return "", false }
func (g *Graph) Nodes(string) string                        { return "[]" }
func (g *Graph) Edges(string) string                        { return "[]" }
func (g *Graph) NodeTypes() []string                        { return nil }
func (g *Graph) Neighbors(string) string                    { return "[]" }
func (g *Graph) ContainsNode(string) bool                   { return false }
func (g *Graph) NodeCount() int                             { return 0 }
func (g *Graph) EdgeCount() int                             { return 0 }
func (g *Graph) NodeIDs() []string                          { return nil }
func (g *Graph) Stats() string                              { return "{}" }
func (g *Graph) ToJSON() string                             { return "{}" }
func (g *Graph) AllNodesJSON() string                       { return "[]" }
func (g *Graph) NodePayload(string) (string, bool)          { return "", false }
func Transform(string, []byte) ([]byte, error)              { return nil, errStub }
func Parse(string, any) ([]SCONode, error)                  { return nil, errStub }
func Version() string                                       { return "stub" }
func SupportedArtifacts() []string                          { return nil }

// ── Traversal stubs ──

func (g *Graph) NeighborIDs(string, string) []string                     { return nil }
func (g *Graph) BFS(string, uint32, bool) []string                       { return nil }
func (g *Graph) ShortestPaths(string, string, uint32) [][]string         { return nil }
func (g *Graph) Degree(string, string) int                               { return 0 }
func (g *Graph) SubgraphNodeIDs([]string, uint32) ([]string, []string, error) { return nil, nil, errStub }

// ── Query DSL stubs ──

func (g *Graph) QueryDSL(string, int, int) (string, error)  { return "", errStub }
func (g *Graph) QueryNodeIDs(string, int, int) ([]string, error) { return nil, errStub }
func (g *Graph) IsPathExpression(string) bool                { return false }

// ── Filter / Export stubs ──

func (g *Graph) EdgesFiltered(string, string, string) string { return "[]" }
func (g *Graph) ExportSnapshotJSON() string                  { return "{}" }
func (g *Graph) UpdateNodeFlags(string, uint64, uint64, int64) bool { return false }
func FlagsAllMask() uint64                                   { return 0 }
func FlagsDefaultExcludeMask() uint64                        { return 0 }
