/* tslint:disable */
/* eslint-disable */

export class CSTX {
    free(): void;
    [Symbol.dispose](): void;
    close(): void;
    constructor(config?: any | null);
    readonly closed: boolean;
    readonly graph: Graph;
    readonly repository: Repository;
    readonly schemas: Schemas;
}

export class EdgeCursor {
    private constructor();
    free(): void;
    [Symbol.dispose](): void;
    close(): void;
    next(): any;
    readonly closed: boolean;
}

export class Graph {
    private constructor();
    free(): void;
    [Symbol.dispose](): void;
    addEdge(edge: any): bigint;
    addEdges(edges: any): bigint;
    addEdgesJson(data: Uint8Array): bigint;
    addNode(node: any): bigint;
    addNodes(nodes: any): bigint;
    addNodesJson(data: Uint8Array): bigint;
    contains(node_id: string): boolean;
    edgeCount(): bigint;
    edges(options?: any | null): EdgeCursor;
    ingest(source: string, data: Uint8Array): bigint;
    neighbors(node_id: string, direction?: string | null, options?: any | null): NodeCursor;
    node(node_id: string): any;
    nodeCount(): bigint;
    nodes(options?: any | null): NodeCursor;
    query(expression: string, options?: any | null): NodeCursor;
    stats(): any;
}

export class NodeCursor {
    private constructor();
    free(): void;
    [Symbol.dispose](): void;
    close(): void;
    next(): any;
    readonly closed: boolean;
}

export class Repository {
    private constructor();
    free(): void;
    [Symbol.dispose](): void;
    commit(ref_name: string | null | undefined, message: string, metadata?: any | null): any;
    diff(base: string, head: string): any;
    dump(compression?: string | null): Uint8Array;
    head(ref_name?: string | null): any;
    load(data: Uint8Array, compression?: string | null): bigint;
    refs(): any;
    snapshotFingerprint(): string;
}

export class Schemas {
    private constructor();
    free(): void;
    [Symbol.dispose](): void;
    availablePlugins(): any;
    contains(node_type: string): boolean;
    get(node_type: string): any;
    list(): any;
    loadAllPlugins(): void;
    loadPlugin(name: string): void;
    register(node_type: string, schema: any, value_field?: string | null): void;
}

export function version(): string;

export type InitInput = RequestInfo | URL | Response | BufferSource | WebAssembly.Module;

export interface InitOutput {
    readonly memory: WebAssembly.Memory;
    readonly __wbg_cstx_free: (a: number, b: number) => void;
    readonly __wbg_edgecursor_free: (a: number, b: number) => void;
    readonly __wbg_graph_free: (a: number, b: number) => void;
    readonly __wbg_nodecursor_free: (a: number, b: number) => void;
    readonly cstx_close: (a: number) => void;
    readonly cstx_closed: (a: number) => number;
    readonly cstx_graph: (a: number) => number;
    readonly cstx_new: (a: number, b: number) => void;
    readonly edgecursor_close: (a: number) => void;
    readonly edgecursor_closed: (a: number) => number;
    readonly edgecursor_next: (a: number, b: number) => void;
    readonly graph_addEdge: (a: number, b: number, c: number) => void;
    readonly graph_addEdges: (a: number, b: number, c: number) => void;
    readonly graph_addEdgesJson: (a: number, b: number, c: number, d: number) => void;
    readonly graph_addNode: (a: number, b: number, c: number) => void;
    readonly graph_addNodes: (a: number, b: number, c: number) => void;
    readonly graph_addNodesJson: (a: number, b: number, c: number, d: number) => void;
    readonly graph_contains: (a: number, b: number, c: number, d: number) => void;
    readonly graph_edgeCount: (a: number, b: number) => void;
    readonly graph_edges: (a: number, b: number, c: number) => void;
    readonly graph_ingest: (a: number, b: number, c: number, d: number, e: number, f: number) => void;
    readonly graph_neighbors: (a: number, b: number, c: number, d: number, e: number, f: number, g: number) => void;
    readonly graph_node: (a: number, b: number, c: number, d: number) => void;
    readonly graph_nodeCount: (a: number, b: number) => void;
    readonly graph_nodes: (a: number, b: number, c: number) => void;
    readonly graph_query: (a: number, b: number, c: number, d: number, e: number) => void;
    readonly graph_stats: (a: number, b: number) => void;
    readonly nodecursor_close: (a: number) => void;
    readonly nodecursor_closed: (a: number) => number;
    readonly nodecursor_next: (a: number, b: number) => void;
    readonly repository_commit: (a: number, b: number, c: number, d: number, e: number, f: number, g: number) => void;
    readonly repository_diff: (a: number, b: number, c: number, d: number, e: number, f: number) => void;
    readonly repository_dump: (a: number, b: number, c: number, d: number) => void;
    readonly repository_head: (a: number, b: number, c: number, d: number) => void;
    readonly repository_load: (a: number, b: number, c: number, d: number, e: number, f: number) => void;
    readonly repository_refs: (a: number, b: number) => void;
    readonly repository_snapshotFingerprint: (a: number, b: number) => void;
    readonly schemas_availablePlugins: (a: number, b: number) => void;
    readonly schemas_contains: (a: number, b: number, c: number, d: number) => void;
    readonly schemas_get: (a: number, b: number, c: number, d: number) => void;
    readonly schemas_list: (a: number, b: number) => void;
    readonly schemas_loadAllPlugins: (a: number, b: number) => void;
    readonly schemas_loadPlugin: (a: number, b: number, c: number, d: number) => void;
    readonly schemas_register: (a: number, b: number, c: number, d: number, e: number, f: number, g: number) => void;
    readonly version: (a: number) => void;
    readonly rust_zstd_wasm_shim_calloc: (a: number, b: number) => number;
    readonly rust_zstd_wasm_shim_free: (a: number) => void;
    readonly rust_zstd_wasm_shim_malloc: (a: number) => number;
    readonly rust_zstd_wasm_shim_memcmp: (a: number, b: number, c: number) => number;
    readonly rust_zstd_wasm_shim_memcpy: (a: number, b: number, c: number) => number;
    readonly rust_zstd_wasm_shim_memmove: (a: number, b: number, c: number) => number;
    readonly rust_zstd_wasm_shim_memset: (a: number, b: number, c: number) => number;
    readonly rust_zstd_wasm_shim_qsort: (a: number, b: number, c: number, d: number) => void;
    readonly cstx_repository: (a: number) => number;
    readonly cstx_schemas: (a: number) => number;
    readonly __wbg_repository_free: (a: number, b: number) => void;
    readonly __wbg_schemas_free: (a: number, b: number) => void;
    readonly __wbindgen_export: (a: number, b: number) => number;
    readonly __wbindgen_export2: (a: number, b: number, c: number, d: number) => number;
    readonly __wbindgen_export3: (a: number) => void;
    readonly __wbindgen_add_to_stack_pointer: (a: number) => number;
    readonly __wbindgen_export4: (a: number, b: number, c: number) => void;
}

export type SyncInitInput = BufferSource | WebAssembly.Module;

/**
 * Instantiates the given `module`, which can either be bytes or
 * a precompiled `WebAssembly.Module`.
 *
 * @param {{ module: SyncInitInput }} module - Passing `SyncInitInput` directly is deprecated.
 *
 * @returns {InitOutput}
 */
export function initSync(module: { module: SyncInitInput } | SyncInitInput): InitOutput;

/**
 * If `module_or_path` is {RequestInfo} or {URL}, makes a request and
 * for everything else, calls `WebAssembly.instantiate` directly.
 *
 * @param {{ module_or_path: InitInput | Promise<InitInput> }} module_or_path - Passing `InitInput` directly is deprecated.
 *
 * @returns {Promise<InitOutput>}
 */
export default function __wbg_init (module_or_path?: { module_or_path: InitInput | Promise<InitInput> } | InitInput | Promise<InitInput>): Promise<InitOutput>;
