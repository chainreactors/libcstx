package cstx

/*
#include "cstx_ffi.h"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"unsafe"
)

type CASStore struct {
	ptr *C.struct_CstxCasStore
}

type TreeDiff struct {
	Added    map[string][]string `json:"added"`
	Removed  map[string][]string `json:"removed"`
	Modified map[string][]string `json:"modified"`
}

func NewCASStore() *CASStore {
	return &CASStore{ptr: C.cstx_cas_store_new()}
}

func (s *CASStore) Close() {
	if s.ptr != nil {
		C.cstx_cas_store_free(s.ptr)
		s.ptr = nil
	}
}

func CanonicalHash(data []byte) string {
	cJSON := C.CString(string(data))
	defer C.free(unsafe.Pointer(cJSON))
	out := C.cstx_cas_canonical_hash(cJSON)
	if out == nil {
		return ""
	}
	defer C.cstx_free_string(out)
	return C.GoString(out)
}

func (s *CASStore) LoadRoot(hash string, data []byte) error {
	if s == nil || s.ptr == nil {
		return fmt.Errorf("cstx: CAS store is closed")
	}
	cHash := C.CString(hash)
	cJSON := C.CString(string(data))
	defer C.free(unsafe.Pointer(cHash))
	defer C.free(unsafe.Pointer(cJSON))
	if C.cstx_cas_store_load_root(s.ptr, cHash, cJSON) != 0 {
		return fmt.Errorf("cstx: load root failed")
	}
	return nil
}

func (s *CASStore) LoadMap(hash string, data []byte) error {
	if s == nil || s.ptr == nil {
		return fmt.Errorf("cstx: CAS store is closed")
	}
	cHash := C.CString(hash)
	cJSON := C.CString(string(data))
	defer C.free(unsafe.Pointer(cHash))
	defer C.free(unsafe.Pointer(cJSON))
	if C.cstx_cas_store_load_map(s.ptr, cHash, cJSON) != 0 {
		return fmt.Errorf("cstx: load map failed")
	}
	return nil
}

func (s *CASStore) Export(rootHash string) ([]byte, error) {
	if s == nil || s.ptr == nil {
		return nil, fmt.Errorf("cstx: CAS store is closed")
	}
	cRoot := C.CString(rootHash)
	defer C.free(unsafe.Pointer(cRoot))
	out := C.cstx_cas_store_export(s.ptr, cRoot)
	if out == nil {
		return nil, fmt.Errorf("cstx: export failed")
	}
	defer C.cstx_free_string(out)
	return []byte(C.GoString(out)), nil
}

func (s *CASStore) TreeBuild(entriesJSON []byte) (string, error) {
	if s == nil || s.ptr == nil {
		return "", fmt.Errorf("cstx: CAS store is closed")
	}
	cJSON := C.CString(string(entriesJSON))
	defer C.free(unsafe.Pointer(cJSON))
	out := C.cstx_cas_tree_build(s.ptr, cJSON)
	if out == nil {
		return "", fmt.Errorf("cstx: tree build failed")
	}
	defer C.cstx_free_string(out)
	return C.GoString(out), nil
}

func (s *CASStore) TreeUpdate(rootHash string, changesJSON []byte) (string, error) {
	if s == nil || s.ptr == nil {
		return "", fmt.Errorf("cstx: CAS store is closed")
	}
	cRoot := C.CString(rootHash)
	cJSON := C.CString(string(changesJSON))
	defer C.free(unsafe.Pointer(cRoot))
	defer C.free(unsafe.Pointer(cJSON))
	out := C.cstx_cas_tree_update(s.ptr, cRoot, cJSON)
	if out == nil {
		return "", fmt.Errorf("cstx: tree update failed")
	}
	defer C.cstx_free_string(out)
	return C.GoString(out), nil
}

func (s *CASStore) TreeDiff(baseHash, headHash string) (*TreeDiff, error) {
	if s == nil || s.ptr == nil {
		return nil, fmt.Errorf("cstx: CAS store is closed")
	}
	cBase := C.CString(baseHash)
	cHead := C.CString(headHash)
	defer C.free(unsafe.Pointer(cBase))
	defer C.free(unsafe.Pointer(cHead))
	out := C.cstx_cas_tree_diff(s.ptr, cBase, cHead)
	if out == nil {
		return nil, fmt.Errorf("cstx: tree diff failed")
	}
	defer C.cstx_free_string(out)

	var diff TreeDiff
	if err := json.Unmarshal([]byte(C.GoString(out)), &diff); err != nil {
		return nil, err
	}
	return &diff, nil
}

func (s *CASStore) TreeFind(rootHash, elementID string) (string, bool) {
	if s == nil || s.ptr == nil {
		return "", false
	}
	cRoot := C.CString(rootHash)
	cID := C.CString(elementID)
	defer C.free(unsafe.Pointer(cRoot))
	defer C.free(unsafe.Pointer(cID))
	out := C.cstx_cas_tree_find(s.ptr, cRoot, cID)
	if out == nil {
		return "", false
	}
	defer C.cstx_free_string(out)
	return C.GoString(out), true
}

func (s *CASStore) TreeEntries(rootHash string) ([]byte, error) {
	if s == nil || s.ptr == nil {
		return nil, fmt.Errorf("cstx: CAS store is closed")
	}
	cRoot := C.CString(rootHash)
	defer C.free(unsafe.Pointer(cRoot))
	out := C.cstx_cas_tree_entries(s.ptr, cRoot)
	if out == nil {
		return nil, fmt.Errorf("cstx: tree entries failed")
	}
	defer C.cstx_free_string(out)
	return []byte(C.GoString(out)), nil
}
