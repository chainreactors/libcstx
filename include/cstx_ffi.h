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

CstxStatusCode cstx_open(struct CstxSlice config_json,
                         struct CstxHandle **output,
                         struct CstxBuffer *error);

void cstx_free(struct CstxHandle *handle);

CstxStatusCode cstx_last_change_json(struct CstxHandle *handle,
                                     struct CstxBuffer *output,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_schema_register(struct CstxHandle *handle,
                                    struct CstxSlice node_type,
                                    struct CstxSlice schema_json,
                                    struct CstxSlice value_field,
                                    struct CstxBuffer *error);

CstxStatusCode cstx_schema_import_schema(struct CstxHandle *handle,
                                         struct CstxSlice contract_json,
                                         struct CstxBuffer *error);

CstxStatusCode cstx_schema_export_schema_json(struct CstxHandle *handle,
                                              struct CstxBuffer *output,
                                              struct CstxBuffer *error);

CstxStatusCode cstx_schema_contains(struct CstxHandle *handle,
                                    struct CstxSlice node_type,
                                    uint8_t *output,
                                    struct CstxBuffer *error);

CstxStatusCode cstx_schema_list_json(struct CstxHandle *handle,
                                     struct CstxBuffer *output,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_schema_get_json(struct CstxHandle *handle,
                                    struct CstxSlice node_type,
                                    struct CstxBuffer *output,
                                    struct CstxBuffer *error);

CstxStatusCode cstx_schema_load_plugin(struct CstxHandle *handle,
                                       struct CstxSlice name,
                                       struct CstxBuffer *error);

CstxStatusCode cstx_schema_load_all_plugins(struct CstxHandle *handle, struct CstxBuffer *error);

CstxStatusCode cstx_schema_available_plugins_json(struct CstxHandle *handle,
                                                  struct CstxBuffer *output,
                                                  struct CstxBuffer *error);

CstxStatusCode cstx_schema_plugin_artifacts_json(struct CstxHandle *handle,
                                                 struct CstxSlice name,
                                                 struct CstxBuffer *output,
                                                 struct CstxBuffer *error);

CstxStatusCode cstx_schema_register_join_rule(struct CstxHandle *handle,
                                              struct CstxSlice rule_json,
                                              struct CstxBuffer *error);

CstxStatusCode cstx_schema_has_native_artifact(struct CstxHandle *handle,
                                               struct CstxSlice artifact,
                                               uint8_t *output,
                                               struct CstxBuffer *error);

CstxStatusCode cstx_schema_anchor_concepts_json(struct CstxHandle *handle,
                                                struct CstxBuffer *output,
                                                struct CstxBuffer *error);

CstxStatusCode cstx_graph_add_nodes(struct CstxHandle *handle,
                                    struct CstxSlice data,
                                    uint64_t *affected,
                                    struct CstxBuffer *error);

/**
 * Write each node as its current state, replacing the stored record.
 *
 * The merge path (`cstx_graph_add_nodes`) owns bulk ingest and keeps its JSON
 * fast path. A replace batch is a caller restating records it already holds —
 * a task's oracles, a document's current revision — so it goes through the
 * shared `Value` path rather than earning a second parser.
 */
CstxStatusCode cstx_graph_replace_nodes(struct CstxHandle *handle,
                                        struct CstxSlice data,
                                        uint64_t *affected,
                                        struct CstxBuffer *error);

CstxStatusCode cstx_graph_add_edges(struct CstxHandle *handle,
                                    struct CstxSlice data,
                                    uint64_t *affected,
                                    struct CstxBuffer *error);

CstxStatusCode cstx_graph_delete_nodes(struct CstxHandle *handle,
                                       struct CstxSlice node_ids_json,
                                       uint64_t *output,
                                       struct CstxBuffer *error);

CstxStatusCode cstx_graph_delete_edges(struct CstxHandle *handle,
                                       struct CstxSlice edge_ids_json,
                                       uint64_t *output,
                                       struct CstxBuffer *error);

CstxStatusCode cstx_graph_ingest(struct CstxHandle *handle,
                                 struct CstxSlice source,
                                 struct CstxSlice data,
                                 uint64_t *affected,
                                 struct CstxBuffer *error);

CstxStatusCode cstx_graph_node(struct CstxHandle *handle,
                               struct CstxSlice node_id,
                               struct CstxBuffer *output,
                               struct CstxBuffer *error);

CstxStatusCode cstx_graph_edge(struct CstxHandle *handle,
                               struct CstxSlice edge_id,
                               struct CstxBuffer *output,
                               struct CstxBuffer *error);

CstxStatusCode cstx_graph_contains(struct CstxHandle *handle,
                                   struct CstxSlice node_id,
                                   uint8_t *output,
                                   struct CstxBuffer *error);

CstxStatusCode cstx_graph_node_count(struct CstxHandle *handle,
                                     uint64_t *output,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_graph_edge_count(struct CstxHandle *handle,
                                     uint64_t *output,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_graph_nodes(struct CstxHandle *handle,
                                struct CstxSlice filter_json,
                                struct CstxSlice options_json,
                                struct CstxGraphCursor **output,
                                struct CstxBuffer *error);

CstxStatusCode cstx_graph_edges(struct CstxHandle *handle,
                                struct CstxSlice filter_json,
                                struct CstxSlice options_json,
                                struct CstxGraphCursor **output,
                                struct CstxBuffer *error);

CstxStatusCode cstx_graph_neighbors(struct CstxHandle *handle,
                                    struct CstxSlice node_id,
                                    struct CstxSlice direction,
                                    struct CstxSlice options_json,
                                    struct CstxGraphCursor **output,
                                    struct CstxBuffer *error);

CstxStatusCode cstx_graph_query(struct CstxHandle *handle,
                                struct CstxSlice expression,
                                struct CstxSlice options_json,
                                struct CstxGraphCursor **output,
                                struct CstxBuffer *error);

CstxStatusCode cstx_graph_ingest_native_json(struct CstxHandle *handle,
                                             struct CstxSlice plugin,
                                             struct CstxSlice artifact,
                                             struct CstxSlice data,
                                             struct CstxBuffer *output,
                                             struct CstxBuffer *error);

CstxStatusCode cstx_graph_find_node_json(struct CstxHandle *handle,
                                         struct CstxSlice identifier,
                                         struct CstxBuffer *output,
                                         struct CstxBuffer *error);

CstxStatusCode cstx_graph_patch_node_extras(struct CstxHandle *handle,
                                            struct CstxSlice node_ids_json,
                                            struct CstxSlice patch_json,
                                            uint64_t *affected,
                                            struct CstxBuffer *error);

CstxStatusCode cstx_graph_create_relationship_json(struct CstxHandle *handle,
                                                   struct CstxSlice request_json,
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

CstxStatusCode cstx_graph_node_types_json(struct CstxHandle *handle,
                                          struct CstxBuffer *output,
                                          struct CstxBuffer *error);

CstxStatusCode cstx_graph_link_json(struct CstxHandle *handle,
                                    struct CstxSlice node_ids_json,
                                    struct CstxSlice data_source,
                                    struct CstxBuffer *output,
                                    struct CstxBuffer *error);

CstxStatusCode cstx_graph_update_node_flags(struct CstxHandle *handle,
                                            struct CstxSlice request_json,
                                            uint64_t *affected,
                                            struct CstxBuffer *error);

CstxStatusCode cstx_graph_analyze(struct CstxHandle *handle,
                                  struct CstxSlice algorithm_json,
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
                                   struct CstxSlice seed_ids_json,
                                   uint32_t depth,
                                   struct CstxHandle **output,
                                   struct CstxBuffer *error);

CstxStatusCode cstx_graph_query_subgraph(struct CstxHandle *handle,
                                         struct CstxSlice request_json,
                                         struct CstxHandle **output,
                                         struct CstxBuffer *error);

CstxStatusCode cstx_graph_induced_subgraph(struct CstxHandle *handle,
                                           struct CstxSlice request_json,
                                           struct CstxHandle **output,
                                           struct CstxBuffer *error);

CstxStatusCode cstx_graph_filter(struct CstxHandle *handle,
                                 struct CstxSlice request_json,
                                 struct CstxHandle **output,
                                 struct CstxBuffer *error);

CstxStatusCode cstx_graph_filter_with_reasons_json(struct CstxHandle *handle,
                                                   struct CstxSlice request_json,
                                                   struct CstxHandle **output,
                                                   struct CstxBuffer *details_json,
                                                   struct CstxBuffer *error);

CstxStatusCode cstx_graph_find_anchors_json(struct CstxHandle *handle,
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

CstxStatusCode cstx_graph_nodes_page_json(struct CstxHandle *handle,
                                          struct CstxSlice request_json,
                                          struct CstxBuffer *output,
                                          struct CstxBuffer *error);

CstxStatusCode cstx_graph_nodes_json(struct CstxHandle *handle,
                                     struct CstxSlice node_type,
                                     struct CstxBuffer *output,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_graph_edges_json(struct CstxHandle *handle,
                                     struct CstxSlice source_id,
                                     struct CstxSlice target_id,
                                     struct CstxSlice relation,
                                     struct CstxBuffer *output,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_graph_neighbors_json(struct CstxHandle *handle,
                                         struct CstxSlice node_id,
                                         struct CstxSlice direction,
                                         struct CstxBuffer *output,
                                         struct CstxBuffer *error);

CstxStatusCode cstx_graph_query_json(struct CstxHandle *handle,
                                     struct CstxSlice expression,
                                     size_t limit,
                                     uint8_t has_limit,
                                     struct CstxBuffer *output,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_graph_cursor_page(struct CstxGraphCursor *cursor,
                                      size_t limit,
                                      size_t page,
                                      struct CstxBuffer *output,
                                      struct CstxBuffer *error);

void cstx_graph_cursor_free(struct CstxGraphCursor *cursor);

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
                                struct CstxSlice metadata_json,
                                struct CstxBuffer *output,
                                struct CstxBuffer *error);

CstxStatusCode cstx_repo_prepare(struct CstxHandle *handle,
                                 struct CstxSlice message,
                                 struct CstxSlice ref_name,
                                 struct CstxSlice expected_head,
                                 struct CstxSlice metadata_json,
                                 int64_t timestamp,
                                 uint8_t has_timestamp,
                                 struct CstxBuffer *output,
                                 struct CstxBuffer *error);

CstxStatusCode cstx_repo_accept(struct CstxHandle *handle,
                                struct CstxSlice commit,
                                struct CstxBuffer *error);

CstxStatusCode cstx_repo_discard(struct CstxHandle *handle, struct CstxBuffer *error);

CstxStatusCode cstx_repo_synchronize(struct CstxHandle *handle,
                                     struct CstxSlice payload_json,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_repo_contains(struct CstxHandle *handle,
                                  struct CstxSlice object,
                                  uint8_t *output,
                                  struct CstxBuffer *error);

CstxStatusCode cstx_repo_missing_tree(struct CstxHandle *handle,
                                      struct CstxSlice commit,
                                      struct CstxBuffer *output,
                                      struct CstxBuffer *error);

CstxStatusCode cstx_repo_object_closure(struct CstxHandle *handle,
                                        struct CstxSlice commit,
                                        struct CstxBuffer *output,
                                        struct CstxBuffer *error);

CstxStatusCode cstx_repo_missing_prepare(struct CstxHandle *handle,
                                         struct CstxSlice commit,
                                         struct CstxBuffer *output,
                                         struct CstxBuffer *error);

CstxStatusCode cstx_repo_missing_history(struct CstxHandle *handle,
                                         struct CstxSlice commit,
                                         struct CstxSlice entity_id,
                                         struct CstxBuffer *output,
                                         struct CstxBuffer *error);

CstxStatusCode cstx_repo_missing_stat(struct CstxHandle *handle,
                                      struct CstxSlice commit,
                                      struct CstxBuffer *output,
                                      struct CstxBuffer *error);

CstxStatusCode cstx_repo_missing_commits(struct CstxHandle *handle,
                                         struct CstxSlice commit,
                                         size_t limit,
                                         struct CstxBuffer *output,
                                         struct CstxBuffer *error);

CstxStatusCode cstx_repo_missing_diff(struct CstxHandle *handle,
                                      struct CstxSlice base,
                                      struct CstxSlice head,
                                      struct CstxSlice detail,
                                      struct CstxBuffer *output,
                                      struct CstxBuffer *error);

CstxStatusCode cstx_repo_missing_delta(struct CstxHandle *handle,
                                       struct CstxSlice commit,
                                       int64_t start_timestamp,
                                       uint8_t has_start,
                                       int64_t end_timestamp,
                                       uint8_t has_end,
                                       struct CstxBuffer *output,
                                       struct CstxBuffer *error);

CstxStatusCode cstx_repo_missing_merge(struct CstxHandle *handle,
                                       struct CstxSlice source,
                                       struct CstxSlice target,
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
                              struct CstxSlice request_json,
                              struct CstxRagIndexSession **output,
                              struct CstxBuffer *error);

CstxStatusCode cstx_rag_index_session_metadata_json(struct CstxRagIndexSession *session,
                                                    struct CstxBuffer *output,
                                                    struct CstxBuffer *error);

CstxStatusCode cstx_rag_index_session_pending_json(struct CstxRagIndexSession *session,
                                                   size_t offset,
                                                   size_t limit,
                                                   struct CstxBuffer *output,
                                                   struct CstxBuffer *error);

CstxStatusCode cstx_rag_index_session_deletes_json(struct CstxRagIndexSession *session,
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
                                 struct CstxSlice query_json,
                                 struct CstxRagRetrieval **output,
                                 struct CstxBuffer *error);

CstxStatusCode cstx_rag_retrieval_requests_json(struct CstxRagRetrieval *retrieval,
                                                struct CstxBuffer *output,
                                                struct CstxBuffer *error);

CstxStatusCode cstx_rag_retrieval_complete_json(struct CstxRagRetrieval *retrieval,
                                                struct CstxSlice batches_json,
                                                struct CstxBuffer *output,
                                                struct CstxBuffer *error);

void cstx_rag_retrieval_close(struct CstxRagRetrieval *retrieval);

void cstx_rag_retrieval_free(struct CstxRagRetrieval *retrieval);

#endif  /* CSTX_FFI_H */
