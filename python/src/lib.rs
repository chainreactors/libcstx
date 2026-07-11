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
}

#[pymodule]
fn _cstx_core(m: &Bound<'_, PyModule>) -> PyResult<()> {
    m.add_class::<Graph>()?;
    Ok(())
}
