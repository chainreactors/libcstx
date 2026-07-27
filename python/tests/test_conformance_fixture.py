"""Cross-language checks backed by the canonical v0.3 fixture."""

import json
from pathlib import Path

import cstxpy


def test_canonical_v03_fixture_matches_python_contract() -> None:
    fixture_path = Path(__file__).resolve().parents[3] / "tests/fixtures/v03_conformance.json"
    fixture = json.loads(fixture_path.read_text(encoding="utf-8"))
    runtime = cstxpy.CSTX()
    schema = fixture["schema"]
    runtime.schemas.register(
        schema["node_type"],
        schema["json_schema"],
        schema["value_field"],
    )
    runtime.graph.add_nodes(fixture["nodes"])
    runtime.graph.add_edges(fixture["edges"])

    expected = fixture["expected"]
    assert runtime.graph.node_ids() == expected["node_ids"]
    assert runtime.graph.node_count() == expected["node_count"]
    assert runtime.graph.edge_count() == expected["edge_count"]
    assert [node["id"] for node in runtime.graph.query(fixture["query"])] == expected["node_ids"]
