# CSTX protobuf packages

The generated Python modules are the native wire model. `cstxpy.proto`
contains shared semantic runtime messages; `cstxpy.easmproto` exposes
the built-in EASM extension as one package even though its source is split
between `sco.proto` and `sro.proto`.

The C ABI accepts and returns these messages as protobuf bytes. `make proto`
generates the Python package with the official protobuf compiler from the
same sources used by Rust and Go. `GraphCursor.page()` returns the generated `GraphResultPage`
directly; iterator methods are the explicit opt-in semantic/domain adapter.
At the low-level ``cstxpy`` package a node carries its payload in one of two
spellings, exactly one of them set:
`Node.entity` is the stored `Any`, readable only with a message class this
build was generated for, and `Node.value` is the same content named by the
extension's schema document, which is the only form a type declared at runtime
can take. Which one reads return is `RuntimeConfig.payload_format`. A
relationship's `Relationship.relation` is a field-less marker — the document
names the message and the payload is empty.

The high-level ``cstx`` package installs convenience properties on these same
generated classes. There, ``node.value`` is the semantic identity string,
``node.payload`` is the raw ``EntityValue``, and ``node.type``/``node.model``/
``node.attrs``/``node.extras`` plus ``relationship.relation_type`` and
``relationship.attrs`` provide the concise Python spelling. No wrapper object
is introduced and the protobuf wire contract is unchanged.

Open-ended annotations are represented by `google.protobuf.Struct`, and
repository metadata uses that same `Struct` wire. No legacy envelope
compatibility path is provided.
