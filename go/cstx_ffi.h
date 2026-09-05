/* Generated from cstx-abi by cbindgen — do not edit */

#ifndef CSTX_FFI_H
#define CSTX_FFI_H

#include <stdarg.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>

enum CstxStatusCode
#if __STDC_VERSION__ >= 202311L
  : int32_t
#endif // __STDC_VERSION__ >= 202311L
 {
  CSTX_OK = 0,
  CSTX_INVALID_ARGUMENT = 1,
  CSTX_NOT_FOUND = 2,
  CSTX_CONFLICT = 3,
  CSTX_NOT_INITIALIZED = 4,
  CSTX_UNSUPPORTED = 5,
  CSTX_INTERNAL = 6,
  CSTX_VALIDATION = 7,
  CSTX_PARSE = 8,
  CSTX_IO = 9,
  CSTX_CORRUPT_DATA = 10,
  CSTX_CURSOR_INVALIDATED = 11,
  CSTX_STALE_OPERATION = 12,
  CSTX_STALE_RECALL = 13,
};
#if __STDC_VERSION__ >= 202311L
typedef enum CstxStatusCode CstxStatusCode;
#else
typedef int32_t CstxStatusCode;
#endif // __STDC_VERSION__ >= 202311L

typedef struct CstxGraphCursor CstxGraphCursor;

typedef struct CstxHandle CstxHandle;

typedef struct CstxRagIndexSession CstxRagIndexSession;

typedef struct CstxRagRecordIterator CstxRagRecordIterator;

typedef struct CstxRagRetrieval CstxRagRetrieval;

typedef struct CstxBuffer {
  uint8_t *data;
  size_t len;
} CstxBuffer;

typedef struct CstxSlice {
  const uint8_t *data;
  size_t len;
} CstxSlice;

/**
 * Release a Rust-owned output or error buffer.
 */
void cstx_buffer_free(struct CstxBuffer *buffer);

/**
 * Open a runtime from the canonical protobuf configuration message.
 */
CstxStatusCode cstx_open(struct CstxSlice config,
                         struct CstxHandle **output,
                         struct CstxBuffer *error);

void cstx_free(struct CstxHandle *handle);

/**
 * Return the last graph mutation as a protobuf message.
 */
CstxStatusCode cstx_last_change(struct CstxHandle *handle,
                                struct CstxBuffer *output,
                                struct CstxBuffer *error);

/**
 * Register an extension contract encoded as protobuf.
 */
CstxStatusCode cstx_extension_register(struct CstxHandle *handle,
                                       struct CstxSlice contract,
                                       struct CstxBuffer *error);

/**
 * Explicitly enable one linked native Rust extension.
 */
CstxStatusCode cstx_extension_enable(struct CstxHandle *handle,
                                     struct CstxSlice name,
                                     struct CstxBuffer *error);

/**
 * List extension metadata as protobuf.
 */
CstxStatusCode cstx_extension_list(struct CstxHandle *handle,
                                   struct CstxBuffer *output,
                                   struct CstxBuffer *error);

/**
 * Return extension metadata as protobuf.
 */
CstxStatusCode cstx_extension_info(struct CstxHandle *handle,
                                   struct CstxSlice name,
                                   struct CstxBuffer *output,
                                   struct CstxBuffer *error);

/**
 * Export the extension contract as protobuf for low-level synchronization.
 */
CstxStatusCode cstx_extension_export_contract(struct CstxHandle *handle,
                                              struct CstxBuffer *output,
                                              struct CstxBuffer *error);

/**
 * Test whether an extension has registered a schema for a node type.
 */
CstxStatusCode cstx_extension_contains(struct CstxHandle *handle,
                                       struct CstxSlice node_type,
                                       uint8_t *output,
                                       struct CstxBuffer *error);

CstxStatusCode cstx_extension_schema(struct CstxHandle *handle,
                                     struct CstxSlice node_type,
                                     struct CstxBuffer *output,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_extension_schemas(struct CstxHandle *handle,
                                      struct CstxBuffer *output,
                                      struct CstxBuffer *error);

/**
 * Test whether an enabled native extension provides an artifact parser.
 */
CstxStatusCode cstx_extension_has_native_artifact(struct CstxHandle *handle,
                                                  struct CstxSlice artifact,
                                                  uint8_t *output,
                                                  struct CstxBuffer *error);

CstxStatusCode cstx_extension_anchor_concepts(struct CstxHandle *handle,
                                              struct CstxBuffer *output,
                                              struct CstxBuffer *error);

/**
 * Add or merge a protobuf graph aggregate at the Rust-owned semantic boundary.
 */
CstxStatusCode cstx_graph_add_nodes(struct CstxHandle *handle,
                                    struct CstxSlice data,
                                    uint64_t *affected,
                                    struct CstxBuffer *error);

/**
 * Replace the current graph content from a protobuf aggregate.
 */
CstxStatusCode cstx_graph_replace_nodes(struct CstxHandle *handle,
                                        struct CstxSlice data,
                                        uint64_t *affected,
                                        struct CstxBuffer *error);

/**
 * Add or merge relationships from a protobuf graph aggregate.
 */
CstxStatusCode cstx_graph_add_relationships(struct CstxHandle *handle,
                                            struct CstxSlice data,
                                            uint64_t *affected,
                                            struct CstxBuffer *error);

CstxStatusCode cstx_graph_delete_nodes(struct CstxHandle *handle,
                                       struct CstxSlice node_ids,
                                       uint64_t *output,
                                       struct CstxBuffer *error);

CstxStatusCode cstx_graph_delete_relationships(struct CstxHandle *handle,
                                               struct CstxSlice relationship_ids,
                                               uint64_t *output,
                                               struct CstxBuffer *error);

/**
 * Return one node as a protobuf envelope.
 */
CstxStatusCode cstx_graph_node(struct CstxHandle *handle,
                               struct CstxSlice node_id,
                               struct CstxBuffer *output,
                               struct CstxBuffer *error);

/**
 * Return one relationship as a protobuf envelope.
 */
CstxStatusCode cstx_graph_relationship(struct CstxHandle *handle,
                                       struct CstxSlice relationship_id,
                                       struct CstxBuffer *output,
                                       struct CstxBuffer *error);

CstxStatusCode cstx_graph_contains(struct CstxHandle *handle,
                                   struct CstxSlice node_id,
                                   uint8_t *output,
                                   struct CstxBuffer *error);

CstxStatusCode cstx_graph_node_count(struct CstxHandle *handle,
                                     uint64_t *output,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_graph_relationship_count(struct CstxHandle *handle,
                                             uint64_t *output,
                                             struct CstxBuffer *error);

/**
 * Create a node cursor from a protobuf `NodeQuery` (filter + window).
 */
CstxStatusCode cstx_graph_nodes(struct CstxHandle *handle,
                                struct CstxSlice request,
                                struct CstxGraphCursor **output,
                                struct CstxBuffer *error);

/**
 * Create a relationship cursor from a protobuf `RelationshipQuery` (filter + window).
 */
CstxStatusCode cstx_graph_relationships(struct CstxHandle *handle,
                                        struct CstxSlice request,
                                        struct CstxGraphCursor **output,
                                        struct CstxBuffer *error);

/**
 * Create a neighbor cursor from a semantic query.
 */
CstxStatusCode cstx_graph_neighbors(struct CstxHandle *handle,
                                    struct CstxSlice request,
                                    struct CstxGraphCursor **output,
                                    struct CstxBuffer *error);

/**
 * Create a query cursor from a semantic query.
 */
CstxStatusCode cstx_graph_query(struct CstxHandle *handle,
                                struct CstxSlice request,
                                struct CstxGraphCursor **output,
                                struct CstxBuffer *error);

CstxStatusCode cstx_graph_ingest(struct CstxHandle *handle,
                                 struct CstxSlice request,
                                 struct CstxBuffer *output,
                                 struct CstxBuffer *error);

/**
 * Resolve an identifier and return the matching node as protobuf.
 */
CstxStatusCode cstx_graph_find_node(struct CstxHandle *handle,
                                    struct CstxSlice identifier,
                                    struct CstxBuffer *output,
                                    struct CstxBuffer *error);

CstxStatusCode cstx_graph_patch_node_extras(struct CstxHandle *handle,
                                            struct CstxSlice request,
                                            uint64_t *affected,
                                            struct CstxBuffer *error);

CstxStatusCode cstx_graph_add_relationship(struct CstxHandle *handle,
                                           struct CstxSlice request_bytes,
                                           struct CstxBuffer *output,
                                           struct CstxBuffer *error);

CstxStatusCode cstx_is_path_expression(struct CstxSlice expression,
                                       uint8_t *output,
                                       struct CstxBuffer *error);

CstxStatusCode cstx_graph_union(struct CstxHandle *left,
                                struct CstxHandle *right,
                                struct CstxHandle **output,
                                struct CstxBuffer *error);

CstxStatusCode cstx_graph_merge(struct CstxHandle *target,
                                struct CstxHandle *source,
                                uint64_t *affected,
                                struct CstxBuffer *error);

CstxStatusCode cstx_graph_difference(struct CstxHandle *left,
                                     struct CstxHandle *right,
                                     struct CstxSlice node_type,
                                     struct CstxHandle **output,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_graph_node_types(struct CstxHandle *handle,
                                     struct CstxBuffer *output,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_graph_link(struct CstxHandle *handle,
                               struct CstxSlice node_ids,
                               struct CstxSlice data_source,
                               struct CstxBuffer *output,
                               struct CstxBuffer *error);

CstxStatusCode cstx_graph_update_node_flags(struct CstxHandle *handle,
                                            struct CstxSlice request_bytes,
                                            uint64_t *affected,
                                            struct CstxBuffer *error);

CstxStatusCode cstx_graph_analyze(struct CstxHandle *handle,
                                  struct CstxSlice algorithm_bytes,
                                  struct CstxSlice selection,
                                  uint8_t *kind,
                                  uint8_t *boolean,
                                  struct CstxGraphCursor **cursor,
                                  struct CstxBuffer *error);

CstxStatusCode cstx_graph_degree(struct CstxHandle *handle,
                                 struct CstxSlice node_id,
                                 struct CstxSlice direction,
                                 uint64_t *output,
                                 struct CstxBuffer *error);

CstxStatusCode cstx_graph_subgraph(struct CstxHandle *handle,
                                   struct CstxSlice seed_ids,
                                   uint32_t depth,
                                   struct CstxHandle **output,
                                   struct CstxBuffer *error);

CstxStatusCode cstx_graph_query_subgraph(struct CstxHandle *handle,
                                         struct CstxSlice request_bytes,
                                         struct CstxHandle **output,
                                         struct CstxBuffer *error);

CstxStatusCode cstx_graph_induced_subgraph(struct CstxHandle *handle,
                                           struct CstxSlice request_bytes,
                                           struct CstxHandle **output,
                                           struct CstxBuffer *error);

CstxStatusCode cstx_graph_filter(struct CstxHandle *handle,
                                 struct CstxSlice request_bytes,
                                 struct CstxHandle **output,
                                 struct CstxBuffer *error);

CstxStatusCode cstx_graph_filter_with_reasons(struct CstxHandle *handle,
                                              struct CstxSlice request_bytes,
                                              struct CstxHandle **output,
                                              struct CstxBuffer *details,
                                              struct CstxBuffer *error);

CstxStatusCode cstx_graph_find_anchors(struct CstxHandle *handle,
                                       struct CstxSlice concept_name,
                                       struct CstxBuffer *output,
                                       struct CstxBuffer *error);

CstxStatusCode cstx_graph_elevate(struct CstxHandle *handle,
                                  struct CstxSlice concept_name,
                                  struct CstxHandle **output,
                                  struct CstxBuffer *error);

CstxStatusCode cstx_graph_stats(struct CstxHandle *handle,
                                uint64_t exclude_mask,
                                uint64_t include_mask,
                                struct CstxBuffer *output,
                                struct CstxBuffer *error);

/**
 * Materialize one cursor page as protobuf bytes.
 */
CstxStatusCode cstx_graph_cursor_page(struct CstxGraphCursor *cursor,
                                      size_t limit,
                                      size_t page,
                                      struct CstxBuffer *output,
                                      struct CstxBuffer *error);

void cstx_graph_cursor_free(struct CstxGraphCursor *cursor);

/**
 * Resolve a revision and return its UTF-8 commit id in `output`.
 */
CstxStatusCode cstx_repo_resolve(struct CstxHandle *handle,
                                 struct CstxSlice revision,
                                 struct CstxBuffer *output,
                                 struct CstxBuffer *error);

CstxStatusCode cstx_repo_checkout(struct CstxHandle *handle,
                                  struct CstxSlice revision,
                                  uint8_t force,
                                  struct CstxBuffer *output,
                                  struct CstxBuffer *error);

CstxStatusCode cstx_repo_commit(struct CstxHandle *handle,
                                struct CstxSlice message,
                                struct CstxSlice ref_name,
                                struct CstxSlice expected_head,
                                struct CstxSlice metadata,
                                struct CstxBuffer *output,
                                struct CstxBuffer *error);

CstxStatusCode cstx_repo_prepare(struct CstxHandle *handle,
                                 struct CstxSlice message,
                                 struct CstxSlice ref_name,
                                 struct CstxSlice expected_head,
                                 struct CstxSlice metadata,
                                 int64_t timestamp,
                                 uint8_t has_timestamp,
                                 struct CstxBuffer *output,
                                 struct CstxBuffer *error);

CstxStatusCode cstx_repo_accept(struct CstxHandle *handle,
                                struct CstxSlice commit,
                                struct CstxBuffer *error);

CstxStatusCode cstx_repo_discard(struct CstxHandle *handle, struct CstxBuffer *error);

CstxStatusCode cstx_repo_synchronize(struct CstxHandle *handle,
                                     struct CstxSlice payload_bytes,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_repo_contains(struct CstxHandle *handle,
                                  struct CstxSlice object,
                                  uint8_t *output,
                                  struct CstxBuffer *error);

CstxStatusCode cstx_repo_missing(struct CstxHandle *handle,
                                 struct CstxSlice request_bytes,
                                 struct CstxBuffer *output,
                                 struct CstxBuffer *error);

CstxStatusCode cstx_repo_release_transient_objects(struct CstxHandle *handle,
                                                   struct CstxBuffer *error);

CstxStatusCode cstx_repo_diff(struct CstxHandle *handle,
                              struct CstxSlice base_ref,
                              struct CstxSlice head_ref,
                              size_t limit,
                              uint8_t has_limit,
                              struct CstxSlice detail,
                              struct CstxBuffer *output,
                              struct CstxBuffer *error);

/**
 * Return the UTF-8 commit id at a ref, or an empty buffer when it is absent.
 */
CstxStatusCode cstx_repo_head(struct CstxHandle *handle,
                              struct CstxSlice ref_name,
                              struct CstxBuffer *output,
                              struct CstxBuffer *error);

CstxStatusCode cstx_repo_log(struct CstxHandle *handle,
                             struct CstxSlice revision,
                             size_t limit,
                             struct CstxBuffer *output,
                             struct CstxBuffer *error);

CstxStatusCode cstx_repo_history(struct CstxHandle *handle,
                                 struct CstxSlice entity_id,
                                 struct CstxSlice revision,
                                 size_t limit,
                                 uint8_t has_limit,
                                 struct CstxBuffer *output,
                                 struct CstxBuffer *error);

/**
 * Create a ref and return the target UTF-8 commit id in `output`.
 */
CstxStatusCode cstx_repo_branch(struct CstxHandle *handle,
                                struct CstxSlice name,
                                struct CstxSlice start_point,
                                struct CstxBuffer *output,
                                struct CstxBuffer *error);

CstxStatusCode cstx_repo_merge(struct CstxHandle *handle,
                               struct CstxSlice source,
                               struct CstxSlice target,
                               struct CstxSlice expected_head,
                               struct CstxSlice message,
                               struct CstxBuffer *output,
                               struct CstxBuffer *error);

CstxStatusCode cstx_repo_stat(struct CstxHandle *handle,
                              struct CstxSlice revision,
                              uint64_t exclude_mask,
                              uint64_t include_mask,
                              struct CstxBuffer *output,
                              struct CstxBuffer *error);

CstxStatusCode cstx_repo_delta(struct CstxHandle *handle,
                               struct CstxSlice revision,
                               int64_t start_timestamp,
                               uint8_t has_start_timestamp,
                               int64_t end_timestamp,
                               uint8_t has_end_timestamp,
                               struct CstxBuffer *output,
                               struct CstxBuffer *error);

CstxStatusCode cstx_rag_index(struct CstxHandle *handle,
                              struct CstxSlice request_bytes,
                              struct CstxRagIndexSession **output,
                              struct CstxBuffer *error);

CstxStatusCode cstx_rag_index_session_metadata(struct CstxRagIndexSession *session,
                                               struct CstxBuffer *output,
                                               struct CstxBuffer *error);

CstxStatusCode cstx_rag_index_session_pending(struct CstxRagIndexSession *session,
                                              size_t offset,
                                              size_t limit,
                                              struct CstxBuffer *output,
                                              struct CstxBuffer *error);

CstxStatusCode cstx_rag_index_session_deletes(struct CstxRagIndexSession *session,
                                              struct CstxBuffer *output,
                                              struct CstxBuffer *error);

CstxStatusCode cstx_rag_index_session_records(struct CstxRagIndexSession *session,
                                              struct CstxRagRecordIterator **output,
                                              struct CstxBuffer *error);

CstxStatusCode cstx_rag_record_iterator_next(struct CstxRagRecordIterator *iterator,
                                             struct CstxBuffer *output,
                                             uint8_t *has_value,
                                             struct CstxBuffer *error);

void cstx_rag_record_iterator_close(struct CstxRagRecordIterator *iterator);

void cstx_rag_record_iterator_free(struct CstxRagRecordIterator *iterator);

void cstx_rag_index_session_close(struct CstxRagIndexSession *session);

void cstx_rag_index_session_free(struct CstxRagIndexSession *session);

CstxStatusCode cstx_rag_retrieve(struct CstxHandle *handle,
                                 struct CstxSlice query_bytes,
                                 struct CstxRagRetrieval **output,
                                 struct CstxBuffer *error);

CstxStatusCode cstx_rag_retrieval_requests(struct CstxRagRetrieval *retrieval,
                                           struct CstxBuffer *output,
                                           struct CstxBuffer *error);

CstxStatusCode cstx_rag_retrieval_complete(struct CstxRagRetrieval *retrieval,
                                           struct CstxSlice batches_bytes,
                                           struct CstxBuffer *output,
                                           struct CstxBuffer *error);

void cstx_rag_retrieval_close(struct CstxRagRetrieval *retrieval);

void cstx_rag_retrieval_free(struct CstxRagRetrieval *retrieval);

#endif  /* CSTX_FFI_H */
