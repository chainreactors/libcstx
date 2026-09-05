package cstx

import (
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
	CodeStaleOperation    Code = "STALE_OPERATION"
	CodeStaleRecall       Code = "STALE_RECALL"
	CodeInternal          Code = "INTERNAL"
)

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

// parseError decodes the error channel.  Structured runtime payloads use
// protobuf; failures intentionally remain a compact UTF-8 diagnostic paired
// with the CstxStatusCode so callers never need a second error codec.
func parseError(data []byte, fallback Code) *Error {
	return &Error{Code: fallback, Message: string(data)}
}
