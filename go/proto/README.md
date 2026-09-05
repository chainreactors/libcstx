# CSTX protobuf packages

The generated packages are the only wire model used by the native SDK:

- `cstxproto` contains shared semantic runtime messages, cursors, and
  extension contracts.
- `easmproto` contains the built-in EASM extension. `sco.proto` and
  `sro.proto` are source-file organization only; callers import one
  extension package.

`Node.entity` and `Relationship.relation` are typed `google.protobuf.Any`
messages. Extension payloads use the canonical URL
`type.googleapis.com/<fully-qualified-message-name>`. Open-ended annotations and
schema metadata use `google.protobuf.Struct`; no envelope or JSON transport is
part of this package. Repository commit metadata uses the same `Struct` wire;
it is not encoded as a standalone `google.protobuf.Value`.
