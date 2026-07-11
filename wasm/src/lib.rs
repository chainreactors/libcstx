use std::sync::Arc;
use wasm_bindgen::prelude::*;

use cstx_graph::engine::CstxEngine;
use cstx_graph::plugin;

#[wasm_bindgen]
pub struct CSTXGraph {
    engine: CstxEngine,
}

#[wasm_bindgen]
impl CSTXGraph {
    #[wasm_bindgen(constructor)]
    pub fn new() -> Self {
        Self { engine: CstxEngine::new() }
    }

    #[wasm_bindgen(js_name = loadAllPlugins)]
    pub fn load_all_plugins(&mut self) -> usize {
        let plugins = plugin::all_plugins();
        let count = plugins.len();
        for p in plugins {
            self.engine.register_plugin(p);
        }
        count
    }

    #[wasm_bindgen(js_name = loadPlugin)]
    pub fn load_plugin(&mut self, name: &str) -> Result<(), JsError> {
        let p = plugin::get_plugin(name)
            .ok_or_else(|| JsError::new(&format!("plugin '{}' not found", name)))?;
        self.engine.register_plugin(p);
        Ok(())
    }

    // ── Schema ──

    #[wasm_bindgen(js_name = registerSchema)]
    pub fn register_schema(&mut self, node_type: &str, json_schema: &str, value_field: Option<String>) -> Result<(), JsError> {
        self.engine.register_schema(node_type, json_schema, value_field.as_deref())
            .map_err(|e| JsError::new(&e.to_string()))
    }

    // ── Node / Edge ──

    #[wasm_bindgen(js_name = addNode)]
    pub fn add_node(&mut self, node_type: &str, cstx_id: &str, payload_json: &str, sources_json: &str, flags: u64) -> Result<(), JsError> {
        let sources: Vec<String> = serde_json::from_str(sources_json).unwrap_or_default();
        self.engine.add_node(node_type, cstx_id, payload_json, &sources, flags, None)
            .map_err(|e| JsError::new(&e.to_string()))
    }

    #[wasm_bindgen(js_name = addEdge)]
    pub fn add_edge(&mut self, source_id: &str, target_id: &str, relation: &str, data_source: &str) -> Result<(), JsError> {
        self.engine.add_edge(source_id, target_id, relation, data_source, None)
            .map_err(|e| JsError::new(&e.to_string()))
    }

    #[wasm_bindgen(js_name = addNodesBatch)]
    pub fn add_nodes_batch(&mut self, nodes_json: &str) -> Result<String, JsError> {
        let r = self.engine.add_nodes_batch(nodes_json)
            .map_err(|e| JsError::new(&e.to_string()))?;
        Ok(format!(r#"{{"new_nodes":{},"updated_nodes":{},"count":{}}}"#, r.new_nodes, r.updated_nodes, r.node_ids.len()))
    }

    #[wasm_bindgen(js_name = addEdgesBatch)]
    pub fn add_edges_batch(&mut self, edges_json: &str) -> Result<usize, JsError> {
        self.engine.add_edges_batch(edges_json)
            .map_err(|e| JsError::new(&e.to_string()))
    }

    // ── Ingest ──

    #[wasm_bindgen(js_name = ingestNative)]
    pub fn ingest_native(&mut self, plugin_name: &str, artifact: &str, data: &[u8]) -> Result<String, JsError> {
        let r = self.engine.ingest_native(plugin_name, artifact, data)
            .map_err(|e| JsError::new(&e.to_string()))?;
        Ok(format!(
            r#"{{"records_parsed":{},"new_nodes":{},"updated_nodes":{},"new_edges":{},"node_count":{},"edge_count":{}}}"#,
            r.records_parsed, r.link.new_nodes, r.link.updated_nodes, r.link.new_edges, r.node_count, r.edge_count
        ))
    }

    #[wasm_bindgen(js_name = hasNativeArtifact)]
    pub fn has_native_artifact(&self, artifact: &str) -> bool {
        self.engine.has_native_artifact(artifact)
    }

    #[wasm_bindgen(js_name = ingestJsonl)]
    pub fn ingest_jsonl(&mut self, node_type: &str, data: &[u8], id_expr: &str, data_source: &str) -> Result<String, JsError> {
        let r = self.engine.ingest_jsonl(node_type, data, id_expr, data_source)
            .map_err(|e| JsError::new(&e.to_string()))?;
        Ok(format!(r#"{{"records_parsed":{},"new_nodes":{},"updated_nodes":{}}}"#, r.records_parsed, r.new_nodes, r.updated_nodes))
    }

    // ── Queries ──

    #[wasm_bindgen(js_name = containsNode)]
    pub fn contains_node(&self, node_id: &str) -> bool { self.engine.contains_node(node_id) }

    #[wasm_bindgen(js_name = nodeCount)]
    pub fn node_count(&self) -> usize { self.engine.node_count() }

    #[wasm_bindgen(js_name = edgeCount)]
    pub fn edge_count(&self) -> usize { self.engine.edge_count() }

    #[wasm_bindgen(js_name = nodeTypes)]
    pub fn node_types(&self) -> String {
        serde_json::to_string(&self.engine.node_types()).unwrap_or_else(|_| "[]".into())
    }

    #[wasm_bindgen(js_name = nodesJson)]
    pub fn nodes_json(&self, type_filter: Option<String>) -> String {
        self.engine.nodes_json(type_filter.as_deref())
    }

    #[wasm_bindgen(js_name = edgesJson)]
    pub fn edges_json(&self, relation: Option<String>) -> String {
        self.engine.edges_json(relation.as_deref())
    }

    #[wasm_bindgen(js_name = neighborsJson)]
    pub fn neighbors_json(&self, node_id: &str) -> String {
        self.engine.neighbors_json(node_id)
    }

    #[wasm_bindgen(js_name = toJson)]
    pub fn to_json(&self) -> String { self.engine.to_json() }

    #[wasm_bindgen(js_name = allNodesJson)]
    pub fn all_nodes_json(&self) -> String { self.engine.all_nodes_json() }

    #[wasm_bindgen(js_name = nodePayload)]
    pub fn node_payload(&self, node_id: &str) -> Option<String> { self.engine.node_payload(node_id) }

    pub fn stats(&self) -> String {
        let s = self.engine.stats();
        serde_json::to_string(&serde_json::json!({"nodes": s.nodes, "edges": s.edges})).unwrap_or_else(|_| "{}".into())
    }

    pub fn version(&self) -> String { env!("CARGO_PKG_VERSION").to_string() }

    #[wasm_bindgen(js_name = supportedArtifacts)]
    pub fn supported_artifacts(&self) -> String {
        let all: Vec<String> = plugin::all_plugins()
            .iter()
            .flat_map(|p| p.artifacts().iter().map(|a| a.to_string()))
            .collect();
        serde_json::to_string(&all).unwrap_or_else(|_| "[]".into())
    }
}

// ── Stateless transform ──

#[wasm_bindgen(js_name = cstxTransform)]
pub fn transform(source_type: &str, data: &[u8]) -> Result<String, JsError> {
    let mut engine = CstxEngine::new();
    for p in plugin::all_plugins() {
        engine.register_plugin(p);
    }
    engine.transform(source_type, data)
        .map_err(|e| JsError::new(&e.to_string()))
}
