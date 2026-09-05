from google.protobuf import descriptor_pb2 as _descriptor_pb2
from google.protobuf import any_pb2 as _any_pb2
from google.protobuf import struct_pb2 as _struct_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class PayloadFormat(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    PAYLOAD_FORMAT_ENTITY: _ClassVar[PayloadFormat]
    PAYLOAD_FORMAT_VALUE: _ClassVar[PayloadFormat]

class NodeFlag(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    NODE_FLAG_UNSPECIFIED: _ClassVar[NodeFlag]
    NODE_FLAG_HONEYPOT: _ClassVar[NodeFlag]
    NODE_FLAG_NOISE: _ClassVar[NodeFlag]
    NODE_FLAG_FALSE_POSITIVE: _ClassVar[NodeFlag]
    NODE_FLAG_MANUAL_IGNORED: _ClassVar[NodeFlag]
    NODE_FLAG_THREAT_PRESENT: _ClassVar[NodeFlag]
    NODE_FLAG_HISTORIC_VULNERABLE: _ClassVar[NodeFlag]
    NODE_FLAG_INTERNAL: _ClassVar[NodeFlag]

class ChangeOperation(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    CHANGE_OPERATION_UNSPECIFIED: _ClassVar[ChangeOperation]
    CHANGE_OPERATION_ADDED: _ClassVar[ChangeOperation]
    CHANGE_OPERATION_UPDATED: _ClassVar[ChangeOperation]
    CHANGE_OPERATION_REMOVED: _ClassVar[ChangeOperation]

class SortOrder(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    SORT_ORDER_UNSPECIFIED: _ClassVar[SortOrder]
    SORT_ORDER_ID_ASC: _ClassVar[SortOrder]
    SORT_ORDER_ID_DESC: _ClassVar[SortOrder]

class Direction(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    DIRECTION_UNSPECIFIED: _ClassVar[Direction]
    DIRECTION_OUT: _ClassVar[Direction]
    DIRECTION_IN: _ClassVar[Direction]
    DIRECTION_BOTH: _ClassVar[Direction]

class ParameterlessAlgorithm(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    PARAMETERLESS_ALGORITHM_UNSPECIFIED: _ClassVar[ParameterlessAlgorithm]
    PARAMETERLESS_WEAK_COMPONENTS: _ClassVar[ParameterlessAlgorithm]
    PARAMETERLESS_STRONG_COMPONENTS: _ClassVar[ParameterlessAlgorithm]
    PARAMETERLESS_CYCLE_BASIS: _ClassVar[ParameterlessAlgorithm]
    PARAMETERLESS_BRIDGES: _ClassVar[ParameterlessAlgorithm]
    PARAMETERLESS_ARTICULATION_POINTS: _ClassVar[ParameterlessAlgorithm]
    PARAMETERLESS_CORE_NUMBERS: _ClassVar[ParameterlessAlgorithm]
    PARAMETERLESS_IS_DAG: _ClassVar[ParameterlessAlgorithm]
    PARAMETERLESS_TOPOLOGICAL_ORDER: _ClassVar[ParameterlessAlgorithm]

class NodeFlagUpdateMode(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    NODE_FLAG_UPDATE_UNSPECIFIED: _ClassVar[NodeFlagUpdateMode]
    NODE_FLAG_UPDATE_MERGE: _ClassVar[NodeFlagUpdateMode]
    NODE_FLAG_UPDATE_REPLACE: _ClassVar[NodeFlagUpdateMode]

class ObjectKind(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    OBJECT_KIND_UNSPECIFIED: _ClassVar[ObjectKind]
    OBJECT_KIND_TREE: _ClassVar[ObjectKind]
    OBJECT_KIND_STAT: _ClassVar[ObjectKind]
    OBJECT_KIND_MERGE: _ClassVar[ObjectKind]
    OBJECT_KIND_DELTA: _ClassVar[ObjectKind]
    OBJECT_KIND_PREPARE: _ClassVar[ObjectKind]
    OBJECT_KIND_HISTORY: _ClassVar[ObjectKind]
    OBJECT_KIND_COMMITS: _ClassVar[ObjectKind]
    OBJECT_KIND_DIFF: _ClassVar[ObjectKind]
    OBJECT_KIND_CLOSURE: _ClassVar[ObjectKind]

class RepositoryObjectKind(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    REPOSITORY_OBJECT_KIND_UNSPECIFIED: _ClassVar[RepositoryObjectKind]
    REPOSITORY_OBJECT_KIND_TREE: _ClassVar[RepositoryObjectKind]
    REPOSITORY_OBJECT_KIND_COMMIT: _ClassVar[RepositoryObjectKind]
    REPOSITORY_OBJECT_KIND_INDEX: _ClassVar[RepositoryObjectKind]
    REPOSITORY_OBJECT_KIND_BLOB: _ClassVar[RepositoryObjectKind]

class RepositoryPlanKind(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    REPOSITORY_PLAN_UNSPECIFIED: _ClassVar[RepositoryPlanKind]
    REPOSITORY_PLAN_TREE: _ClassVar[RepositoryPlanKind]
    REPOSITORY_PLAN_STAT: _ClassVar[RepositoryPlanKind]
    REPOSITORY_PLAN_PREPARE: _ClassVar[RepositoryPlanKind]
    REPOSITORY_PLAN_COMMITS: _ClassVar[RepositoryPlanKind]
    REPOSITORY_PLAN_DELTA: _ClassVar[RepositoryPlanKind]
    REPOSITORY_PLAN_CLOSURE: _ClassVar[RepositoryPlanKind]
    REPOSITORY_PLAN_HISTORY: _ClassVar[RepositoryPlanKind]
    REPOSITORY_PLAN_MERGE: _ClassVar[RepositoryPlanKind]
    REPOSITORY_PLAN_DIFF: _ClassVar[RepositoryPlanKind]

class DiffDetail(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    DIFF_DETAIL_UNSPECIFIED: _ClassVar[DiffDetail]
    DIFF_DETAIL_ENTITIES: _ClassVar[DiffDetail]
    DIFF_DETAIL_COUNTS: _ClassVar[DiffDetail]

class RagRecordKind(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    RAG_RECORD_KIND_UNSPECIFIED: _ClassVar[RagRecordKind]
    RAG_RECORD_NODE: _ClassVar[RagRecordKind]
    RAG_RECORD_RELATIONSHIP: _ClassVar[RagRecordKind]

class RagIndexMode(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    RAG_INDEX_MODE_UNSPECIFIED: _ClassVar[RagIndexMode]
    RAG_INDEX_INCREMENTAL: _ClassVar[RagIndexMode]
    RAG_INDEX_FULL: _ClassVar[RagIndexMode]
PAYLOAD_FORMAT_ENTITY: PayloadFormat
PAYLOAD_FORMAT_VALUE: PayloadFormat
NODE_FLAG_UNSPECIFIED: NodeFlag
NODE_FLAG_HONEYPOT: NodeFlag
NODE_FLAG_NOISE: NodeFlag
NODE_FLAG_FALSE_POSITIVE: NodeFlag
NODE_FLAG_MANUAL_IGNORED: NodeFlag
NODE_FLAG_THREAT_PRESENT: NodeFlag
NODE_FLAG_HISTORIC_VULNERABLE: NodeFlag
NODE_FLAG_INTERNAL: NodeFlag
CHANGE_OPERATION_UNSPECIFIED: ChangeOperation
CHANGE_OPERATION_ADDED: ChangeOperation
CHANGE_OPERATION_UPDATED: ChangeOperation
CHANGE_OPERATION_REMOVED: ChangeOperation
SORT_ORDER_UNSPECIFIED: SortOrder
SORT_ORDER_ID_ASC: SortOrder
SORT_ORDER_ID_DESC: SortOrder
DIRECTION_UNSPECIFIED: Direction
DIRECTION_OUT: Direction
DIRECTION_IN: Direction
DIRECTION_BOTH: Direction
PARAMETERLESS_ALGORITHM_UNSPECIFIED: ParameterlessAlgorithm
PARAMETERLESS_WEAK_COMPONENTS: ParameterlessAlgorithm
PARAMETERLESS_STRONG_COMPONENTS: ParameterlessAlgorithm
PARAMETERLESS_CYCLE_BASIS: ParameterlessAlgorithm
PARAMETERLESS_BRIDGES: ParameterlessAlgorithm
PARAMETERLESS_ARTICULATION_POINTS: ParameterlessAlgorithm
PARAMETERLESS_CORE_NUMBERS: ParameterlessAlgorithm
PARAMETERLESS_IS_DAG: ParameterlessAlgorithm
PARAMETERLESS_TOPOLOGICAL_ORDER: ParameterlessAlgorithm
NODE_FLAG_UPDATE_UNSPECIFIED: NodeFlagUpdateMode
NODE_FLAG_UPDATE_MERGE: NodeFlagUpdateMode
NODE_FLAG_UPDATE_REPLACE: NodeFlagUpdateMode
OBJECT_KIND_UNSPECIFIED: ObjectKind
OBJECT_KIND_TREE: ObjectKind
OBJECT_KIND_STAT: ObjectKind
OBJECT_KIND_MERGE: ObjectKind
OBJECT_KIND_DELTA: ObjectKind
OBJECT_KIND_PREPARE: ObjectKind
OBJECT_KIND_HISTORY: ObjectKind
OBJECT_KIND_COMMITS: ObjectKind
OBJECT_KIND_DIFF: ObjectKind
OBJECT_KIND_CLOSURE: ObjectKind
REPOSITORY_OBJECT_KIND_UNSPECIFIED: RepositoryObjectKind
REPOSITORY_OBJECT_KIND_TREE: RepositoryObjectKind
REPOSITORY_OBJECT_KIND_COMMIT: RepositoryObjectKind
REPOSITORY_OBJECT_KIND_INDEX: RepositoryObjectKind
REPOSITORY_OBJECT_KIND_BLOB: RepositoryObjectKind
REPOSITORY_PLAN_UNSPECIFIED: RepositoryPlanKind
REPOSITORY_PLAN_TREE: RepositoryPlanKind
REPOSITORY_PLAN_STAT: RepositoryPlanKind
REPOSITORY_PLAN_PREPARE: RepositoryPlanKind
REPOSITORY_PLAN_COMMITS: RepositoryPlanKind
REPOSITORY_PLAN_DELTA: RepositoryPlanKind
REPOSITORY_PLAN_CLOSURE: RepositoryPlanKind
REPOSITORY_PLAN_HISTORY: RepositoryPlanKind
REPOSITORY_PLAN_MERGE: RepositoryPlanKind
REPOSITORY_PLAN_DIFF: RepositoryPlanKind
DIFF_DETAIL_UNSPECIFIED: DiffDetail
DIFF_DETAIL_ENTITIES: DiffDetail
DIFF_DETAIL_COUNTS: DiffDetail
RAG_RECORD_KIND_UNSPECIFIED: RagRecordKind
RAG_RECORD_NODE: RagRecordKind
RAG_RECORD_RELATIONSHIP: RagRecordKind
RAG_INDEX_MODE_UNSPECIFIED: RagIndexMode
RAG_INDEX_INCREMENTAL: RagIndexMode
RAG_INDEX_FULL: RagIndexMode
CSTX_NODE_FIELD_NUMBER: _ClassVar[int]
cstx_node: _descriptor.FieldDescriptor
CSTX_RELATIONSHIP_FIELD_NUMBER: _ClassVar[int]
cstx_relationship: _descriptor.FieldDescriptor
CSTX_FIELD_FIELD_NUMBER: _ClassVar[int]
cstx_field: _descriptor.FieldDescriptor
CSTX_FLAG_FIELD_NUMBER: _ClassVar[int]
cstx_flag: _descriptor.FieldDescriptor

class CstxNodeOptions(_message.Message):
    __slots__ = ()
    NODE_TYPE_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_FIELD_NUMBER: _ClassVar[int]
    IDENTITY_COMPUTED_FIELD_NUMBER: _ClassVar[int]
    LABEL_FIELD_FIELD_NUMBER: _ClassVar[int]
    node_type: str
    value_field: str
    identity_computed: bool
    label_field: str
    def __init__(self, node_type: _Optional[str] = ..., value_field: _Optional[str] = ..., identity_computed: _Optional[bool] = ..., label_field: _Optional[str] = ...) -> None: ...

class CstxFieldOptions(_message.Message):
    __slots__ = ()
    IDENTITY_FIELD_NUMBER: _ClassVar[int]
    IDENTITY_FORMAT_FIELD_NUMBER: _ClassVar[int]
    SEMANTIC_FIELD_NUMBER: _ClassVar[int]
    SEMANTIC_LABEL_FIELD_NUMBER: _ClassVar[int]
    COLUMN_FIELD_NUMBER: _ClassVar[int]
    ORDERED_VALUES_FIELD_NUMBER: _ClassVar[int]
    identity: bool
    identity_format: str
    semantic: bool
    semantic_label: str
    column: str
    ordered_values: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, identity: _Optional[bool] = ..., identity_format: _Optional[str] = ..., semantic: _Optional[bool] = ..., semantic_label: _Optional[str] = ..., column: _Optional[str] = ..., ordered_values: _Optional[_Iterable[str]] = ...) -> None: ...

class CstxRelationshipOptions(_message.Message):
    __slots__ = ()
    RELATIONSHIP_TYPE_FIELD_NUMBER: _ClassVar[int]
    relationship_type: str
    def __init__(self, relationship_type: _Optional[str] = ...) -> None: ...

class CstxFlagOptions(_message.Message):
    __slots__ = ()
    BIT_FIELD_NUMBER: _ClassVar[int]
    DEFAULT_EXCLUDE_FIELD_NUMBER: _ClassVar[int]
    LABEL_FIELD_NUMBER: _ClassVar[int]
    bit: int
    default_exclude: bool
    label: str
    def __init__(self, bit: _Optional[int] = ..., default_exclude: _Optional[bool] = ..., label: _Optional[str] = ...) -> None: ...

class RuntimeConfig(_message.Message):
    __slots__ = ()
    PROJECT_ID_FIELD_NUMBER: _ClassVar[int]
    CURSOR_PAGE_SIZE_FIELD_NUMBER: _ClassVar[int]
    PAYLOAD_FORMAT_FIELD_NUMBER: _ClassVar[int]
    project_id: str
    cursor_page_size: int
    payload_format: PayloadFormat
    def __init__(self, project_id: _Optional[str] = ..., cursor_page_size: _Optional[int] = ..., payload_format: _Optional[_Union[PayloadFormat, str]] = ...) -> None: ...

class StringList(_message.Message):
    __slots__ = ()
    VALUES_FIELD_NUMBER: _ClassVar[int]
    values: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, values: _Optional[_Iterable[str]] = ...) -> None: ...

class EntityField(_message.Message):
    __slots__ = ()
    NAME_FIELD_NUMBER: _ClassVar[int]
    TEXT_FIELD_NUMBER: _ClassVar[int]
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    FLAG_FIELD_NUMBER: _ClassVar[int]
    REAL_FIELD_NUMBER: _ClassVar[int]
    LIST_FIELD_NUMBER: _ClassVar[int]
    name: str
    text: str
    number: int
    flag: bool
    real: float
    list: StringList
    def __init__(self, name: _Optional[str] = ..., text: _Optional[str] = ..., number: _Optional[int] = ..., flag: _Optional[bool] = ..., real: _Optional[float] = ..., list: _Optional[_Union[StringList, _Mapping]] = ...) -> None: ...

class EntityValue(_message.Message):
    __slots__ = ()
    NODE_TYPE_FIELD_NUMBER: _ClassVar[int]
    FIELDS_FIELD_NUMBER: _ClassVar[int]
    node_type: str
    fields: _containers.RepeatedCompositeFieldContainer[EntityField]
    def __init__(self, node_type: _Optional[str] = ..., fields: _Optional[_Iterable[_Union[EntityField, _Mapping]]] = ...) -> None: ...

class Node(_message.Message):
    __slots__ = ()
    ID_FIELD_NUMBER: _ClassVar[int]
    ENTITY_FIELD_NUMBER: _ClassVar[int]
    SOURCES_FIELD_NUMBER: _ClassVar[int]
    ANNOTATIONS_FIELD_NUMBER: _ClassVar[int]
    FLAGS_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    id: str
    entity: _any_pb2.Any
    sources: _containers.RepeatedScalarFieldContainer[str]
    annotations: _struct_pb2.Struct
    flags: _containers.RepeatedScalarFieldContainer[NodeFlag]
    value: EntityValue
    def __init__(self, id: _Optional[str] = ..., entity: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., sources: _Optional[_Iterable[str]] = ..., annotations: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., flags: _Optional[_Iterable[_Union[NodeFlag, str]]] = ..., value: _Optional[_Union[EntityValue, _Mapping]] = ...) -> None: ...

class Relationship(_message.Message):
    __slots__ = ()
    ID_FIELD_NUMBER: _ClassVar[int]
    SOURCE_ID_FIELD_NUMBER: _ClassVar[int]
    TARGET_ID_FIELD_NUMBER: _ClassVar[int]
    RELATION_FIELD_NUMBER: _ClassVar[int]
    SOURCES_FIELD_NUMBER: _ClassVar[int]
    ANNOTATIONS_FIELD_NUMBER: _ClassVar[int]
    id: str
    source_id: str
    target_id: str
    relation: _any_pb2.Any
    sources: _containers.RepeatedScalarFieldContainer[str]
    annotations: _struct_pb2.Struct
    def __init__(self, id: _Optional[str] = ..., source_id: _Optional[str] = ..., target_id: _Optional[str] = ..., relation: _Optional[_Union[_any_pb2.Any, _Mapping]] = ..., sources: _Optional[_Iterable[str]] = ..., annotations: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class Graph(_message.Message):
    __slots__ = ()
    NODES_FIELD_NUMBER: _ClassVar[int]
    RELATIONSHIPS_FIELD_NUMBER: _ClassVar[int]
    nodes: _containers.RepeatedCompositeFieldContainer[Node]
    relationships: _containers.RepeatedCompositeFieldContainer[Relationship]
    def __init__(self, nodes: _Optional[_Iterable[_Union[Node, _Mapping]]] = ..., relationships: _Optional[_Iterable[_Union[Relationship, _Mapping]]] = ...) -> None: ...

class GraphChangeSet(_message.Message):
    __slots__ = ()
    ADDED_NODE_IDS_FIELD_NUMBER: _ClassVar[int]
    UPDATED_NODE_IDS_FIELD_NUMBER: _ClassVar[int]
    REMOVED_NODE_IDS_FIELD_NUMBER: _ClassVar[int]
    ADDED_RELATIONSHIP_IDS_FIELD_NUMBER: _ClassVar[int]
    UPDATED_RELATIONSHIP_IDS_FIELD_NUMBER: _ClassVar[int]
    REMOVED_RELATIONSHIP_IDS_FIELD_NUMBER: _ClassVar[int]
    RESET_FIELD_NUMBER: _ClassVar[int]
    added_node_ids: _containers.RepeatedScalarFieldContainer[str]
    updated_node_ids: _containers.RepeatedScalarFieldContainer[str]
    removed_node_ids: _containers.RepeatedScalarFieldContainer[str]
    added_relationship_ids: _containers.RepeatedScalarFieldContainer[str]
    updated_relationship_ids: _containers.RepeatedScalarFieldContainer[str]
    removed_relationship_ids: _containers.RepeatedScalarFieldContainer[str]
    reset: bool
    def __init__(self, added_node_ids: _Optional[_Iterable[str]] = ..., updated_node_ids: _Optional[_Iterable[str]] = ..., removed_node_ids: _Optional[_Iterable[str]] = ..., added_relationship_ids: _Optional[_Iterable[str]] = ..., updated_relationship_ids: _Optional[_Iterable[str]] = ..., removed_relationship_ids: _Optional[_Iterable[str]] = ..., reset: _Optional[bool] = ...) -> None: ...

class GraphChangeSummary(_message.Message):
    __slots__ = ()
    ADDED_NODES_FIELD_NUMBER: _ClassVar[int]
    UPDATED_NODES_FIELD_NUMBER: _ClassVar[int]
    REMOVED_NODES_FIELD_NUMBER: _ClassVar[int]
    ADDED_RELATIONSHIPS_FIELD_NUMBER: _ClassVar[int]
    UPDATED_RELATIONSHIPS_FIELD_NUMBER: _ClassVar[int]
    REMOVED_RELATIONSHIPS_FIELD_NUMBER: _ClassVar[int]
    added_nodes: int
    updated_nodes: int
    removed_nodes: int
    added_relationships: int
    updated_relationships: int
    removed_relationships: int
    def __init__(self, added_nodes: _Optional[int] = ..., updated_nodes: _Optional[int] = ..., removed_nodes: _Optional[int] = ..., added_relationships: _Optional[int] = ..., updated_relationships: _Optional[int] = ..., removed_relationships: _Optional[int] = ...) -> None: ...

class GraphStats(_message.Message):
    __slots__ = ()
    class NodesByTypeEntry(_message.Message):
        __slots__ = ()
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: int
        def __init__(self, key: _Optional[str] = ..., value: _Optional[int] = ...) -> None: ...
    class RelationshipsByTypeEntry(_message.Message):
        __slots__ = ()
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: int
        def __init__(self, key: _Optional[str] = ..., value: _Optional[int] = ...) -> None: ...
    class ObjectsBySourceEntry(_message.Message):
        __slots__ = ()
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: int
        def __init__(self, key: _Optional[str] = ..., value: _Optional[int] = ...) -> None: ...
    class AnchorsByKindEntry(_message.Message):
        __slots__ = ()
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: int
        def __init__(self, key: _Optional[str] = ..., value: _Optional[int] = ...) -> None: ...
    NODES_BY_TYPE_FIELD_NUMBER: _ClassVar[int]
    RELATIONSHIPS_BY_TYPE_FIELD_NUMBER: _ClassVar[int]
    OBJECTS_BY_SOURCE_FIELD_NUMBER: _ClassVar[int]
    ANCHORS_BY_KIND_FIELD_NUMBER: _ClassVar[int]
    nodes_by_type: _containers.ScalarMap[str, int]
    relationships_by_type: _containers.ScalarMap[str, int]
    objects_by_source: _containers.ScalarMap[str, int]
    anchors_by_kind: _containers.ScalarMap[str, int]
    def __init__(self, nodes_by_type: _Optional[_Mapping[str, int]] = ..., relationships_by_type: _Optional[_Mapping[str, int]] = ..., objects_by_source: _Optional[_Mapping[str, int]] = ..., anchors_by_kind: _Optional[_Mapping[str, int]] = ...) -> None: ...

class Commit(_message.Message):
    __slots__ = ()
    ID_FIELD_NUMBER: _ClassVar[int]
    PARENTS_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    STATS_FIELD_NUMBER: _ClassVar[int]
    CREATED_AT_FIELD_NUMBER: _ClassVar[int]
    id: str
    parents: _containers.RepeatedScalarFieldContainer[str]
    message: str
    metadata: _struct_pb2.Struct
    stats: GraphChangeSummary
    created_at: int
    def __init__(self, id: _Optional[str] = ..., parents: _Optional[_Iterable[str]] = ..., message: _Optional[str] = ..., metadata: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., stats: _Optional[_Union[GraphChangeSummary, _Mapping]] = ..., created_at: _Optional[int] = ...) -> None: ...

class CommitLog(_message.Message):
    __slots__ = ()
    COMMITS_FIELD_NUMBER: _ClassVar[int]
    commits: _containers.RepeatedCompositeFieldContainer[Commit]
    def __init__(self, commits: _Optional[_Iterable[_Union[Commit, _Mapping]]] = ...) -> None: ...

class EntityChange(_message.Message):
    __slots__ = ()
    COMMIT_ID_FIELD_NUMBER: _ClassVar[int]
    ORDINAL_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    OPERATION_FIELD_NUMBER: _ClassVar[int]
    BEFORE_OBJECT_ID_FIELD_NUMBER: _ClassVar[int]
    AFTER_OBJECT_ID_FIELD_NUMBER: _ClassVar[int]
    commit_id: str
    ordinal: int
    timestamp: int
    operation: ChangeOperation
    before_object_id: str
    after_object_id: str
    def __init__(self, commit_id: _Optional[str] = ..., ordinal: _Optional[int] = ..., timestamp: _Optional[int] = ..., operation: _Optional[_Union[ChangeOperation, str]] = ..., before_object_id: _Optional[str] = ..., after_object_id: _Optional[str] = ...) -> None: ...

class EntityHistory(_message.Message):
    __slots__ = ()
    CHANGES_FIELD_NUMBER: _ClassVar[int]
    changes: _containers.RepeatedCompositeFieldContainer[EntityChange]
    def __init__(self, changes: _Optional[_Iterable[_Union[EntityChange, _Mapping]]] = ...) -> None: ...

class GraphSelection(_message.Message):
    __slots__ = ()
    NODE_IDS_FIELD_NUMBER: _ClassVar[int]
    RELATIONSHIP_IDS_FIELD_NUMBER: _ClassVar[int]
    ALL_NODES_FIELD_NUMBER: _ClassVar[int]
    node_ids: _containers.RepeatedScalarFieldContainer[str]
    relationship_ids: _containers.RepeatedScalarFieldContainer[str]
    all_nodes: bool
    def __init__(self, node_ids: _Optional[_Iterable[str]] = ..., relationship_ids: _Optional[_Iterable[str]] = ..., all_nodes: _Optional[bool] = ...) -> None: ...

class GraphDiff(_message.Message):
    __slots__ = ()
    ADDED_FIELD_NUMBER: _ClassVar[int]
    REMOVED_FIELD_NUMBER: _ClassVar[int]
    MODIFIED_FIELD_NUMBER: _ClassVar[int]
    TRUNCATED_FIELD_NUMBER: _ClassVar[int]
    STATS_FIELD_NUMBER: _ClassVar[int]
    added: GraphSelection
    removed: GraphSelection
    modified: GraphSelection
    truncated: bool
    stats: GraphChangeSummary
    def __init__(self, added: _Optional[_Union[GraphSelection, _Mapping]] = ..., removed: _Optional[_Union[GraphSelection, _Mapping]] = ..., modified: _Optional[_Union[GraphSelection, _Mapping]] = ..., truncated: _Optional[bool] = ..., stats: _Optional[_Union[GraphChangeSummary, _Mapping]] = ...) -> None: ...

class QueryWindow(_message.Message):
    __slots__ = ()
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    ORDER_FIELD_NUMBER: _ClassVar[int]
    limit: int
    page: int
    order: SortOrder
    def __init__(self, limit: _Optional[int] = ..., page: _Optional[int] = ..., order: _Optional[_Union[SortOrder, str]] = ...) -> None: ...

class NodeFilter(_message.Message):
    __slots__ = ()
    NODE_TYPES_FIELD_NUMBER: _ClassVar[int]
    NODE_IDS_FIELD_NUMBER: _ClassVar[int]
    SOURCES_FIELD_NUMBER: _ClassVar[int]
    NAME_CONTAINS_FIELD_NUMBER: _ClassVar[int]
    FLAGS_ALL_FIELD_NUMBER: _ClassVar[int]
    FLAGS_ANY_FIELD_NUMBER: _ClassVar[int]
    FLAGS_NONE_FIELD_NUMBER: _ClassVar[int]
    node_types: _containers.RepeatedScalarFieldContainer[str]
    node_ids: _containers.RepeatedScalarFieldContainer[str]
    sources: _containers.RepeatedScalarFieldContainer[str]
    name_contains: str
    flags_all: _containers.RepeatedScalarFieldContainer[NodeFlag]
    flags_any: _containers.RepeatedScalarFieldContainer[NodeFlag]
    flags_none: _containers.RepeatedScalarFieldContainer[NodeFlag]
    def __init__(self, node_types: _Optional[_Iterable[str]] = ..., node_ids: _Optional[_Iterable[str]] = ..., sources: _Optional[_Iterable[str]] = ..., name_contains: _Optional[str] = ..., flags_all: _Optional[_Iterable[_Union[NodeFlag, str]]] = ..., flags_any: _Optional[_Iterable[_Union[NodeFlag, str]]] = ..., flags_none: _Optional[_Iterable[_Union[NodeFlag, str]]] = ...) -> None: ...

class RelationshipFilter(_message.Message):
    __slots__ = ()
    SOURCE_ID_FIELD_NUMBER: _ClassVar[int]
    TARGET_ID_FIELD_NUMBER: _ClassVar[int]
    RELATIONSHIP_TYPES_FIELD_NUMBER: _ClassVar[int]
    SOURCES_FIELD_NUMBER: _ClassVar[int]
    source_id: str
    target_id: str
    relationship_types: _containers.RepeatedScalarFieldContainer[str]
    sources: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, source_id: _Optional[str] = ..., target_id: _Optional[str] = ..., relationship_types: _Optional[_Iterable[str]] = ..., sources: _Optional[_Iterable[str]] = ...) -> None: ...

class NodeQuery(_message.Message):
    __slots__ = ()
    FILTER_FIELD_NUMBER: _ClassVar[int]
    WINDOW_FIELD_NUMBER: _ClassVar[int]
    filter: NodeFilter
    window: QueryWindow
    def __init__(self, filter: _Optional[_Union[NodeFilter, _Mapping]] = ..., window: _Optional[_Union[QueryWindow, _Mapping]] = ...) -> None: ...

class RelationshipQuery(_message.Message):
    __slots__ = ()
    FILTER_FIELD_NUMBER: _ClassVar[int]
    WINDOW_FIELD_NUMBER: _ClassVar[int]
    filter: RelationshipFilter
    window: QueryWindow
    def __init__(self, filter: _Optional[_Union[RelationshipFilter, _Mapping]] = ..., window: _Optional[_Union[QueryWindow, _Mapping]] = ...) -> None: ...

class GraphProjection(_message.Message):
    __slots__ = ()
    NODE_FILTER_FIELD_NUMBER: _ClassVar[int]
    EXCLUDED_FIELD_NUMBER: _ClassVar[int]
    node_filter: NodeFilter
    excluded: GraphSelection
    def __init__(self, node_filter: _Optional[_Union[NodeFilter, _Mapping]] = ..., excluded: _Optional[_Union[GraphSelection, _Mapping]] = ...) -> None: ...

class QueryOptions(_message.Message):
    __slots__ = ()
    WINDOW_FIELD_NUMBER: _ClassVar[int]
    RESULT_FILTER_FIELD_NUMBER: _ClassVar[int]
    PROJECTION_FIELD_NUMBER: _ClassVar[int]
    window: QueryWindow
    result_filter: NodeFilter
    projection: GraphProjection
    def __init__(self, window: _Optional[_Union[QueryWindow, _Mapping]] = ..., result_filter: _Optional[_Union[NodeFilter, _Mapping]] = ..., projection: _Optional[_Union[GraphProjection, _Mapping]] = ...) -> None: ...

class NodeTypeCatalog(_message.Message):
    __slots__ = ()
    NODE_TYPES_FIELD_NUMBER: _ClassVar[int]
    SCHEMAS_FIELD_NUMBER: _ClassVar[int]
    node_types: _containers.RepeatedScalarFieldContainer[str]
    schemas: _containers.RepeatedCompositeFieldContainer[NodeType]
    def __init__(self, node_types: _Optional[_Iterable[str]] = ..., schemas: _Optional[_Iterable[_Union[NodeType, _Mapping]]] = ...) -> None: ...

class NeighborQuery(_message.Message):
    __slots__ = ()
    NODE_ID_FIELD_NUMBER: _ClassVar[int]
    DIRECTION_FIELD_NUMBER: _ClassVar[int]
    WINDOW_FIELD_NUMBER: _ClassVar[int]
    node_id: str
    direction: Direction
    window: QueryWindow
    def __init__(self, node_id: _Optional[str] = ..., direction: _Optional[_Union[Direction, str]] = ..., window: _Optional[_Union[QueryWindow, _Mapping]] = ...) -> None: ...

class GraphQuery(_message.Message):
    __slots__ = ()
    EXPRESSION_FIELD_NUMBER: _ClassVar[int]
    OPTIONS_FIELD_NUMBER: _ClassVar[int]
    expression: str
    options: QueryOptions
    def __init__(self, expression: _Optional[str] = ..., options: _Optional[_Union[QueryOptions, _Mapping]] = ...) -> None: ...

class NodeAnnotationUpdate(_message.Message):
    __slots__ = ()
    SELECTION_FIELD_NUMBER: _ClassVar[int]
    ANNOTATIONS_FIELD_NUMBER: _ClassVar[int]
    selection: GraphSelection
    annotations: _struct_pb2.Struct
    def __init__(self, selection: _Optional[_Union[GraphSelection, _Mapping]] = ..., annotations: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class NodeFlagChange(_message.Message):
    __slots__ = ()
    SELECTION_FIELD_NUMBER: _ClassVar[int]
    UPDATE_FIELD_NUMBER: _ClassVar[int]
    selection: GraphSelection
    update: NodeFlagUpdate
    def __init__(self, selection: _Optional[_Union[GraphSelection, _Mapping]] = ..., update: _Optional[_Union[NodeFlagUpdate, _Mapping]] = ...) -> None: ...

class BfsAlgorithm(_message.Message):
    __slots__ = ()
    SEED_ID_FIELD_NUMBER: _ClassVar[int]
    DEPTH_FIELD_NUMBER: _ClassVar[int]
    DIRECTION_FIELD_NUMBER: _ClassVar[int]
    MAX_VISITED_NODES_FIELD_NUMBER: _ClassVar[int]
    TIMEOUT_MS_FIELD_NUMBER: _ClassVar[int]
    seed_id: str
    depth: int
    direction: Direction
    max_visited_nodes: int
    timeout_ms: int
    def __init__(self, seed_id: _Optional[str] = ..., depth: _Optional[int] = ..., direction: _Optional[_Union[Direction, str]] = ..., max_visited_nodes: _Optional[int] = ..., timeout_ms: _Optional[int] = ...) -> None: ...

class BetweennessAlgorithm(_message.Message):
    __slots__ = ()
    INCLUDE_ENDPOINTS_FIELD_NUMBER: _ClassVar[int]
    NORMALIZED_FIELD_NUMBER: _ClassVar[int]
    TOP_K_FIELD_NUMBER: _ClassVar[int]
    include_endpoints: bool
    normalized: bool
    top_k: int
    def __init__(self, include_endpoints: _Optional[bool] = ..., normalized: _Optional[bool] = ..., top_k: _Optional[int] = ...) -> None: ...

class ClosenessAlgorithm(_message.Message):
    __slots__ = ()
    WF_IMPROVED_FIELD_NUMBER: _ClassVar[int]
    TOP_K_FIELD_NUMBER: _ClassVar[int]
    wf_improved: bool
    top_k: int
    def __init__(self, wf_improved: _Optional[bool] = ..., top_k: _Optional[int] = ...) -> None: ...

class LeidenAlgorithm(_message.Message):
    __slots__ = ()
    RESOLUTION_FIELD_NUMBER: _ClassVar[int]
    MIN_COMMUNITY_SIZE_FIELD_NUMBER: _ClassVar[int]
    TOP_K_FIELD_NUMBER: _ClassVar[int]
    resolution: float
    min_community_size: int
    top_k: int
    def __init__(self, resolution: _Optional[float] = ..., min_community_size: _Optional[int] = ..., top_k: _Optional[int] = ...) -> None: ...

class ShortestPathsAlgorithm(_message.Message):
    __slots__ = ()
    START_ID_FIELD_NUMBER: _ClassVar[int]
    END_ID_FIELD_NUMBER: _ClassVar[int]
    DIRECTION_FIELD_NUMBER: _ClassVar[int]
    MAX_DEPTH_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    MAX_VISITED_NODES_FIELD_NUMBER: _ClassVar[int]
    TIMEOUT_MS_FIELD_NUMBER: _ClassVar[int]
    start_id: str
    end_id: str
    direction: Direction
    max_depth: int
    limit: int
    max_visited_nodes: int
    timeout_ms: int
    def __init__(self, start_id: _Optional[str] = ..., end_id: _Optional[str] = ..., direction: _Optional[_Union[Direction, str]] = ..., max_depth: _Optional[int] = ..., limit: _Optional[int] = ..., max_visited_nodes: _Optional[int] = ..., timeout_ms: _Optional[int] = ...) -> None: ...

class Algorithm(_message.Message):
    __slots__ = ()
    BFS_FIELD_NUMBER: _ClassVar[int]
    PARAMETERLESS_FIELD_NUMBER: _ClassVar[int]
    BETWEENNESS_FIELD_NUMBER: _ClassVar[int]
    CLOSENESS_FIELD_NUMBER: _ClassVar[int]
    LEIDEN_FIELD_NUMBER: _ClassVar[int]
    SHORTEST_PATHS_FIELD_NUMBER: _ClassVar[int]
    bfs: BfsAlgorithm
    parameterless: ParameterlessAlgorithm
    betweenness: BetweennessAlgorithm
    closeness: ClosenessAlgorithm
    leiden: LeidenAlgorithm
    shortest_paths: ShortestPathsAlgorithm
    def __init__(self, bfs: _Optional[_Union[BfsAlgorithm, _Mapping]] = ..., parameterless: _Optional[_Union[ParameterlessAlgorithm, str]] = ..., betweenness: _Optional[_Union[BetweennessAlgorithm, _Mapping]] = ..., closeness: _Optional[_Union[ClosenessAlgorithm, _Mapping]] = ..., leiden: _Optional[_Union[LeidenAlgorithm, _Mapping]] = ..., shortest_paths: _Optional[_Union[ShortestPathsAlgorithm, _Mapping]] = ...) -> None: ...

class NodePage(_message.Message):
    __slots__ = ()
    VALUES_FIELD_NUMBER: _ClassVar[int]
    values: _containers.RepeatedCompositeFieldContainer[Node]
    def __init__(self, values: _Optional[_Iterable[_Union[Node, _Mapping]]] = ...) -> None: ...

class RelationshipPage(_message.Message):
    __slots__ = ()
    VALUES_FIELD_NUMBER: _ClassVar[int]
    values: _containers.RepeatedCompositeFieldContainer[Relationship]
    def __init__(self, values: _Optional[_Iterable[_Union[Relationship, _Mapping]]] = ...) -> None: ...

class ComponentMembership(_message.Message):
    __slots__ = ()
    NODE_ID_FIELD_NUMBER: _ClassVar[int]
    COMPONENT_ID_FIELD_NUMBER: _ClassVar[int]
    node_id: str
    component_id: int
    def __init__(self, node_id: _Optional[str] = ..., component_id: _Optional[int] = ...) -> None: ...

class ComponentMembershipPage(_message.Message):
    __slots__ = ()
    VALUES_FIELD_NUMBER: _ClassVar[int]
    values: _containers.RepeatedCompositeFieldContainer[ComponentMembership]
    def __init__(self, values: _Optional[_Iterable[_Union[ComponentMembership, _Mapping]]] = ...) -> None: ...

class NodeScore(_message.Message):
    __slots__ = ()
    NODE_ID_FIELD_NUMBER: _ClassVar[int]
    METRIC_FIELD_NUMBER: _ClassVar[int]
    SCORE_FIELD_NUMBER: _ClassVar[int]
    node_id: str
    metric: str
    score: float
    def __init__(self, node_id: _Optional[str] = ..., metric: _Optional[str] = ..., score: _Optional[float] = ...) -> None: ...

class NodeScorePage(_message.Message):
    __slots__ = ()
    VALUES_FIELD_NUMBER: _ClassVar[int]
    values: _containers.RepeatedCompositeFieldContainer[NodeScore]
    def __init__(self, values: _Optional[_Iterable[_Union[NodeScore, _Mapping]]] = ...) -> None: ...

class NodePair(_message.Message):
    __slots__ = ()
    SOURCE_ID_FIELD_NUMBER: _ClassVar[int]
    TARGET_ID_FIELD_NUMBER: _ClassVar[int]
    source_id: str
    target_id: str
    def __init__(self, source_id: _Optional[str] = ..., target_id: _Optional[str] = ...) -> None: ...

class NodePairPage(_message.Message):
    __slots__ = ()
    VALUES_FIELD_NUMBER: _ClassVar[int]
    values: _containers.RepeatedCompositeFieldContainer[NodePair]
    def __init__(self, values: _Optional[_Iterable[_Union[NodePair, _Mapping]]] = ...) -> None: ...

class NodeCycle(_message.Message):
    __slots__ = ()
    NODE_IDS_FIELD_NUMBER: _ClassVar[int]
    node_ids: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, node_ids: _Optional[_Iterable[str]] = ...) -> None: ...

class CyclePage(_message.Message):
    __slots__ = ()
    VALUES_FIELD_NUMBER: _ClassVar[int]
    values: _containers.RepeatedCompositeFieldContainer[NodeCycle]
    def __init__(self, values: _Optional[_Iterable[_Union[NodeCycle, _Mapping]]] = ...) -> None: ...

class NodePath(_message.Message):
    __slots__ = ()
    NODE_IDS_FIELD_NUMBER: _ClassVar[int]
    node_ids: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, node_ids: _Optional[_Iterable[str]] = ...) -> None: ...

class PathPage(_message.Message):
    __slots__ = ()
    VALUES_FIELD_NUMBER: _ClassVar[int]
    values: _containers.RepeatedCompositeFieldContainer[NodePath]
    def __init__(self, values: _Optional[_Iterable[_Union[NodePath, _Mapping]]] = ...) -> None: ...

class CommunityMembership(_message.Message):
    __slots__ = ()
    NODE_ID_FIELD_NUMBER: _ClassVar[int]
    COMMUNITY_ID_FIELD_NUMBER: _ClassVar[int]
    node_id: str
    community_id: int
    def __init__(self, node_id: _Optional[str] = ..., community_id: _Optional[int] = ...) -> None: ...

class CommunityMembershipPage(_message.Message):
    __slots__ = ()
    VALUES_FIELD_NUMBER: _ClassVar[int]
    values: _containers.RepeatedCompositeFieldContainer[CommunityMembership]
    def __init__(self, values: _Optional[_Iterable[_Union[CommunityMembership, _Mapping]]] = ...) -> None: ...

class QuerySummary(_message.Message):
    __slots__ = ()
    class NodesByTypeEntry(_message.Message):
        __slots__ = ()
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: int
        def __init__(self, key: _Optional[str] = ..., value: _Optional[int] = ...) -> None: ...
    NODES_BY_TYPE_FIELD_NUMBER: _ClassVar[int]
    nodes_by_type: _containers.ScalarMap[str, int]
    def __init__(self, nodes_by_type: _Optional[_Mapping[str, int]] = ...) -> None: ...

class TraversalSummary(_message.Message):
    __slots__ = ()
    ALGORITHM_FIELD_NUMBER: _ClassVar[int]
    DIRECTION_FIELD_NUMBER: _ClassVar[int]
    TRUNCATED_FIELD_NUMBER: _ClassVar[int]
    PROJECTION_FIELD_NUMBER: _ClassVar[int]
    algorithm: str
    direction: Direction
    truncated: bool
    projection: str
    def __init__(self, algorithm: _Optional[str] = ..., direction: _Optional[_Union[Direction, str]] = ..., truncated: _Optional[bool] = ..., projection: _Optional[str] = ...) -> None: ...

class ComponentSummary(_message.Message):
    __slots__ = ()
    ALGORITHM_FIELD_NUMBER: _ClassVar[int]
    COMPONENT_COUNT_FIELD_NUMBER: _ClassVar[int]
    PROJECTION_FIELD_NUMBER: _ClassVar[int]
    algorithm: str
    component_count: int
    projection: str
    def __init__(self, algorithm: _Optional[str] = ..., component_count: _Optional[int] = ..., projection: _Optional[str] = ...) -> None: ...

class ScoreSummary(_message.Message):
    __slots__ = ()
    METRIC_FIELD_NUMBER: _ClassVar[int]
    INCLUDE_ENDPOINTS_FIELD_NUMBER: _ClassVar[int]
    NORMALIZED_FIELD_NUMBER: _ClassVar[int]
    WF_IMPROVED_FIELD_NUMBER: _ClassVar[int]
    TOP_K_FIELD_NUMBER: _ClassVar[int]
    PROJECTION_FIELD_NUMBER: _ClassVar[int]
    metric: str
    include_endpoints: bool
    normalized: bool
    wf_improved: bool
    top_k: int
    projection: str
    def __init__(self, metric: _Optional[str] = ..., include_endpoints: _Optional[bool] = ..., normalized: _Optional[bool] = ..., wf_improved: _Optional[bool] = ..., top_k: _Optional[int] = ..., projection: _Optional[str] = ...) -> None: ...

class CommunitySummary(_message.Message):
    __slots__ = ()
    class CommunitySizesEntry(_message.Message):
        __slots__ = ()
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: int
        value: int
        def __init__(self, key: _Optional[int] = ..., value: _Optional[int] = ...) -> None: ...
    NUM_COMMUNITIES_FIELD_NUMBER: _ClassVar[int]
    TOTAL_COMMUNITIES_FIELD_NUMBER: _ClassVar[int]
    COMMUNITIES_TRUNCATED_FIELD_NUMBER: _ClassVar[int]
    MODULARITY_FIELD_NUMBER: _ClassVar[int]
    RESOLUTION_FIELD_NUMBER: _ClassVar[int]
    MIN_COMMUNITY_SIZE_FIELD_NUMBER: _ClassVar[int]
    TOP_K_FIELD_NUMBER: _ClassVar[int]
    COMMUNITY_SIZES_FIELD_NUMBER: _ClassVar[int]
    PROJECTION_FIELD_NUMBER: _ClassVar[int]
    ALGORITHM_FIELD_NUMBER: _ClassVar[int]
    num_communities: int
    total_communities: int
    communities_truncated: bool
    modularity: float
    resolution: float
    min_community_size: int
    top_k: int
    community_sizes: _containers.ScalarMap[int, int]
    projection: str
    algorithm: str
    def __init__(self, num_communities: _Optional[int] = ..., total_communities: _Optional[int] = ..., communities_truncated: _Optional[bool] = ..., modularity: _Optional[float] = ..., resolution: _Optional[float] = ..., min_community_size: _Optional[int] = ..., top_k: _Optional[int] = ..., community_sizes: _Optional[_Mapping[int, int]] = ..., projection: _Optional[str] = ..., algorithm: _Optional[str] = ...) -> None: ...

class PathSummary(_message.Message):
    __slots__ = ()
    ALGORITHM_FIELD_NUMBER: _ClassVar[int]
    START_ID_FIELD_NUMBER: _ClassVar[int]
    END_ID_FIELD_NUMBER: _ClassVar[int]
    DIRECTION_FIELD_NUMBER: _ClassVar[int]
    MAX_DEPTH_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    algorithm: str
    start_id: str
    end_id: str
    direction: Direction
    max_depth: int
    limit: int
    def __init__(self, algorithm: _Optional[str] = ..., start_id: _Optional[str] = ..., end_id: _Optional[str] = ..., direction: _Optional[_Union[Direction, str]] = ..., max_depth: _Optional[int] = ..., limit: _Optional[int] = ...) -> None: ...

class GraphResultPage(_message.Message):
    __slots__ = ()
    PAGE_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    HAS_NEXT_FIELD_NUMBER: _ClassVar[int]
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    NODES_FIELD_NUMBER: _ClassVar[int]
    RELATIONSHIPS_FIELD_NUMBER: _ClassVar[int]
    COMPONENTS_FIELD_NUMBER: _ClassVar[int]
    SCORES_FIELD_NUMBER: _ClassVar[int]
    PAIRS_FIELD_NUMBER: _ClassVar[int]
    CYCLES_FIELD_NUMBER: _ClassVar[int]
    PATHS_FIELD_NUMBER: _ClassVar[int]
    COMMUNITIES_FIELD_NUMBER: _ClassVar[int]
    QUERY_FIELD_NUMBER: _ClassVar[int]
    TRAVERSAL_FIELD_NUMBER: _ClassVar[int]
    COMPONENT_FIELD_NUMBER: _ClassVar[int]
    SCORE_FIELD_NUMBER: _ClassVar[int]
    COMMUNITY_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    page: int
    limit: int
    has_next: bool
    total: int
    nodes: NodePage
    relationships: RelationshipPage
    components: ComponentMembershipPage
    scores: NodeScorePage
    pairs: NodePairPage
    cycles: CyclePage
    paths: PathPage
    communities: CommunityMembershipPage
    query: QuerySummary
    traversal: TraversalSummary
    component: ComponentSummary
    score: ScoreSummary
    community: CommunitySummary
    path: PathSummary
    def __init__(self, page: _Optional[int] = ..., limit: _Optional[int] = ..., has_next: _Optional[bool] = ..., total: _Optional[int] = ..., nodes: _Optional[_Union[NodePage, _Mapping]] = ..., relationships: _Optional[_Union[RelationshipPage, _Mapping]] = ..., components: _Optional[_Union[ComponentMembershipPage, _Mapping]] = ..., scores: _Optional[_Union[NodeScorePage, _Mapping]] = ..., pairs: _Optional[_Union[NodePairPage, _Mapping]] = ..., cycles: _Optional[_Union[CyclePage, _Mapping]] = ..., paths: _Optional[_Union[PathPage, _Mapping]] = ..., communities: _Optional[_Union[CommunityMembershipPage, _Mapping]] = ..., query: _Optional[_Union[QuerySummary, _Mapping]] = ..., traversal: _Optional[_Union[TraversalSummary, _Mapping]] = ..., component: _Optional[_Union[ComponentSummary, _Mapping]] = ..., score: _Optional[_Union[ScoreSummary, _Mapping]] = ..., community: _Optional[_Union[CommunitySummary, _Mapping]] = ..., path: _Optional[_Union[PathSummary, _Mapping]] = ...) -> None: ...

class ParserPayload(_message.Message):
    __slots__ = ()
    PLUGIN_FIELD_NUMBER: _ClassVar[int]
    ARTIFACT_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    CONTENT_TYPE_FIELD_NUMBER: _ClassVar[int]
    plugin: str
    artifact: str
    data: bytes
    content_type: str
    def __init__(self, plugin: _Optional[str] = ..., artifact: _Optional[str] = ..., data: _Optional[bytes] = ..., content_type: _Optional[str] = ...) -> None: ...

class GraphIngestResult(_message.Message):
    __slots__ = ()
    class NodesByTypeEntry(_message.Message):
        __slots__ = ()
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: int
        def __init__(self, key: _Optional[str] = ..., value: _Optional[int] = ...) -> None: ...
    RECORDS_PARSED_FIELD_NUMBER: _ClassVar[int]
    NEW_NODES_FIELD_NUMBER: _ClassVar[int]
    UPDATED_NODES_FIELD_NUMBER: _ClassVar[int]
    NEW_RELATIONSHIPS_FIELD_NUMBER: _ClassVar[int]
    NODE_IDS_FIELD_NUMBER: _ClassVar[int]
    NODE_COUNT_FIELD_NUMBER: _ClassVar[int]
    RELATIONSHIP_COUNT_FIELD_NUMBER: _ClassVar[int]
    NODES_BY_TYPE_FIELD_NUMBER: _ClassVar[int]
    records_parsed: int
    new_nodes: int
    updated_nodes: int
    new_relationships: int
    node_ids: _containers.RepeatedScalarFieldContainer[str]
    node_count: int
    relationship_count: int
    nodes_by_type: _containers.ScalarMap[str, int]
    def __init__(self, records_parsed: _Optional[int] = ..., new_nodes: _Optional[int] = ..., updated_nodes: _Optional[int] = ..., new_relationships: _Optional[int] = ..., node_ids: _Optional[_Iterable[str]] = ..., node_count: _Optional[int] = ..., relationship_count: _Optional[int] = ..., nodes_by_type: _Optional[_Mapping[str, int]] = ...) -> None: ...

class GraphLinkResult(_message.Message):
    __slots__ = ()
    NEW_NODES_FIELD_NUMBER: _ClassVar[int]
    UPDATED_NODES_FIELD_NUMBER: _ClassVar[int]
    NEW_RELATIONSHIPS_FIELD_NUMBER: _ClassVar[int]
    RELATIONSHIP_IDS_FIELD_NUMBER: _ClassVar[int]
    new_nodes: int
    updated_nodes: int
    new_relationships: int
    relationship_ids: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, new_nodes: _Optional[int] = ..., updated_nodes: _Optional[int] = ..., new_relationships: _Optional[int] = ..., relationship_ids: _Optional[_Iterable[str]] = ...) -> None: ...

class GraphAnchor(_message.Message):
    __slots__ = ()
    CONCEPT_FIELD_NUMBER: _ClassVar[int]
    ANCHOR_ID_FIELD_NUMBER: _ClassVar[int]
    ANCHOR_TYPE_FIELD_NUMBER: _ClassVar[int]
    SOURCE_ID_FIELD_NUMBER: _ClassVar[int]
    TARGET_ID_FIELD_NUMBER: _ClassVar[int]
    INBOUND_RELATIONSHIP_ID_FIELD_NUMBER: _ClassVar[int]
    OUTBOUND_RELATIONSHIP_ID_FIELD_NUMBER: _ClassVar[int]
    concept: str
    anchor_id: str
    anchor_type: str
    source_id: str
    target_id: str
    inbound_relationship_id: str
    outbound_relationship_id: str
    def __init__(self, concept: _Optional[str] = ..., anchor_id: _Optional[str] = ..., anchor_type: _Optional[str] = ..., source_id: _Optional[str] = ..., target_id: _Optional[str] = ..., inbound_relationship_id: _Optional[str] = ..., outbound_relationship_id: _Optional[str] = ...) -> None: ...

class GraphAnchorCatalog(_message.Message):
    __slots__ = ()
    ANCHORS_FIELD_NUMBER: _ClassVar[int]
    anchors: _containers.RepeatedCompositeFieldContainer[GraphAnchor]
    def __init__(self, anchors: _Optional[_Iterable[_Union[GraphAnchor, _Mapping]]] = ...) -> None: ...

class NodeFlagUpdate(_message.Message):
    __slots__ = ()
    MODE_FIELD_NUMBER: _ClassVar[int]
    ADD_FIELD_NUMBER: _ClassVar[int]
    REMOVE_FIELD_NUMBER: _ClassVar[int]
    REPLACE_FIELD_NUMBER: _ClassVar[int]
    mode: NodeFlagUpdateMode
    add: _containers.RepeatedScalarFieldContainer[NodeFlag]
    remove: _containers.RepeatedScalarFieldContainer[NodeFlag]
    replace: _containers.RepeatedScalarFieldContainer[NodeFlag]
    def __init__(self, mode: _Optional[_Union[NodeFlagUpdateMode, str]] = ..., add: _Optional[_Iterable[_Union[NodeFlag, str]]] = ..., remove: _Optional[_Iterable[_Union[NodeFlag, str]]] = ..., replace: _Optional[_Iterable[_Union[NodeFlag, str]]] = ...) -> None: ...

class GraphProjectionReport(_message.Message):
    __slots__ = ()
    class NodeExclusion(_message.Message):
        __slots__ = ()
        NODE_ID_FIELD_NUMBER: _ClassVar[int]
        REASON_FIELD_NUMBER: _ClassVar[int]
        node_id: str
        reason: str
        def __init__(self, node_id: _Optional[str] = ..., reason: _Optional[str] = ...) -> None: ...
    EXCLUDED_NODES_FIELD_NUMBER: _ClassVar[int]
    REUSED_FIELD_NUMBER: _ClassVar[int]
    excluded_nodes: _containers.RepeatedCompositeFieldContainer[GraphProjectionReport.NodeExclusion]
    reused: bool
    def __init__(self, excluded_nodes: _Optional[_Iterable[_Union[GraphProjectionReport.NodeExclusion, _Mapping]]] = ..., reused: _Optional[bool] = ...) -> None: ...

class RepositoryObject(_message.Message):
    __slots__ = ()
    ID_FIELD_NUMBER: _ClassVar[int]
    KIND_FIELD_NUMBER: _ClassVar[int]
    PAYLOAD_FIELD_NUMBER: _ClassVar[int]
    id: str
    kind: RepositoryObjectKind
    payload: bytes
    def __init__(self, id: _Optional[str] = ..., kind: _Optional[_Union[RepositoryObjectKind, str]] = ..., payload: _Optional[bytes] = ...) -> None: ...

class PublicationPlan(_message.Message):
    __slots__ = ()
    COMMIT_FIELD_NUMBER: _ClassVar[int]
    INDEX_ROOT_FIELD_NUMBER: _ClassVar[int]
    OBJECTS_FIELD_NUMBER: _ClassVar[int]
    commit: Commit
    index_root: str
    objects: _containers.RepeatedCompositeFieldContainer[RepositoryObject]
    def __init__(self, commit: _Optional[_Union[Commit, _Mapping]] = ..., index_root: _Optional[str] = ..., objects: _Optional[_Iterable[_Union[RepositoryObject, _Mapping]]] = ...) -> None: ...

class RepositoryState(_message.Message):
    __slots__ = ()
    class Object(_message.Message):
        __slots__ = ()
        ID_FIELD_NUMBER: _ClassVar[int]
        PAYLOAD_FIELD_NUMBER: _ClassVar[int]
        id: str
        payload: bytes
        def __init__(self, id: _Optional[str] = ..., payload: _Optional[bytes] = ...) -> None: ...
    class Ref(_message.Message):
        __slots__ = ()
        NAME_FIELD_NUMBER: _ClassVar[int]
        COMMIT_ID_FIELD_NUMBER: _ClassVar[int]
        name: str
        commit_id: str
        def __init__(self, name: _Optional[str] = ..., commit_id: _Optional[str] = ...) -> None: ...
    class Index(_message.Message):
        __slots__ = ()
        COMMIT_ID_FIELD_NUMBER: _ClassVar[int]
        INDEX_ROOT_FIELD_NUMBER: _ClassVar[int]
        commit_id: str
        index_root: str
        def __init__(self, commit_id: _Optional[str] = ..., index_root: _Optional[str] = ...) -> None: ...
    OBJECTS_FIELD_NUMBER: _ClassVar[int]
    REFS_FIELD_NUMBER: _ClassVar[int]
    INDEXES_FIELD_NUMBER: _ClassVar[int]
    objects: _containers.RepeatedCompositeFieldContainer[RepositoryState.Object]
    refs: _containers.RepeatedCompositeFieldContainer[RepositoryState.Ref]
    indexes: _containers.RepeatedCompositeFieldContainer[RepositoryState.Index]
    def __init__(self, objects: _Optional[_Iterable[_Union[RepositoryState.Object, _Mapping]]] = ..., refs: _Optional[_Iterable[_Union[RepositoryState.Ref, _Mapping]]] = ..., indexes: _Optional[_Iterable[_Union[RepositoryState.Index, _Mapping]]] = ...) -> None: ...

class ObjectSelection(_message.Message):
    __slots__ = ()
    OBJECT_IDS_FIELD_NUMBER: _ClassVar[int]
    object_ids: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, object_ids: _Optional[_Iterable[str]] = ...) -> None: ...

class RepositoryObjectPlan(_message.Message):
    __slots__ = ()
    KIND_FIELD_NUMBER: _ClassVar[int]
    COMMIT_ID_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    START_TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    END_TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    ENTITY_ID_FIELD_NUMBER: _ClassVar[int]
    SOURCE_ID_FIELD_NUMBER: _ClassVar[int]
    TARGET_ID_FIELD_NUMBER: _ClassVar[int]
    DETAIL_FIELD_NUMBER: _ClassVar[int]
    kind: RepositoryPlanKind
    commit_id: str
    limit: int
    start_timestamp: int
    end_timestamp: int
    entity_id: str
    source_id: str
    target_id: str
    detail: DiffDetail
    def __init__(self, kind: _Optional[_Union[RepositoryPlanKind, str]] = ..., commit_id: _Optional[str] = ..., limit: _Optional[int] = ..., start_timestamp: _Optional[int] = ..., end_timestamp: _Optional[int] = ..., entity_id: _Optional[str] = ..., source_id: _Optional[str] = ..., target_id: _Optional[str] = ..., detail: _Optional[_Union[DiffDetail, str]] = ...) -> None: ...

class RagFilter(_message.Message):
    __slots__ = ()
    NODE_TYPES_FIELD_NUMBER: _ClassVar[int]
    RELATIONSHIP_TYPES_FIELD_NUMBER: _ClassVar[int]
    EXCLUDE_FLAGS_FIELD_NUMBER: _ClassVar[int]
    INCLUDE_FLAGS_FIELD_NUMBER: _ClassVar[int]
    node_types: _containers.RepeatedScalarFieldContainer[str]
    relationship_types: _containers.RepeatedScalarFieldContainer[str]
    exclude_flags: _containers.RepeatedScalarFieldContainer[NodeFlag]
    include_flags: _containers.RepeatedScalarFieldContainer[NodeFlag]
    def __init__(self, node_types: _Optional[_Iterable[str]] = ..., relationship_types: _Optional[_Iterable[str]] = ..., exclude_flags: _Optional[_Iterable[_Union[NodeFlag, str]]] = ..., include_flags: _Optional[_Iterable[_Union[NodeFlag, str]]] = ...) -> None: ...

class RagGraphChanges(_message.Message):
    __slots__ = ()
    CHANGED_NODE_IDS_FIELD_NUMBER: _ClassVar[int]
    DELETED_NODE_IDS_FIELD_NUMBER: _ClassVar[int]
    CHANGED_RELATIONSHIP_IDS_FIELD_NUMBER: _ClassVar[int]
    DELETED_RELATIONSHIP_IDS_FIELD_NUMBER: _ClassVar[int]
    changed_node_ids: _containers.RepeatedScalarFieldContainer[str]
    deleted_node_ids: _containers.RepeatedScalarFieldContainer[str]
    changed_relationship_ids: _containers.RepeatedScalarFieldContainer[str]
    deleted_relationship_ids: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, changed_node_ids: _Optional[_Iterable[str]] = ..., deleted_node_ids: _Optional[_Iterable[str]] = ..., changed_relationship_ids: _Optional[_Iterable[str]] = ..., deleted_relationship_ids: _Optional[_Iterable[str]] = ...) -> None: ...

class RagRecord(_message.Message):
    __slots__ = ()
    ID_FIELD_NUMBER: _ClassVar[int]
    KIND_FIELD_NUMBER: _ClassVar[int]
    TEXT_FIELD_NUMBER: _ClassVar[int]
    CONTENT_HASH_FIELD_NUMBER: _ClassVar[int]
    NODE_IDS_FIELD_NUMBER: _ClassVar[int]
    RELATIONSHIP_IDS_FIELD_NUMBER: _ClassVar[int]
    NODE_TYPE_FIELD_NUMBER: _ClassVar[int]
    RELATIONSHIP_TYPE_FIELD_NUMBER: _ClassVar[int]
    id: str
    kind: RagRecordKind
    text: str
    content_hash: str
    node_ids: _containers.RepeatedScalarFieldContainer[str]
    relationship_ids: _containers.RepeatedScalarFieldContainer[str]
    node_type: str
    relationship_type: str
    def __init__(self, id: _Optional[str] = ..., kind: _Optional[_Union[RagRecordKind, str]] = ..., text: _Optional[str] = ..., content_hash: _Optional[str] = ..., node_ids: _Optional[_Iterable[str]] = ..., relationship_ids: _Optional[_Iterable[str]] = ..., node_type: _Optional[str] = ..., relationship_type: _Optional[str] = ...) -> None: ...

class RagIndexResult(_message.Message):
    __slots__ = ()
    OPERATION_ID_FIELD_NUMBER: _ClassVar[int]
    COMMIT_FIELD_NUMBER: _ClassVar[int]
    MODE_FIELD_NUMBER: _ClassVar[int]
    UPSERT_COUNT_FIELD_NUMBER: _ClassVar[int]
    DELETE_COUNT_FIELD_NUMBER: _ClassVar[int]
    operation_id: str
    commit: str
    mode: RagIndexMode
    upsert_count: int
    delete_count: int
    def __init__(self, operation_id: _Optional[str] = ..., commit: _Optional[str] = ..., mode: _Optional[_Union[RagIndexMode, str]] = ..., upsert_count: _Optional[int] = ..., delete_count: _Optional[int] = ...) -> None: ...

class RagIndexPlan(_message.Message):
    __slots__ = ()
    COMMIT_FIELD_NUMBER: _ClassVar[int]
    MODE_FIELD_NUMBER: _ClassVar[int]
    CHANGES_FIELD_NUMBER: _ClassVar[int]
    commit: str
    mode: RagIndexMode
    changes: RagGraphChanges
    def __init__(self, commit: _Optional[str] = ..., mode: _Optional[_Union[RagIndexMode, str]] = ..., changes: _Optional[_Union[RagGraphChanges, _Mapping]] = ...) -> None: ...

class RagRecordPage(_message.Message):
    __slots__ = ()
    RECORDS_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    HAS_NEXT_FIELD_NUMBER: _ClassVar[int]
    records: _containers.RepeatedCompositeFieldContainer[RagRecord]
    page: int
    limit: int
    has_next: bool
    def __init__(self, records: _Optional[_Iterable[_Union[RagRecord, _Mapping]]] = ..., page: _Optional[int] = ..., limit: _Optional[int] = ..., has_next: _Optional[bool] = ...) -> None: ...

class RecallQuery(_message.Message):
    __slots__ = ()
    ID_FIELD_NUMBER: _ClassVar[int]
    TEXT_FIELD_NUMBER: _ClassVar[int]
    KIND_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    FILTER_FIELD_NUMBER: _ClassVar[int]
    id: str
    text: str
    kind: RagRecordKind
    limit: int
    filter: RagFilter
    def __init__(self, id: _Optional[str] = ..., text: _Optional[str] = ..., kind: _Optional[_Union[RagRecordKind, str]] = ..., limit: _Optional[int] = ..., filter: _Optional[_Union[RagFilter, _Mapping]] = ...) -> None: ...

class RecallHit(_message.Message):
    __slots__ = ()
    RECORD_ID_FIELD_NUMBER: _ClassVar[int]
    RANK_FIELD_NUMBER: _ClassVar[int]
    SCORE_FIELD_NUMBER: _ClassVar[int]
    record_id: str
    rank: int
    score: float
    def __init__(self, record_id: _Optional[str] = ..., rank: _Optional[int] = ..., score: _Optional[float] = ...) -> None: ...

class ExtensionRecallResult(_message.Message):
    __slots__ = ()
    QUERY_ID_FIELD_NUMBER: _ClassVar[int]
    EXTENSION_FIELD_NUMBER: _ClassVar[int]
    HITS_FIELD_NUMBER: _ClassVar[int]
    query_id: str
    extension: str
    hits: _containers.RepeatedCompositeFieldContainer[RecallHit]
    def __init__(self, query_id: _Optional[str] = ..., extension: _Optional[str] = ..., hits: _Optional[_Iterable[_Union[RecallHit, _Mapping]]] = ...) -> None: ...

class RecallResults(_message.Message):
    __slots__ = ()
    RESULTS_FIELD_NUMBER: _ClassVar[int]
    results: _containers.RepeatedCompositeFieldContainer[ExtensionRecallResult]
    def __init__(self, results: _Optional[_Iterable[_Union[ExtensionRecallResult, _Mapping]]] = ...) -> None: ...

class RecallPlan(_message.Message):
    __slots__ = ()
    QUERIES_FIELD_NUMBER: _ClassVar[int]
    queries: _containers.RepeatedCompositeFieldContainer[RecallQuery]
    def __init__(self, queries: _Optional[_Iterable[_Union[RecallQuery, _Mapping]]] = ...) -> None: ...

class RagPolicy(_message.Message):
    __slots__ = ()
    RRF_K_FIELD_NUMBER: _ClassVar[int]
    CANDIDATE_MULTIPLIER_FIELD_NUMBER: _ClassVar[int]
    DAMPING_FIELD_NUMBER: _ClassVar[int]
    PROPAGATION_ITERATIONS_FIELD_NUMBER: _ClassVar[int]
    MAX_PATH_DEPTH_FIELD_NUMBER: _ClassVar[int]
    EPSILON_FIELD_NUMBER: _ClassVar[int]
    COMMUNITIES_FIELD_NUMBER: _ClassVar[int]
    USE_LEXICAL_FIELD_NUMBER: _ClassVar[int]
    rrf_k: float
    candidate_multiplier: int
    damping: float
    propagation_iterations: int
    max_path_depth: int
    epsilon: float
    communities: bool
    use_lexical: bool
    def __init__(self, rrf_k: _Optional[float] = ..., candidate_multiplier: _Optional[int] = ..., damping: _Optional[float] = ..., propagation_iterations: _Optional[int] = ..., max_path_depth: _Optional[int] = ..., epsilon: _Optional[float] = ..., communities: _Optional[bool] = ..., use_lexical: _Optional[bool] = ...) -> None: ...

class RagQuery(_message.Message):
    __slots__ = ()
    TEXT_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    FILTER_FIELD_NUMBER: _ClassVar[int]
    POLICY_FIELD_NUMBER: _ClassVar[int]
    CONTEXT_BUDGET_FIELD_NUMBER: _ClassVar[int]
    text: str
    limit: int
    filter: RagFilter
    policy: RagPolicy
    context_budget: int
    def __init__(self, text: _Optional[str] = ..., limit: _Optional[int] = ..., filter: _Optional[_Union[RagFilter, _Mapping]] = ..., policy: _Optional[_Union[RagPolicy, _Mapping]] = ..., context_budget: _Optional[int] = ...) -> None: ...

class RankedNode(_message.Message):
    __slots__ = ()
    NODE_ID_FIELD_NUMBER: _ClassVar[int]
    SCORE_FIELD_NUMBER: _ClassVar[int]
    DIRECT_FIELD_NUMBER: _ClassVar[int]
    PROVENANCE_FIELD_NUMBER: _ClassVar[int]
    node_id: str
    score: float
    direct: bool
    provenance: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, node_id: _Optional[str] = ..., score: _Optional[float] = ..., direct: _Optional[bool] = ..., provenance: _Optional[_Iterable[str]] = ...) -> None: ...

class RankedRelationship(_message.Message):
    __slots__ = ()
    RELATIONSHIP_ID_FIELD_NUMBER: _ClassVar[int]
    SCORE_FIELD_NUMBER: _ClassVar[int]
    DIRECT_FIELD_NUMBER: _ClassVar[int]
    PROVENANCE_FIELD_NUMBER: _ClassVar[int]
    relationship_id: str
    score: float
    direct: bool
    provenance: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, relationship_id: _Optional[str] = ..., score: _Optional[float] = ..., direct: _Optional[bool] = ..., provenance: _Optional[_Iterable[str]] = ...) -> None: ...

class RagPath(_message.Message):
    __slots__ = ()
    NODE_IDS_FIELD_NUMBER: _ClassVar[int]
    RELATIONSHIP_IDS_FIELD_NUMBER: _ClassVar[int]
    SCORE_FIELD_NUMBER: _ClassVar[int]
    node_ids: _containers.RepeatedScalarFieldContainer[str]
    relationship_ids: _containers.RepeatedScalarFieldContainer[str]
    score: float
    def __init__(self, node_ids: _Optional[_Iterable[str]] = ..., relationship_ids: _Optional[_Iterable[str]] = ..., score: _Optional[float] = ...) -> None: ...

class RagCommunityHit(_message.Message):
    __slots__ = ()
    ID_FIELD_NUMBER: _ClassVar[int]
    LEVEL_FIELD_NUMBER: _ClassVar[int]
    MEMBER_NODE_IDS_FIELD_NUMBER: _ClassVar[int]
    SCORE_FIELD_NUMBER: _ClassVar[int]
    id: str
    level: int
    member_node_ids: _containers.RepeatedScalarFieldContainer[str]
    score: float
    def __init__(self, id: _Optional[str] = ..., level: _Optional[int] = ..., member_node_ids: _Optional[_Iterable[str]] = ..., score: _Optional[float] = ...) -> None: ...

class RagContextBlock(_message.Message):
    __slots__ = ()
    TEXT_FIELD_NUMBER: _ClassVar[int]
    RECORD_IDS_FIELD_NUMBER: _ClassVar[int]
    ESTIMATED_TOKENS_FIELD_NUMBER: _ClassVar[int]
    text: str
    record_ids: _containers.RepeatedScalarFieldContainer[str]
    estimated_tokens: int
    def __init__(self, text: _Optional[str] = ..., record_ids: _Optional[_Iterable[str]] = ..., estimated_tokens: _Optional[int] = ...) -> None: ...

class EvidenceProvenance(_message.Message):
    __slots__ = ()
    RESULT_ID_FIELD_NUMBER: _ClassVar[int]
    RECORD_IDS_FIELD_NUMBER: _ClassVar[int]
    result_id: str
    record_ids: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, result_id: _Optional[str] = ..., record_ids: _Optional[_Iterable[str]] = ...) -> None: ...

class RagResult(_message.Message):
    __slots__ = ()
    COMMIT_FIELD_NUMBER: _ClassVar[int]
    NODES_FIELD_NUMBER: _ClassVar[int]
    RELATIONSHIPS_FIELD_NUMBER: _ClassVar[int]
    PATHS_FIELD_NUMBER: _ClassVar[int]
    COMMUNITIES_FIELD_NUMBER: _ClassVar[int]
    CONTEXT_FIELD_NUMBER: _ClassVar[int]
    PROVENANCE_FIELD_NUMBER: _ClassVar[int]
    DROPPED_RECORDS_FIELD_NUMBER: _ClassVar[int]
    EXTENSIONS_FIELD_NUMBER: _ClassVar[int]
    commit: str
    nodes: _containers.RepeatedCompositeFieldContainer[RankedNode]
    relationships: _containers.RepeatedCompositeFieldContainer[RankedRelationship]
    paths: _containers.RepeatedCompositeFieldContainer[RagPath]
    communities: _containers.RepeatedCompositeFieldContainer[RagCommunityHit]
    context: _containers.RepeatedCompositeFieldContainer[RagContextBlock]
    provenance: _containers.RepeatedCompositeFieldContainer[EvidenceProvenance]
    dropped_records: _containers.RepeatedScalarFieldContainer[str]
    extensions: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, commit: _Optional[str] = ..., nodes: _Optional[_Iterable[_Union[RankedNode, _Mapping]]] = ..., relationships: _Optional[_Iterable[_Union[RankedRelationship, _Mapping]]] = ..., paths: _Optional[_Iterable[_Union[RagPath, _Mapping]]] = ..., communities: _Optional[_Iterable[_Union[RagCommunityHit, _Mapping]]] = ..., context: _Optional[_Iterable[_Union[RagContextBlock, _Mapping]]] = ..., provenance: _Optional[_Iterable[_Union[EvidenceProvenance, _Mapping]]] = ..., dropped_records: _Optional[_Iterable[str]] = ..., extensions: _Optional[_Iterable[str]] = ...) -> None: ...

class ExtensionContract(_message.Message):
    __slots__ = ()
    class ExtensionsEntry(_message.Message):
        __slots__ = ()
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: ExtensionDefinition
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[ExtensionDefinition, _Mapping]] = ...) -> None: ...
    CONTRACT_VERSION_FIELD_NUMBER: _ClassVar[int]
    EXTENSIONS_FIELD_NUMBER: _ClassVar[int]
    contract_version: int
    extensions: _containers.MessageMap[str, ExtensionDefinition]
    def __init__(self, contract_version: _Optional[int] = ..., extensions: _Optional[_Mapping[str, ExtensionDefinition]] = ...) -> None: ...

class ExtensionDefinition(_message.Message):
    __slots__ = ()
    class ParsersEntry(_message.Message):
        __slots__ = ()
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: ParserType
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[ParserType, _Mapping]] = ...) -> None: ...
    NAME_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    PARSERS_FIELD_NUMBER: _ClassVar[int]
    RULES_FIELD_NUMBER: _ClassVar[int]
    SCHEMA_FIELD_NUMBER: _ClassVar[int]
    name: str
    version: str
    parsers: _containers.MessageMap[str, ParserType]
    rules: _containers.RepeatedCompositeFieldContainer[JoinRule]
    schema: str
    def __init__(self, name: _Optional[str] = ..., version: _Optional[str] = ..., parsers: _Optional[_Mapping[str, ParserType]] = ..., rules: _Optional[_Iterable[_Union[JoinRule, _Mapping]]] = ..., schema: _Optional[str] = ...) -> None: ...

class NodeType(_message.Message):
    __slots__ = ()
    TYPE_URL_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    type_url: str
    metadata: _struct_pb2.Struct
    def __init__(self, type_url: _Optional[str] = ..., metadata: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class RelationshipType(_message.Message):
    __slots__ = ()
    TYPE_URL_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    type_url: str
    metadata: _struct_pb2.Struct
    def __init__(self, type_url: _Optional[str] = ..., metadata: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class ParserType(_message.Message):
    __slots__ = ()
    ARTIFACT_FIELD_NUMBER: _ClassVar[int]
    INPUT_SCHEMA_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    artifact: str
    input_schema: _struct_pb2.Struct
    metadata: _struct_pb2.Struct
    def __init__(self, artifact: _Optional[str] = ..., input_schema: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., metadata: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class JoinRule(_message.Message):
    __slots__ = ()
    LEFT_TYPE_URL_FIELD_NUMBER: _ClassVar[int]
    RIGHT_TYPE_URL_FIELD_NUMBER: _ClassVar[int]
    RELATIONSHIP_TYPE_URL_FIELD_NUMBER: _ClassVar[int]
    LEFT_KEY_FIELD_NUMBER: _ClassVar[int]
    RIGHT_KEY_FIELD_NUMBER: _ClassVar[int]
    PREDICTED_FIELD_NUMBER: _ClassVar[int]
    LEFT_TARGET_ID_FIELD_NUMBER: _ClassVar[int]
    RIGHT_SOURCE_ID_FIELD_NUMBER: _ClassVar[int]
    left_type_url: str
    right_type_url: str
    relationship_type_url: str
    left_key: str
    right_key: str
    predicted: bool
    left_target_id: str
    right_source_id: str
    def __init__(self, left_type_url: _Optional[str] = ..., right_type_url: _Optional[str] = ..., relationship_type_url: _Optional[str] = ..., left_key: _Optional[str] = ..., right_key: _Optional[str] = ..., predicted: _Optional[bool] = ..., left_target_id: _Optional[str] = ..., right_source_id: _Optional[str] = ...) -> None: ...

class ExtensionInfo(_message.Message):
    __slots__ = ()
    NAME_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    KIND_FIELD_NUMBER: _ClassVar[int]
    ENABLED_FIELD_NUMBER: _ClassVar[int]
    ARTIFACTS_FIELD_NUMBER: _ClassVar[int]
    name: str
    version: str
    kind: str
    enabled: bool
    artifacts: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, name: _Optional[str] = ..., version: _Optional[str] = ..., kind: _Optional[str] = ..., enabled: _Optional[bool] = ..., artifacts: _Optional[_Iterable[str]] = ...) -> None: ...

class ExtensionCatalog(_message.Message):
    __slots__ = ()
    EXTENSIONS_FIELD_NUMBER: _ClassVar[int]
    extensions: _containers.RepeatedCompositeFieldContainer[ExtensionInfo]
    def __init__(self, extensions: _Optional[_Iterable[_Union[ExtensionInfo, _Mapping]]] = ...) -> None: ...

class AnchorConcept(_message.Message):
    __slots__ = ()
    NAME_FIELD_NUMBER: _ClassVar[int]
    NODE_TYPES_FIELD_NUMBER: _ClassVar[int]
    name: str
    node_types: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, name: _Optional[str] = ..., node_types: _Optional[_Iterable[str]] = ...) -> None: ...

class AnchorConceptCatalog(_message.Message):
    __slots__ = ()
    CONCEPTS_FIELD_NUMBER: _ClassVar[int]
    concepts: _containers.RepeatedCompositeFieldContainer[AnchorConcept]
    def __init__(self, concepts: _Optional[_Iterable[_Union[AnchorConcept, _Mapping]]] = ...) -> None: ...
