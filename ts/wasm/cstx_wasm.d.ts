/* tslint:disable */
/* eslint-disable */

export class CSTXGraph {
    free(): void;
    [Symbol.dispose](): void;
    addEdge(source_id: string, target_id: string, relation: string, data_source: string): void;
    addEdgesBatch(edges_json: string): number;
    addJoinRule(left_type: string, right_type: string, relation: string, left_key: string, right_key: string, predicted: boolean): void;
    addNode(node_type: string, cstx_id: string, payload_json: string, sources_json: string, flags: bigint): void;
    addNodesBatch(nodes_json: string): string;
    allNodesJson(): string;
    availablePlugins(): string;
    bfs(seed_id: string, depth: number, reverse: boolean): string;
    containsNode(node_id: string): boolean;
    degree(node_id: string, direction: string): number;
    edgeCount(): number;
    edgesFiltered(source_id?: string | null, target_id?: string | null, relation?: string | null): string;
    edgesJson(relation?: string | null): string;
    exportSnapshotJson(): string;
    hasNativeArtifact(artifact: string): boolean;
    ingestJsonl(node_type: string, data: Uint8Array, id_expr: string, data_source: string): string;
    ingestNative(plugin_name: string, artifact: string, data: Uint8Array): string;
    isPathExpression(expression: string): boolean;
    link(source: string, data: Uint8Array): string;
    loadAllPlugins(): number;
    loadPlugin(name: string): void;
    neighborIds(node_id: string, direction: string): string;
    neighborsJson(node_id: string): string;
    constructor();
    nodeCount(): number;
    nodeIds(): string;
    nodePayload(node_id: string): string | undefined;
    nodeTypes(): string;
    nodesJson(type_filter?: string | null): string;
    queryDsl(expression: string, limit: number, offset: number): string;
    queryNodeIds(expression: string, limit: number, offset: number): string;
    registerSchema(node_type: string, json_schema: string, value_field?: string | null): void;
    shortestPaths(start_id: string, end_id: string, max_depth: number): string;
    stats(): string;
    subgraphNodeIds(seed_ids_json: string, depth: number): string;
    supportedArtifacts(): string;
    toJson(): string;
    updateNodeFlags(node_id: string, add: bigint, remove: bigint, set_to: bigint): boolean;
    version(): string;
}

export function cstxTransform(source_type: string, data: Uint8Array): string;

export function flagsAllMask(): bigint;

export function flagsDefaultExcludeMask(): bigint;

export type InitInput = RequestInfo | URL | Response | BufferSource | WebAssembly.Module;

export interface InitOutput {
    readonly memory: WebAssembly.Memory;
    readonly __wbg_cstxgraph_free: (a: number, b: number) => void;
    readonly cstxTransform: (a: number, b: number, c: number, d: number, e: number) => void;
    readonly cstxgraph_addEdge: (a: number, b: number, c: number, d: number, e: number, f: number, g: number, h: number, i: number, j: number) => void;
    readonly cstxgraph_addEdgesBatch: (a: number, b: number, c: number, d: number) => void;
    readonly cstxgraph_addJoinRule: (a: number, b: number, c: number, d: number, e: number, f: number, g: number, h: number, i: number, j: number, k: number, l: number, m: number) => void;
    readonly cstxgraph_addNode: (a: number, b: number, c: number, d: number, e: number, f: number, g: number, h: number, i: number, j: number, k: bigint) => void;
    readonly cstxgraph_addNodesBatch: (a: number, b: number, c: number, d: number) => void;
    readonly cstxgraph_allNodesJson: (a: number, b: number) => void;
    readonly cstxgraph_availablePlugins: (a: number, b: number) => void;
    readonly cstxgraph_bfs: (a: number, b: number, c: number, d: number, e: number, f: number) => void;
    readonly cstxgraph_containsNode: (a: number, b: number, c: number) => number;
    readonly cstxgraph_degree: (a: number, b: number, c: number, d: number, e: number) => number;
    readonly cstxgraph_edgeCount: (a: number) => number;
    readonly cstxgraph_edgesFiltered: (a: number, b: number, c: number, d: number, e: number, f: number, g: number, h: number) => void;
    readonly cstxgraph_edgesJson: (a: number, b: number, c: number, d: number) => void;
    readonly cstxgraph_exportSnapshotJson: (a: number, b: number) => void;
    readonly cstxgraph_hasNativeArtifact: (a: number, b: number, c: number) => number;
    readonly cstxgraph_ingestJsonl: (a: number, b: number, c: number, d: number, e: number, f: number, g: number, h: number, i: number, j: number) => void;
    readonly cstxgraph_ingestNative: (a: number, b: number, c: number, d: number, e: number, f: number, g: number, h: number) => void;
    readonly cstxgraph_isPathExpression: (a: number, b: number, c: number) => number;
    readonly cstxgraph_link: (a: number, b: number, c: number, d: number, e: number, f: number) => void;
    readonly cstxgraph_loadAllPlugins: (a: number) => number;
    readonly cstxgraph_loadPlugin: (a: number, b: number, c: number, d: number) => void;
    readonly cstxgraph_neighborIds: (a: number, b: number, c: number, d: number, e: number, f: number) => void;
    readonly cstxgraph_neighborsJson: (a: number, b: number, c: number, d: number) => void;
    readonly cstxgraph_new: () => number;
    readonly cstxgraph_nodeCount: (a: number) => number;
    readonly cstxgraph_nodeIds: (a: number, b: number) => void;
    readonly cstxgraph_nodePayload: (a: number, b: number, c: number, d: number) => void;
    readonly cstxgraph_nodeTypes: (a: number, b: number) => void;
    readonly cstxgraph_nodesJson: (a: number, b: number, c: number, d: number) => void;
    readonly cstxgraph_queryDsl: (a: number, b: number, c: number, d: number, e: number, f: number) => void;
    readonly cstxgraph_queryNodeIds: (a: number, b: number, c: number, d: number, e: number, f: number) => void;
    readonly cstxgraph_registerSchema: (a: number, b: number, c: number, d: number, e: number, f: number, g: number, h: number) => void;
    readonly cstxgraph_shortestPaths: (a: number, b: number, c: number, d: number, e: number, f: number, g: number) => void;
    readonly cstxgraph_stats: (a: number, b: number) => void;
    readonly cstxgraph_subgraphNodeIds: (a: number, b: number, c: number, d: number, e: number) => void;
    readonly cstxgraph_supportedArtifacts: (a: number, b: number) => void;
    readonly cstxgraph_toJson: (a: number, b: number) => void;
    readonly cstxgraph_updateNodeFlags: (a: number, b: number, c: number, d: number, e: bigint, f: bigint, g: bigint) => void;
    readonly cstxgraph_version: (a: number, b: number) => void;
    readonly flagsAllMask: () => bigint;
    readonly flagsDefaultExcludeMask: () => bigint;
    readonly __wbindgen_export: (a: number) => void;
    readonly __wbindgen_add_to_stack_pointer: (a: number) => number;
    readonly __wbindgen_export2: (a: number, b: number) => number;
    readonly __wbindgen_export3: (a: number, b: number, c: number, d: number) => number;
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
