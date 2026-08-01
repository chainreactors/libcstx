package cstx

import "context"

// Schemas is the schema/plugin namespace of a CSTX runtime.
type Schemas struct{ eng engine }

// Register adds canonical validation metadata for one node type. An empty
// valueField means no designated value field.
func (s *Schemas) Register(ctx context.Context, nodeType string, schema map[string]any, valueField string) error {
	if err := contextError(ctx); err != nil {
		return err
	}
	return s.eng.schemaRegister(ctx, nodeType, schema, valueField)
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
