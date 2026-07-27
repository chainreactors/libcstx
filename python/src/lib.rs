use pyo3::prelude::*;
use pyo3::types::PyModule;

#[pymodule]
fn _cstxpy(module: &Bound<'_, PyModule>) -> PyResult<()> {
    cstx_core::register_python(module)
}
