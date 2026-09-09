package cstx

import (
	"encoding/json"
	"sort"

	"github.com/chainreactors/libcstx/go/proto/cstxproto"
)

// FlagRegistry answers what a node flag bit means, from the schema documents
// the runtime holds.
//
// The counterpart of Python's `cstxpy.flags.NodeFlags`. Neither side compiles
// in a vocabulary: a flag's name and its bit are declared by the extension
// that owns them, so a third party's flags are reachable here on exactly the
// same terms as the built-in ones. Bits 0-55 belong to extensions; 56-63 are
// reserved for the runtime.
type FlagRegistry struct {
	bits           map[string]uint32
	defaultExclude uint64
}

type flagDeclaration struct {
	Bit            uint32 `json:"bit"`
	DefaultExclude bool   `json:"default_exclude"`
}

type schemaDocument struct {
	Flags map[string]flagDeclaration `json:"flags"`
}

// Flags reads every loaded extension's declarations out of an exported
// contract. `Extensions.ExportContract` returns what the core accepted, so
// this needs no ABI call of its own.
func Flags(contract *cstxproto.ExtensionContract) (*FlagRegistry, error) {
	registry := &FlagRegistry{bits: map[string]uint32{}}
	if contract == nil {
		return registry, nil
	}
	for _, definition := range contract.GetExtensions() {
		document := definition.GetSchema()
		if document == "" {
			continue
		}
		var parsed schemaDocument
		if err := json.Unmarshal([]byte(document), &parsed); err != nil {
			return nil, err
		}
		for name, declaration := range parsed.Flags {
			// First claimant keeps the bit, as the core decides at
			// registration: a bit is what a stored mask means, so a second
			// claim would make one stored value ambiguous.
			if _, taken := registry.bits[name]; taken {
				continue
			}
			registry.bits[name] = declaration.Bit
			if declaration.DefaultExclude {
				registry.defaultExclude |= uint64(1) << declaration.Bit
			}
		}
	}
	return registry, nil
}

// Bit returns the bit one declared flag occupies, and whether it is declared.
func (r *FlagRegistry) Bit(name string) (uint32, bool) {
	bit, ok := r.bits[name]
	return bit, ok
}

// Mask returns the single-bit mask for one declared flag; 0 if undeclared.
func (r *FlagRegistry) Mask(name string) uint64 {
	bit, ok := r.bits[name]
	if !ok {
		return 0
	}
	return uint64(1) << bit
}

// AllMask returns every bit any loaded extension declared.
func (r *FlagRegistry) AllMask() uint64 {
	var mask uint64
	for _, bit := range r.bits {
		mask |= uint64(1) << bit
	}
	return mask
}

// DefaultExcludeMask returns the bits extensions advise hiding from an
// ordinary view. Advice, not enforcement: nothing applies it on its own.
func (r *FlagRegistry) DefaultExcludeMask() uint64 { return r.defaultExclude }

// Names returns the declared flag names, lowest bit first.
func (r *FlagRegistry) Names() []string {
	names := make([]string, 0, len(r.bits))
	for name := range r.bits {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool { return r.bits[names[i]] < r.bits[names[j]] })
	return names
}
