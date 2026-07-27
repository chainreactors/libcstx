package cstx

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Code is the stable, language-neutral error category shared with the Rust
// core and every other binding.
type Code string

const (
	CodeInvalidArgument   Code = "INVALID_ARGUMENT"
	CodeNotFound          Code = "NOT_FOUND"
	CodeConflict          Code = "CONFLICT"
	CodeValidation        Code = "VALIDATION"
	CodeParse             Code = "PARSE"
	CodeNotInitialized    Code = "NOT_INITIALIZED"
	CodeUnsupported       Code = "UNSUPPORTED"
	CodeIO                Code = "IO"
	CodeCorruptData       Code = "CORRUPT_DATA"
	CodeCursorInvalidated Code = "CURSOR_INVALIDATED"
	CodeInternal          Code = "INTERNAL"
)

// ErrNativeUnavailable is returned by every engine operation when the package
// is compiled without the cstx_native build tag.
var ErrNativeUnavailable = errors.New("cstx: native engine not compiled (build with -tags cstx_native)")

// Error is a stable CSTX failure whose fields can be handled without parsing
// text. It mirrors the Rust CstxError contract.
type Error struct {
	Code      Code   `json:"code"`
	Operation string `json:"operation"`
	ItemIndex *int   `json:"item_index"`
	Field     string `json:"field"`
	Message   string `json:"message"`
	Expected  string `json:"expected"`
	Actual    string `json:"actual"`
}

func (e *Error) Error() string {
	if e.Operation == "" {
		return fmt.Sprintf("cstx: %s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("cstx: %s: %s: %s", e.Operation, e.Code, e.Message)
}

// IsCode reports whether err is a CSTX *Error with the given code.
func IsCode(err error, code Code) bool {
	var cerr *Error
	return errors.As(err, &cerr) && cerr.Code == code
}

// errorJSON is the FFI transport shape; nullable string fields distinguish
// absent context from empty strings on the Rust side.
type errorJSON struct {
	Code      Code    `json:"code"`
	Operation string  `json:"operation"`
	ItemIndex *int    `json:"item_index"`
	Field     *string `json:"field"`
	Message   string  `json:"message"`
	Expected  *string `json:"expected"`
	Actual    *string `json:"actual"`
}

func parseError(data []byte, fallback Code) *Error {
	var wire errorJSON
	if err := json.Unmarshal(data, &wire); err != nil || wire.Code == "" {
		return &Error{Code: fallback, Message: string(data)}
	}
	return &Error{
		Code:      wire.Code,
		Operation: wire.Operation,
		ItemIndex: wire.ItemIndex,
		Field:     deref(wire.Field),
		Message:   wire.Message,
		Expected:  deref(wire.Expected),
		Actual:    deref(wire.Actual),
	}
}

func deref(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
