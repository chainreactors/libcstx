export class CSTXGraph {
  constructor();
  free(): void;
  loadAllPlugins(): number;
  loadPlugin(name: string): void;
  registerSchema(nodeType: string, jsonSchema: string, valueField?: string): void;
  addNode(nodeType: string, cstxId: string, payloadJson: string, sourcesJson: string, flags: bigint): void;
  addEdge(sourceId: string, targetId: string, relation: string, dataSource: string): void;
  addNodesBatch(nodesJson: string): string;
  addEdgesBatch(edgesJson: string): number;
  ingestNative(pluginName: string, artifact: string, data: Uint8Array): string;
  hasNativeArtifact(artifact: string): boolean;
  ingestJsonl(nodeType: string, data: Uint8Array, idExpr: string, dataSource: string): string;
  containsNode(nodeId: string): boolean;
  nodeCount(): number;
  edgeCount(): number;
  nodeTypes(): string;
  nodesJson(typeFilter?: string): string;
  edgesJson(relation?: string): string;
  neighborsJson(nodeId: string): string;
  toJson(): string;
  allNodesJson(): string;
  nodePayload(nodeId: string): string | undefined;
  stats(): string;
  version(): string;
  supportedArtifacts(): string;

  // Schema/Rules
  addJoinRule(leftType: string, rightType: string, relation: string, leftKey: string, rightKey: string, predicted: boolean): void;

  // Flags
  updateNodeFlags(nodeId: string, add: bigint, remove: bigint, setTo: bigint): boolean;

  // Traversal
  neighborIds(nodeId: string, direction: string): string;
  bfs(seedId: string, depth: number, reverse: boolean): string;
  shortestPaths(startId: string, endId: string, maxDepth: number): string;
  degree(nodeId: string, direction: string): number;
  subgraphNodeIds(seedIdsJson: string, depth: number): string;

  // Query DSL
  queryDsl(expression: string, limit: number, offset: number): string;
  queryNodeIds(expression: string, limit: number, offset: number): string;
  isPathExpression(expression: string): boolean;

  // Filter/Export
  edgesFiltered(sourceId: string | undefined, targetId: string | undefined, relation: string | undefined): string;
  exportSnapshotJson(): string;
  nodeIds(): string;
  availablePlugins(): string;
  link(source: string, data: Uint8Array): string;
}
export function cstxTransform(sourceType: string, data: Uint8Array): string;
