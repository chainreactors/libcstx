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

The native `Graph`, `CASStore`, `Parse`, and `Transform` APIs require CGO. Build
consumers of those APIs with:

```bash
CGO_ENABLED=1 go build ./...
```

With `CGO_ENABLED=0`, the package only exposes the generated SCO/SRO data types
and their JSON helpers. Native graph, CAS, parsing, and transform functionality
is intentionally unavailable rather than replaced by a no-op fallback.

Cross-compile or target older Linux (avoid glibc version issues):

```bash
CGO_ENABLED=1 CC="zig cc -target x86_64-linux-gnu.2.17" CXX="zig c++ -target x86_64-linux-gnu.2.17" go build ./...
```

> Go + CGO 在编译时会绑定宿主系统的 glibc 版本（如 `pthread_create` 在 glibc 2.34 后从 `libpthread` 合并到 `libc.so`）。如果产物需要在较旧的 Linux 上运行，使用 [zig](https://ziglang.org/) 作为 CC 可以将 glibc 依赖锁定到指定版本。

The bundled `libcstx_ffi.a` archives statically package the Rust FFI library,
but they do not guarantee that the final Go executable is fully static. The C
toolchain may still add a dynamic loader or shared system libraries. Scratch or
other libc-free images must use explicit external static linking and verify the
resulting ELF has neither an interpreter nor `NEEDED` entries before release.

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
