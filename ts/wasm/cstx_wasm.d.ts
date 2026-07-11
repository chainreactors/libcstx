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
}
export function cstxTransform(sourceType: string, data: Uint8Array): string;
