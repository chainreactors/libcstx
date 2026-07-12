use pyo3::prelude::*;
use pyo3::types::PyDict;
use std::sync::Arc;

use cstx_graph::engine::CstxEngine;

macro_rules! to_pydict {
    ($py:expr, $($key:literal => $val:expr),* $(,)?) => {{
        let dict = PyDict::new_bound($py);
        $(dict.set_item($key, $val)?;)*
        Ok(dict)
    }};
}

fn to_py_err<E: std::fmt::Display>(e: E) -> PyErr {
    pyo3::exceptions::PyRuntimeError::new_err(format!("{e:#}"))
}

#[pyclass]
struct Graph {
    engine: CstxEngine,
}

#[pymethods]
impl Graph {
    #[new]
    fn new() -> Self {
        Self { engine: CstxEngine::new() }
    }

    fn load_easm(&mut self) {
        self.engine.register_plugin(Arc::new(cstx_easm::EasmPlugin));
    }

    fn load_plugin(&mut self, name: &str) -> PyResult<()> {
        let p = cstx_graph::plugin::get_plugin(name)
            .ok_or_else(|| pyo3::exceptions::PyValueError::new_err(format!("plugin '{}' not found", name)))?;
        self.engine.register_plugin(p);
        Ok(())
    }

    fn load_all_plugins(&mut self) {
        for p in cstx_graph::plugin::all_plugins() {
            self.engine.register_plugin(p);
        }
    }

    #[pyo3(signature = (node_type, json_schema_str, value_field=None))]
    fn register_schema(&mut self, node_type: &str, json_schema_str: &str, value_field: Option<&str>) -> PyResult<()> {
        self.engine.register_schema(node_type, json_schema_str, value_field).map_err(to_py_err)
    }

    #[pyo3(signature = (left_type, right_type, relation, left_key, right_key, predicted=false, left_target_id=None, right_source_id=None))]
    fn add_join_rule(&mut self, left_type: &str, right_type: &str, relation: &str, left_key: &str, right_key: &str, predicted: bool, left_target_id: Option<&str>, right_source_id: Option<&str>) -> PyResult<()> {
        self.engine.add_join_rule(left_type, right_type, relation, left_key, right_key, predicted, left_target_id, right_source_id).map_err(to_py_err)
    }

    #[pyo3(signature = (node_type, cstx_id, payload_json, sources, cstx_flags=0, extras_json=None))]
    fn add_node(&mut self, node_type: &str, cstx_id: &str, payload_json: &str, sources: Vec<String>, cstx_flags: u64, extras_json: Option<&str>) -> PyResult<()> {
        self.engine.add_node(node_type, cstx_id, payload_json, &sources, cstx_flags, extras_json).map_err(to_py_err)
    }

    #[pyo3(signature = (source_id, target_id, relation, data_source, attrs_json=None))]
    fn add_edge(&mut self, source_id: &str, target_id: &str, relation: &str, data_source: &str, attrs_json: Option<&str>) -> PyResult<()> {
        self.engine.add_edge(source_id, target_id, relation, data_source, attrs_json).map_err(to_py_err)
    }

    // ── Transform / Link ──

    fn transform(&mut self, source_type: &str, data: &[u8]) -> PyResult<String> {
        self.engine.transform(source_type, data).map_err(to_py_err)
    }

    fn link<'py>(&mut self, py: Python<'py>, source: &str, data: &[u8]) -> PyResult<Bound<'py, PyDict>> {
        let r = self.engine.link(source, data).map_err(to_py_err)?;
        to_pydict!(py,
            "records_parsed" => r.records_parsed,
            "new_nodes" => r.link.new_nodes,
            "updated_nodes" => r.link.updated_nodes,
            "new_edges" => r.link.new_edges,
            "node_count" => r.node_count,
            "edge_count" => r.edge_count,
        )
    }

    // ── Flags ──

    #[pyo3(signature = (node_id, add=0, remove=0, set_to=None))]
    fn update_node_flags(
        &mut self,
        node_id: &str,
        add: u64,
        remove: u64,
        set_to: Option<u64>,
    ) -> PyResult<bool> {
        self.engine
            .update_node_flags(node_id, add, remove, set_to)
            .map(|r| r.is_some())
            .map_err(to_py_err)
    }

    // ── Batch ──

    fn add_nodes_batch<'py>(&mut self, py: Python<'py>, nodes_json: &str) -> PyResult<Bound<'py, PyDict>> {
        let r = self.engine.add_nodes_batch(nodes_json).map_err(to_py_err)?;
        to_pydict!(py, "new_nodes" => r.new_nodes, "updated_nodes" => r.updated_nodes, "node_ids" => r.node_ids)
    }

    fn add_edges_batch(&mut self, edges_json: &str) -> PyResult<usize> {
        self.engine.add_edges_batch(edges_json).map_err(to_py_err)
    }

    fn ingest_native<'py>(&mut self, py: Python<'py>, plugin_name: &str, artifact: &str, data: &[u8]) -> PyResult<Bound<'py, PyDict>> {
        let r = self.engine.ingest_native(plugin_name, artifact, data).map_err(to_py_err)?;
        to_pydict!(py, "records_parsed" => r.records_parsed, "new_nodes" => r.link.new_nodes, "updated_nodes" => r.link.updated_nodes, "new_edges" => r.link.new_edges, "node_count" => r.node_count, "edge_count" => r.edge_count)
    }

    fn has_native_artifact(&self, artifact: &str) -> bool {
        self.engine.has_native_artifact(artifact)
    }

    #[pyo3(signature = (node_type, data, id_expr, data_source))]
    fn ingest_jsonl<'py>(&mut self, py: Python<'py>, node_type: &str, data: &[u8], id_expr: &str, data_source: &str) -> PyResult<Bound<'py, PyDict>> {
        let r = self.engine.ingest_jsonl(node_type, data, id_expr, data_source).map_err(to_py_err)?;
        to_pydict!(py, "records_parsed" => r.records_parsed, "new_nodes" => r.new_nodes, "updated_nodes" => r.updated_nodes)
    }

    // ── Accessors ──

    fn contains_node(&self, node_id: &str) -> bool { self.engine.contains_node(node_id) }
    fn node_count(&self) -> usize { self.engine.node_count() }
    fn edge_count(&self) -> usize { self.engine.edge_count() }
    fn node_ids(&self) -> Vec<String> { self.engine.node_ids() }
    fn node_types(&self) -> Vec<String> { self.engine.node_types() }

    #[pyo3(signature = (node_type=None))]
    fn nodes_json(&self, node_type: Option<&str>) -> String { self.engine.nodes_json(node_type) }

    #[pyo3(signature = (relation=None))]
    fn edges_json(&self, relation: Option<&str>) -> String { self.engine.edges_json(relation) }

    fn all_nodes_json(&self) -> String { self.engine.all_nodes_json() }
    fn to_json(&self) -> String { self.engine.to_json() }
    fn neighbors_json(&self, node_id: &str) -> String { self.engine.neighbors_json(node_id) }

    fn node_payload(&self, node_id: &str) -> Option<String> { self.engine.node_payload(node_id) }

    fn stats<'py>(&self, py: Python<'py>) -> PyResult<Bound<'py, PyDict>> {
        let s = self.engine.stats();
        let nodes = PyDict::new_bound(py);
        for (name, count) in &s.nodes { nodes.set_item(name.as_str(), *count)?; }
        let edges = PyDict::new_bound(py);
        for (name, count) in &s.edges { edges.set_item(name.as_str(), *count)?; }
        to_pydict!(py, "nodes" => nodes, "edges" => edges)
    }

    fn version(&self) -> &'static str { env!("CARGO_PKG_VERSION") }

    fn supported_artifacts(&self) -> Vec<String> {
        cstx_graph::plugin::all_plugins()
            .iter()
            .flat_map(|p| p.artifacts().iter().map(|a| a.to_string()))
            .collect()
    }

    // ── Traversal ──

    #[pyo3(signature = (node_id, direction="out"))]
    fn neighbor_ids(&self, node_id: &str, direction: &str) -> Vec<String> {
        self.engine.neighbor_ids(node_id, direction)
    }

    #[pyo3(signature = (seed_id, depth=0, reverse=false))]
    fn bfs(&self, seed_id: &str, depth: u32, reverse: bool) -> Vec<String> {
        self.engine.bfs(seed_id, depth, reverse)
    }

    #[pyo3(signature = (start_id, end_id, max_depth=5))]
    fn shortest_paths(&self, start_id: &str, end_id: &str, max_depth: u32) -> Vec<Vec<String>> {
        self.engine.shortest_paths(start_id, end_id, max_depth)
    }

    #[pyo3(signature = (node_id, direction="both"))]
    fn degree(&self, node_id: &str, direction: &str) -> usize {
        self.engine.degree(node_id, direction)
    }

    #[pyo3(signature = (seed_ids, depth=0))]
    fn subgraph_node_ids(&self, seed_ids: Vec<String>, depth: u32) -> (Vec<String>, Vec<String>) {
        let refs: Vec<&str> = seed_ids.iter().map(|s| s.as_str()).collect();
        self.engine.subgraph_node_ids(refs, depth)
    }

    // ── Query DSL ──

    #[pyo3(signature = (expression,))]
    fn parse_query(&self, expression: &str) -> PyResult<(String, Option<String>)> {
        let pq = self.engine.parse_query(expression).map_err(to_py_err)?;
        let semantic = pq.semantic_query.clone();
        let json = serde_json::to_string(&pq).map_err(to_py_err)?;
        Ok((json, semantic))
    }

    #[pyo3(signature = (ast_json, semantic_allowed_ids=None, limit=None, offset=0, flags_exclude_mask=0, flags_include_mask=0))]
    fn execute_query(
        &self,
        ast_json: &str,
        semantic_allowed_ids: Option<Vec<String>>,
        limit: Option<usize>,
        offset: usize,
        flags_exclude_mask: u64,
        flags_include_mask: u64,
    ) -> PyResult<String> {
        self.engine
            .execute_query(
                ast_json, semantic_allowed_ids.as_deref(), limit, offset,
                flags_exclude_mask, flags_include_mask,
            )
            .map_err(to_py_err)
    }

    #[pyo3(signature = (expression, limit=None, offset=0))]
    fn query_dsl(&self, expression: &str, limit: Option<usize>, offset: usize) -> PyResult<String> {
        self.engine.query_dsl(expression, limit, offset).map_err(to_py_err)
    }

    fn is_path_expression(&self, expression: &str) -> bool {
        self.engine.is_path_expression(expression)
    }

    #[pyo3(signature = (expression, limit=None, offset=0))]
    fn query_node_ids(&self, expression: &str, limit: Option<usize>, offset: usize) -> PyResult<Vec<String>> {
        self.engine.query_node_ids(expression, limit, offset).map_err(to_py_err)
    }

    // ── Filter / Subgraph ──

    #[pyo3(signature = (source_id=None, target_id=None, relation=None))]
    fn edges_filtered(&self, source_id: Option<&str>, target_id: Option<&str>, relation: Option<&str>) -> String {
        self.engine.edges_filtered(source_id, target_id, relation)
    }

    fn export_snapshot_json(&self) -> String {
        self.engine.export_snapshot_json()
    }
}

// ── Module-level transform function ──

/// Stateless transform: creates a temporary engine with all plugins,
/// ingests data for the given source_type, returns all nodes as JSON.
#[pyfunction]
#[pyo3(name = "transform")]
fn transform_fn(source_type: &str, data: &[u8]) -> PyResult<String> {
    let mut engine = CstxEngine::new();
    for p in cstx_graph::plugin::all_plugins() {
        engine.register_plugin(p);
    }
    engine.transform(source_type, data).map_err(to_py_err)
}

#[pymodule]
fn _cstxpy(m: &Bound<'_, PyModule>) -> PyResult<()> {
    m.add_class::<Graph>()?;
    m.add_function(wrap_pyfunction!(transform_fn, m)?)?;
    Ok(())
}
