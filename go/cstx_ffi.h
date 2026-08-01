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

typedef struct CstxEdgeIterator CstxEdgeIterator;

typedef struct CstxHandle CstxHandle;

typedef struct CstxNodeIterator CstxNodeIterator;

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

void cstx_close(struct CstxHandle *handle);

void cstx_free(struct CstxHandle *handle);

CstxStatusCode cstx_last_change_json(struct CstxHandle *handle,
                                     struct CstxBuffer *output,
                                     struct CstxBuffer *error);

CstxStatusCode cstx_schema_register(struct CstxHandle *handle,
                                    struct CstxSlice node_type,
                                    struct CstxSlice schema_json,
                                    struct CstxSlice value_field,
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

CstxStatusCode cstx_graph_add_edges(struct CstxHandle *handle,
                                    struct CstxSlice data,
                                    uint64_t *affected,
                                    struct CstxBuffer *error);

CstxStatusCode cstx_graph_add_nodes_json(struct CstxHandle *handle,
                                         struct CstxSlice data,
                                         uint64_t *affected,
                                         struct CstxBuffer *error);

CstxStatusCode cstx_graph_add_edges_json(struct CstxHandle *handle,
                                         struct CstxSlice data,
                                         uint64_t *affected,
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

CstxStatusCode cstx_graph_stats(struct CstxHandle *handle,
                                struct CstxBuffer *output,
                                struct CstxBuffer *error);

CstxStatusCode cstx_graph_nodes(struct CstxHandle *handle,
                                struct CstxSlice filter_json,
                                struct CstxSlice options_json,
                                struct CstxNodeIterator **output,
                                struct CstxBuffer *error);

CstxStatusCode cstx_graph_edges(struct CstxHandle *handle,
                                struct CstxSlice filter_json,
                                struct CstxSlice options_json,
                                struct CstxEdgeIterator **output,
                                struct CstxBuffer *error);

CstxStatusCode cstx_graph_neighbors(struct CstxHandle *handle,
                                    struct CstxSlice node_id,
                                    struct CstxSlice direction,
                                    struct CstxSlice options_json,
                                    struct CstxNodeIterator **output,
                                    struct CstxBuffer *error);

CstxStatusCode cstx_graph_query(struct CstxHandle *handle,
                                struct CstxSlice expression,
                                struct CstxSlice options_json,
                                struct CstxNodeIterator **output,
                                struct CstxBuffer *error);

CstxStatusCode cstx_graph_ingest_native_json(struct CstxHandle *handle,
                                             struct CstxSlice plugin,
                                             struct CstxSlice artifact,
                                             struct CstxSlice data,
                                             struct CstxBuffer *output,
                                             struct CstxBuffer *error);

CstxStatusCode cstx_graph_hydrate_cas_nodes(struct CstxHandle *handle,
                                            struct CstxSlice items_json,
                                            uint64_t *affected,
                                            struct CstxBuffer *error);

CstxStatusCode cstx_graph_hydrate_cas_edges(struct CstxHandle *handle,
                                            struct CstxSlice items_json,
                                            uint64_t *affected,
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

CstxStatusCode cstx_graph_bfs_json(struct CstxHandle *handle,
                                   struct CstxSlice seed_id,
                                   uint32_t depth,
                                   uint8_t reverse,
                                   struct CstxBuffer *output,
                                   struct CstxBuffer *error);

CstxStatusCode cstx_graph_shortest_paths_json(struct CstxHandle *handle,
                                              struct CstxSlice start_id,
                                              struct CstxSlice end_id,
                                              uint32_t max_depth,
                                              struct CstxBuffer *output,
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

CstxStatusCode cstx_graph_stats_filtered_json(struct CstxHandle *handle,
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

CstxStatusCode cstx_node_iterator_next(struct CstxNodeIterator *cursor,
                                       struct CstxBuffer *output,
                                       uint8_t *has_value,
                                       struct CstxBuffer *error);

void cstx_node_iterator_close(struct CstxNodeIterator *cursor);

void cstx_node_iterator_free(struct CstxNodeIterator *cursor);

CstxStatusCode cstx_edge_iterator_next(struct CstxEdgeIterator *cursor,
                                       struct CstxBuffer *output,
                                       uint8_t *has_value,
                                       struct CstxBuffer *error);

void cstx_edge_iterator_close(struct CstxEdgeIterator *cursor);

void cstx_edge_iterator_free(struct CstxEdgeIterator *cursor);

CstxStatusCode cstx_repo_commit(struct CstxHandle *handle,
                                struct CstxSlice ref_name,
                                struct CstxSlice message,
                                struct CstxSlice metadata_json,
                                struct CstxBuffer *output,
                                struct CstxBuffer *error);

CstxStatusCode cstx_repo_diff(struct CstxHandle *handle,
                              struct CstxSlice base_ref,
                              struct CstxSlice head_ref,
                              struct CstxBuffer *output,
                              struct CstxBuffer *error);

CstxStatusCode cstx_repo_dump(struct CstxHandle *handle,
                              struct CstxSlice compression,
                              struct CstxBuffer *output,
                              struct CstxBuffer *error);

CstxStatusCode cstx_repo_load(struct CstxHandle *handle,
                              struct CstxSlice data,
                              struct CstxSlice compression,
                              uint64_t *consumed,
                              struct CstxBuffer *error);

CstxStatusCode cstx_repo_dump_json(struct CstxHandle *handle,
                                   struct CstxBuffer *output,
                                   struct CstxBuffer *error);

CstxStatusCode cstx_repo_load_json(struct CstxHandle *handle,
                                   struct CstxSlice data,
                                   uint64_t *consumed,
                                   struct CstxBuffer *error);

CstxStatusCode cstx_repo_snapshot_fingerprint(struct CstxHandle *handle,
                                              struct CstxBuffer *output,
                                              struct CstxBuffer *error);

CstxStatusCode cstx_repo_head(struct CstxHandle *handle,
                              struct CstxSlice ref_name,
                              struct CstxBuffer *output,
                              struct CstxBuffer *error);

CstxStatusCode cstx_repo_refs(struct CstxHandle *handle,
                              struct CstxBuffer *output,
                              struct CstxBuffer *error);

CstxStatusCode cstx_repo_index_graph_json(struct CstxHandle *handle,
                                          struct CstxBuffer *output,
                                          struct CstxBuffer *error);

CstxStatusCode cstx_repo_validate_ref(struct CstxHandle *handle,
                                      struct CstxSlice ref_name,
                                      struct CstxBuffer *error);

CstxStatusCode cstx_repo_delete_ref(struct CstxHandle *handle,
                                    struct CstxSlice ref_name,
                                    struct CstxBuffer *error);

CstxStatusCode cstx_repo_set_ref(struct CstxHandle *handle,
                                 struct CstxSlice ref_name,
                                 struct CstxSlice commit_hash,
                                 struct CstxBuffer *error);

CstxStatusCode cstx_repo_fork_json(struct CstxHandle *handle,
                                   struct CstxSlice source,
                                   struct CstxSlice target,
                                   struct CstxBuffer *output,
                                   struct CstxBuffer *error);

CstxStatusCode cstx_repo_checkout_entries_json(struct CstxHandle *handle,
                                               struct CstxSlice ref_name,
                                               struct CstxSlice types_json,
                                               struct CstxBuffer *output,
                                               struct CstxBuffer *error);

CstxStatusCode cstx_repo_commit_entries_json(struct CstxHandle *handle,
                                             struct CstxSlice request_json,
                                             struct CstxBuffer *output,
                                             struct CstxBuffer *error);

CstxStatusCode cstx_repo_commit_delta_json(struct CstxHandle *handle,
                                           struct CstxSlice request_json,
                                           struct CstxBuffer *output,
                                           struct CstxBuffer *error);

CstxStatusCode cstx_repo_merge_json(struct CstxHandle *handle,
                                    struct CstxSlice request_json,
                                    struct CstxBuffer *output,
                                    struct CstxBuffer *error);

CstxStatusCode cstx_repo_merge_base_json(struct CstxHandle *handle,
                                         struct CstxSlice ref_a,
                                         struct CstxSlice ref_b,
                                         struct CstxBuffer *output,
                                         struct CstxBuffer *error);

CstxStatusCode cstx_repo_log_json(struct CstxHandle *handle,
                                  struct CstxSlice ref_name,
                                  size_t limit,
                                  struct CstxBuffer *output,
                                  struct CstxBuffer *error);

CstxStatusCode cstx_repo_tree_stats_json(struct CstxHandle *handle,
                                         struct CstxSlice ref_name,
                                         struct CstxBuffer *output,
                                         struct CstxBuffer *error);

CstxStatusCode cstx_repo_cherry_pick_json(struct CstxHandle *handle,
                                          struct CstxSlice request_json,
                                          struct CstxBuffer *output,
                                          struct CstxBuffer *error);

CstxStatusCode cstx_repo_import_commit(struct CstxHandle *handle,
                                       struct CstxSlice hash,
                                       struct CstxSlice data_json,
                                       struct CstxBuffer *error);

CstxStatusCode cstx_repo_export_commit_json(struct CstxHandle *handle,
                                            struct CstxSlice hash,
                                            struct CstxBuffer *output,
                                            struct CstxBuffer *error);

CstxStatusCode cstx_repo_export_tree_objects_json(struct CstxHandle *handle,
                                                  struct CstxSlice root_hash,
                                                  struct CstxBuffer *output,
                                                  struct CstxBuffer *error);

CstxStatusCode cstx_repo_import_tree_objects(struct CstxHandle *handle,
                                             struct CstxSlice objects_json,
                                             struct CstxBuffer *error);

CstxStatusCode cstx_repo_tree_entries_json(struct CstxHandle *handle,
                                           struct CstxSlice root_hash,
                                           struct CstxSlice types_json,
                                           struct CstxBuffer *output,
                                           struct CstxBuffer *error);

CstxStatusCode cstx_repo_find_tree_entry_json(struct CstxHandle *handle,
                                              struct CstxSlice root_hash,
                                              struct CstxSlice element_id,
                                              struct CstxBuffer *output,
                                              struct CstxBuffer *error);

CstxStatusCode cstx_repo_diff_tree_entries_json(struct CstxHandle *handle,
                                                struct CstxSlice base_root,
                                                struct CstxSlice head_root,
                                                struct CstxBuffer *output,
                                                struct CstxBuffer *error);

CstxStatusCode cstx_repo_tree_root_stats_json(struct CstxHandle *handle,
                                              struct CstxSlice root_hash,
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
