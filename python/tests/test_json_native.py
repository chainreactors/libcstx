import json

import pytest

from cstxpy import CSTX, CSTXError


SCHEMA = {"properties": {"ip": {"type": "string"}}}


def node(value: str) -> dict:
    return {
        "id": f"ip:{value}",
        "type": "ip",
        "value": value,
        "model": {"ip": value},
        "sources": ["json"],
        "extras": {},
    }


def db() -> CSTX:
    value = CSTX()
    value.schemas.register("ip", SCHEMA, "ip")
    return value


def test_json_batch_and_direct_nodes_match_native_cursor():
    value = db()
    nodes = [node("1.1.1.1"), node("2.2.2.2")]
    assert value.graph._add_nodes_json(json.dumps(nodes).encode()) == 2
    assert json.loads(value.graph._nodes_json()) == list(value.graph.nodes())


def test_json_edges_query_neighbors_and_snapshot():
    value = db()
    value.graph.add_nodes([node("1.1.1.1"), node("2.2.2.2")])
    relation = {
        "id": "relationship:ip:1.1.1.1:related:ip:2.2.2.2",
        "source_id": "ip:1.1.1.1",
        "target_id": "ip:2.2.2.2",
        "relation_type": "related",
        "sources": ["json"],
        "attrs": {},
    }
    assert value.graph._add_edges_json(json.dumps([relation]).encode()) == 1
    assert json.loads(value.graph._edges_json()) == list(value.graph.edges())
    assert json.loads(value.graph._neighbors_json("ip:1.1.1.1")) == list(
        value.graph.neighbors("ip:1.1.1.1")
    )
    assert json.loads(value.graph._query_json("ip")) == list(value.graph.query("ip"))
    snapshot = value.repo.dump_json()
    restored = db()
    restored.repo.load_json(snapshot)
    assert json.loads(restored.graph._nodes_json()) == json.loads(
        value.graph._nodes_json()
    )


def test_invalid_batch_is_atomic():
    value = db()
    bad = [node("1.1.1.1"), {**node("2.2.2.2"), "model": []}]
    with pytest.raises(CSTXError):
        value.graph._add_nodes_json(json.dumps(bad).encode())
    assert value.graph.node_count() == 0
