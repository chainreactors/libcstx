package cstx

import "context"

// Schemas is the schema/plugin namespace of a CSTX runtime.
type Schemas struct{ eng engine }

// Import atomically validates and registers a portable schema contract.
func (s *Schemas) Import(ctx context.Context, contract SchemaContract) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return s.eng.schemaImport(ctx, contract)
}

// Export returns the complete portable schema contract.
func (s *Schemas) Export(ctx context.Context) (SchemaContract, error) {
	if err := contextError(ctx); err != nil {
		return SchemaContract{}, err
	}
	return s.eng.schemaExport(ctx)
}

// Register adds CSTX validation metadata for one node type. An empty
// valueField means no designated value field.
func (s *Schemas) Register(ctx context.Context, nodeType string, schema map[string]any, valueField string) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return s.eng.schemaRegister(ctx, nodeType, schema, valueField)
}

// RegisterJoinRule registers one declarative native linker rule.
func (s *Schemas) RegisterJoinRule(ctx context.Context, rule JoinRuleSpec) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return s.eng.schemaRegisterJoinRule(ctx, rule)
}

// Contains reports whether a schema exists for the node type.
func (s *Schemas) Contains(ctx context.Context, nodeType string) (bool, error) {
	if err := contextError(ctx); err != nil {
		return false, err
	}
	return s.eng.schemaContains(ctx, nodeType)
}

// Get returns one retained schema.
func (s *Schemas) Get(ctx context.Context, nodeType string) (map[string]any, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return s.eng.schemaGet(ctx, nodeType)
}

// List returns retained schemas in deterministic node-type order.
func (s *Schemas) List(ctx context.Context) ([]map[string]any, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return s.eng.schemaList(ctx)
}

// LoadPlugin loads one linked native plugin into the shared graph engine.
func (s *Schemas) LoadPlugin(ctx context.Context, name string) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return s.eng.schemaLoadPlugin(ctx, name)
}

// LoadAllPlugins loads every linked native plugin.
func (s *Schemas) LoadAllPlugins(ctx context.Context) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return s.eng.schemaLoadAllPlugins(ctx)
}

// AvailablePlugins lists linked plugins without changing runtime state.
func (s *Schemas) AvailablePlugins(ctx context.Context) ([]string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return s.eng.schemaAvailablePlugins(ctx)
}

// PluginArtifacts lists artifacts provided by one linked plugin without
// loading it into the runtime.
func (s *Schemas) PluginArtifacts(ctx context.Context, name string) ([]string, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return s.eng.schemaPluginArtifacts(ctx, name)
}

// HasNativeArtifact reports whether a linked native parser supports an artifact.
func (s *Schemas) HasNativeArtifact(ctx context.Context, artifact string) (bool, error) {
	if err := contextError(ctx); err != nil {
		return false, err
	}
	return s.eng.schemaHasNativeArtifact(ctx, artifact)
}

// AnchorConcepts lists native concepts and their member node types.
func (s *Schemas) AnchorConcepts(ctx context.Context) ([]AnchorConcept, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	return s.eng.schemaAnchorConcepts(ctx)
}
