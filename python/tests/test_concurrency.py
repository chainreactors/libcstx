from __future__ import annotations


from concurrent.futures import ThreadPoolExecutor

from cstxpy import CSTX
from cstxpy.proto import cstx_pb2 as cstx, sco_pb2 as easm
from google.protobuf.any_pb2 import Any
from google.protobuf.struct_pb2 import Struct


_SCHEMA = """{
  "schema_version": 1,
  "extension": "test",
  "nodes": {
    "ip": {
      "message": "test.Ip",
      "value_field": "ip",
      "identity": { "field": "ip" },
      "fields": [
        { "name": "ip", "number": 1, "type": "string", "repeated": false,
          "optional": false, "semantic": false, "semantic_label": "ip" }
      ]
    }
  },
  "relations": {}
}"""


def _node(value: str) -> cstx.Node:
    node = cstx.Node(id=f"ip:{value}", sources=["concurrency"])
    node.entity.CopyFrom(Any(type_url="type.googleapis.com/easm.Ip", value=(easm.Ip(ip=value)).SerializeToString()))
    return node


def _graph(nodes: list[cstx.Node]) -> bytes:
    return (cstx.Graph(nodes=nodes)).SerializeToString()


def _runtime() -> CSTX:
    runtime = CSTX()
    contract = cstx.ExtensionContract(contract_version=1)
    # Fill the definition in place: a protobuf map hands back the stored
    # message, so mutating a copy after inserting it would drop the schema.
    definition = contract.extensions["test"]
    definition.name = "test"
    definition.schema = _SCHEMA
    runtime.extensions.register((contract).SerializeToString())
    runtime.graph.add_nodes(_graph([_node(str(index)) for index in range(100)]))
    return runtime


def test_reads_run_concurrently_with_writes_without_corruption() -> None:
    runtime = _runtime()

    def read(index: int) -> tuple[int, int, bool]:
        return (
            runtime.graph.node_count(),
            runtime.graph.relationship_count(),
            runtime.graph.contains(f"ip:{index % 100}"),
        )

    with ThreadPoolExecutor(max_workers=4) as pool:
        futures = [pool.submit(read, index) for index in range(100)]
        for index in range(20):
            runtime.graph.add_nodes(_graph([_node(f"write-{index}")]))
        results = [future.result() for future in futures]

    assert all(nodes >= 100 and edges == 0 and present for nodes, edges, present in results)
    assert runtime.graph.node_count() == 120


def test_writes_are_atomic_when_called_from_multiple_threads() -> None:
    runtime = _runtime()

    def write(worker: int) -> None:
        for index in range(25):
            runtime.graph.add_nodes(_graph([_node(f"{worker}-{index}")]))

    with ThreadPoolExecutor(max_workers=4) as pool:
        list(pool.map(write, range(4)))

    assert runtime.graph.node_count() == 200
