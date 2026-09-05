package cstx

import (
	"context"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

// ExtensionBuilder constructs one canonical protobuf extension contract. It
// is a convenience for code-defined extensions, not a second wire format.
type ExtensionBuilder struct {
	contract *cstxproto.ExtensionContract
	name     string
}

func newExtensionBuilder(name, version string) *ExtensionBuilder {
	contract := &cstxproto.ExtensionContract{
		ContractVersion: 1,
		Extensions:      map[string]*cstxproto.ExtensionDefinition{},
	}
	contract.Extensions[name] = &cstxproto.ExtensionDefinition{
		Name:    name,
		Version: version,
		Parsers: map[string]*cstxproto.ParserType{},
	}
	return &ExtensionBuilder{contract: contract, name: name}
}

func (b *ExtensionBuilder) definition() *cstxproto.ExtensionDefinition {
	return b.contract.Extensions[b.name]
}

// Schema sets this extension's runtime schema document.
//
// One JSON document declares every node and relation type the extension
// contributes — the same artifact `make codegen` produces for the built-in
// extension. Nothing else is needed to make the types usable.
func (b *ExtensionBuilder) Schema(document string) *ExtensionBuilder {
	b.definition().Schema = document
	return b
}

// Parser adds one generated parser declaration.
func (b *ExtensionBuilder) Parser(name string, parserType *cstxproto.ParserType) *ExtensionBuilder {
	b.definition().Parsers[name] = parserType
	return b
}

// Rule adds one declarative native linker rule.
func (b *ExtensionBuilder) Rule(rule *cstxproto.JoinRule) *ExtensionBuilder {
	b.definition().Rules = append(b.definition().Rules, rule)
	return b
}

// Build returns the canonical generated protobuf message.
func (b *ExtensionBuilder) Build() *cstxproto.ExtensionContract { return b.contract }

// Extensions is the unified extension lifecycle and schema namespace.
type Extensions struct{ eng engine }

// Register atomically registers one canonical protobuf extension contract.
func (e *Extensions) Register(ctx context.Context, contract *cstxproto.ExtensionContract) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	if contract == nil {
		return &Error{Code: CodeInvalidArgument, Operation: "extensions.register", Message: "contract must not be nil"}
	}
	return e.eng.extensionRegister(ctx, contract)
}

// ExportContract returns the canonical registered extension contract.
// Consumers should derive any metadata view from this generated protobuf
// instead of maintaining a second schema registry.
func (e *Extensions) ExportContract(ctx context.Context) (cstxproto.ExtensionContract, error) {
	if err := contextError(ctx); err != nil {
		return cstxproto.ExtensionContract{}, err
	}
	return e.eng.extensionExportContract(ctx)
}

// Enable explicitly enables one linked native Rust extension.
func (e *Extensions) Enable(ctx context.Context, name string) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return e.eng.extensionEnable(ctx, name)
}

// List returns linked and registered extension metadata.
func (e *Extensions) List(ctx context.Context) (*cstxproto.ExtensionCatalog, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return e.eng.extensionList(ctx)
}

// Info returns one linked or registered extension's metadata.
func (e *Extensions) Info(ctx context.Context, name string) (*cstxproto.ExtensionInfo, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return e.eng.extensionInfo(ctx, name)
}

// Contains reports whether a schema exists.
func (e *Extensions) Contains(ctx context.Context, nodeType string) (bool, error) {
	if err := contextError(ctx); err != nil {
		return false, err
	}
	return e.eng.extensionContains(ctx, nodeType)
}

// Schema returns one retained schema.
func (e *Extensions) Schema(ctx context.Context, nodeType string) (cstxproto.NodeType, error) {
	if err := contextError(ctx); err != nil {
		return cstxproto.NodeType{}, err
	}
	return e.eng.extensionSchema(ctx, nodeType)
}

// Schemas returns retained schemas in deterministic order.
func (e *Extensions) Schemas(ctx context.Context) (cstxproto.NodeTypeCatalog, error) {
	if err := contextError(ctx); err != nil {
		return cstxproto.NodeTypeCatalog{}, err
	}
	return e.eng.extensionSchemas(ctx)
}

// HasNativeArtifact reports whether an enabled native parser supports an artifact.
func (e *Extensions) HasNativeArtifact(ctx context.Context, artifact string) (bool, error) {
	if err := contextError(ctx); err != nil {
		return false, err
	}
	return e.eng.extensionHasNativeArtifact(ctx, artifact)
}

// AnchorConcepts lists native concepts and their member node types.
func (e *Extensions) AnchorConcepts(ctx context.Context) (cstxproto.AnchorConceptCatalog, error) {
	if err := contextError(ctx); err != nil {
		return cstxproto.AnchorConceptCatalog{}, err
	}
	return e.eng.extensionAnchorConcepts(ctx)
}
