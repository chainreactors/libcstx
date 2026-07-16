# libcstx

CSTX (CyberSpace Topology eXchange) — 赛博空间拓扑表达协议。将安全工具输出转换为标准化的 SCO (Structured Cyber Observable) 对象。

## Supported Tools

`gogo`, `spray`, `zombie`, `neutron`, `subfinder`, `httpx`, `nuclei`, ...

## Supported SCO Types

`domain`, `subdomain`, `ip`, `cidr`, `port`, `app`, `url`, `framework`, `vuln`, `certificate`, `company`, `icp`, `bucket`, `endpoint`, `host`, `repository`, `secret`

## Install

### Go

```bash
go get github.com/chainreactors/libcstx/go
```

Build (requires CGO for static FFI linkage):

```bash
CGO_ENABLED=1 go build ./...
```

### Python

```bash
pip install cstxpy
```

### TypeScript

```bash
npm install @chainreactors/cstx
```

## Usage

### Go

```go
import cstx "github.com/chainreactors/libcstx/go"

// Parse tool output into typed SCO nodes
nodes, err := cstx.Parse("gogo", serviceResults)
for _, n := range nodes {
    fmt.Println(n.CstxType(), n.CstxID())
}

// Low-level: raw JSON transform
jsonBytes, err := cstx.Transform("spray", rawJSONL)
```

### Python

```python
from cstxpy import Graph, parse_sco_node

g = Graph()
result = g.link("gogo", data)
print(result.new_nodes, result.node_count)
```

### TypeScript

```typescript
import { cstxTransform } from '@chainreactors/cstx'

const nodes = await cstxTransform('gogo', data)
```

## SCO Node Interface

Every node implements:

```
CstxType() string   // e.g. "ip", "port", "vuln"
CstxID()   string   // globally unique identifier
```

Nodes are typed structs with tool-specific fields — e.g. `Ip` has `Country`, `AsnNumber`; `Vuln` has `Severity`, `VulnId`.

## License

AGPL-3.0
