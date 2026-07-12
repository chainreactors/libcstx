from cstxpy import Graph


def test_lifecycle():
    g = Graph()
    g.load_easm()
    assert g.node_count() == 0
    assert g.edge_count() == 0


def test_add_node():
    g = Graph()
    g.load_easm()
    g.add_node("ip", "ip:10.0.0.1", '{"ip":"10.0.0.1","country":"CN"}', ["test"])
    assert g.node_count() == 1
    assert g.contains_node("ip:10.0.0.1")


def test_ingest_native():
    g = Graph()
    g.load_easm()
    data = b'{"host":"www.example.com","a":["1.2.3.4"],"ttl":300}\n'
    result = g.ingest_native("easm", "dnsx", data)
    assert result["records_parsed"] == 1
    assert result["new_nodes"] > 0


def test_node_types():
    g = Graph()
    g.load_easm()
    g.add_node("ip", "ip:1.1.1.1", '{"ip":"1.1.1.1"}', ["test"])
    g.add_node("domain", "domain:a.com", '{"host":"a.com"}', ["test"])
    types = g.node_types()
    assert "ip" in types
    assert "domain" in types


def test_version():
    g = Graph()
    v = g.version()
    assert v and v != ""
