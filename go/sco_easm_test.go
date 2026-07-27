package cstx

import (
	"encoding/json"
	"testing"
)

func TestParseSCONode_Domain(t *testing.T) {
	raw := `{"cstx_type":"domain","cstx_id":"domain:example.com","host":"example.com"}`
	node, err := ParseSCONode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	if node.CstxType() != "domain" {
		t.Errorf("expected type 'domain', got %q", node.CstxType())
	}
	d := node.(*DomainNode)
	if d.Host != "example.com" {
		t.Errorf("expected host 'example.com', got %q", d.Host)
	}
}

func TestParseSCONode_Port(t *testing.T) {
	raw := `{"cstx_type":"port","cstx_id":"port:1.2.3.4:80","ip":"1.2.3.4","port":"80","protocol":"tcp"}`
	node, err := ParseSCONode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	p := node.(*PortNode)
	if p.Ip != "1.2.3.4" {
		t.Errorf("ip: got %q", p.Ip)
	}
	if p.Port != "80" {
		t.Errorf("port: got %q", p.Port)
	}
	if p.Protocol != "tcp" {
		t.Errorf("protocol: got %q", p.Protocol)
	}
}

func TestParseSCONode_Ip(t *testing.T) {
	raw := `{"cstx_type":"ip","cstx_id":"ip:10.0.0.1","ip":"10.0.0.1","country":"CN","cdn":true}`
	node, err := ParseSCONode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	ip := node.(*IpNode)
	if ip.Ip != "10.0.0.1" {
		t.Errorf("ip: got %q", ip.Ip)
	}
	if ip.Country != "CN" {
		t.Errorf("country: got %q", ip.Country)
	}
	if !ip.Cdn {
		t.Error("cdn should be true")
	}
}

func TestParseSCONode_Subdomain(t *testing.T) {
	raw := `{"cstx_type":"subdomain","cstx_id":"subdomain:www.a.com","host":"www.a.com","a":["1.1.1.1","2.2.2.2"],"ttl":300}`
	node, err := ParseSCONode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	s := node.(*SubdomainNode)
	if s.Host != "www.a.com" {
		t.Errorf("host: got %q", s.Host)
	}
	if len(s.A) != 2 || s.A[0] != "1.1.1.1" {
		t.Errorf("a records: got %v", s.A)
	}
	if s.Ttl != 300 {
		t.Errorf("ttl: got %d", s.Ttl)
	}
}

func TestParseSCONode_UnknownType(t *testing.T) {
	raw := `{"cstx_type":"unknown_type","cstx_id":"x"}`
	node, err := ParseSCONode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if node != nil {
		t.Error("unknown type should return nil node")
	}
}

func TestParseSCONode_AllTypes(t *testing.T) {
	types := []string{
		"domain", "subdomain", "ip", "cidr", "port", "app", "url",
		"framework", "vuln", "certificate", "company", "icp",
		"bucket", "endpoint", "host", "repository", "secret",
	}
	for _, typ := range types {
		raw, _ := json.Marshal(map[string]string{"cstx_type": typ, "cstx_id": typ + ":test"})
		node, err := ParseSCONode(raw)
		if err != nil {
			t.Errorf("type %q: parse error: %v", typ, err)
			continue
		}
		if node == nil {
			t.Errorf("type %q: got nil", typ)
			continue
		}
		if node.CstxType() != typ {
			t.Errorf("type %q: CstxType() = %q", typ, node.CstxType())
		}
	}
}

func TestRelationConstants(t *testing.T) {
	expected := map[string]string{
		"RelResolve":      RelResolve,
		"RelOpen":         RelOpen,
		"RelHasSubdomain": RelHasSubdomain,
		"RelContain":      RelContain,
		"RelHosts":        RelHosts,
		"RelUses":         RelUses,
		"RelRefers":       RelRefers,
		"RelSecuredBy":    RelSecuredBy,
		"RelExploit":      RelExploit,
		"RelAffect":       RelAffect,
		"RelInvest":       RelInvest,
		"RelOwn":          RelOwn,
		"RelFiledFor":     RelFiledFor,
	}
	want := map[string]string{
		"RelResolve":      "resolve",
		"RelOpen":         "open",
		"RelHasSubdomain": "has-subdomain",
		"RelContain":      "contain",
		"RelHosts":        "hosts",
		"RelUses":         "uses",
		"RelRefers":       "refers",
		"RelSecuredBy":    "secured_by",
		"RelExploit":      "exploit",
		"RelAffect":       "affect",
		"RelInvest":       "invest",
		"RelOwn":          "own",
		"RelFiledFor":     "filed-for",
	}
	for name, got := range expected {
		if got != want[name] {
			t.Errorf("%s: want %q, got %q", name, want[name], got)
		}
	}
	if len(RelationTypes) != len(want) {
		t.Errorf("RelationTypes length: want %d, got %d", len(want), len(RelationTypes))
	}
}
