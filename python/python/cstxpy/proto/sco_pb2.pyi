import cstx_pb2 as _cstx_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

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
NODE_FLAG_UNSPECIFIED: NodeFlag
NODE_FLAG_HONEYPOT: NodeFlag
NODE_FLAG_NOISE: NodeFlag
NODE_FLAG_FALSE_POSITIVE: NodeFlag
NODE_FLAG_MANUAL_IGNORED: NodeFlag
NODE_FLAG_THREAT_PRESENT: NodeFlag
NODE_FLAG_HISTORIC_VULNERABLE: NodeFlag
NODE_FLAG_INTERNAL: NodeFlag

class Domain(_message.Message):
    __slots__ = ()
    HOST_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    host: str
    extra: str
    def __init__(self, host: _Optional[str] = ..., extra: _Optional[str] = ...) -> None: ...

class Subdomain(_message.Message):
    __slots__ = ()
    HOST_FIELD_NUMBER: _ClassVar[int]
    IS_TLD_FIELD_NUMBER: _ClassVar[int]
    TTL_FIELD_NUMBER: _ClassVar[int]
    RESOLVER_FIELD_NUMBER: _ClassVar[int]
    A_FIELD_NUMBER: _ClassVar[int]
    AAAA_FIELD_NUMBER: _ClassVar[int]
    CNAME_FIELD_NUMBER: _ClassVar[int]
    MX_FIELD_NUMBER: _ClassVar[int]
    NS_FIELD_NUMBER: _ClassVar[int]
    TXT_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    host: str
    is_tld: bool
    ttl: int
    resolver: _containers.RepeatedScalarFieldContainer[str]
    a: _containers.RepeatedScalarFieldContainer[str]
    aaaa: _containers.RepeatedScalarFieldContainer[str]
    cname: _containers.RepeatedScalarFieldContainer[str]
    mx: _containers.RepeatedScalarFieldContainer[str]
    ns: _containers.RepeatedScalarFieldContainer[str]
    txt: _containers.RepeatedScalarFieldContainer[str]
    extra: str
    def __init__(self, host: _Optional[str] = ..., is_tld: _Optional[bool] = ..., ttl: _Optional[int] = ..., resolver: _Optional[_Iterable[str]] = ..., a: _Optional[_Iterable[str]] = ..., aaaa: _Optional[_Iterable[str]] = ..., cname: _Optional[_Iterable[str]] = ..., mx: _Optional[_Iterable[str]] = ..., ns: _Optional[_Iterable[str]] = ..., txt: _Optional[_Iterable[str]] = ..., extra: _Optional[str] = ...) -> None: ...

class Ip(_message.Message):
    __slots__ = ()
    IP_FIELD_NUMBER: _ClassVar[int]
    COUNTRY_FIELD_NUMBER: _ClassVar[int]
    AREA_FIELD_NUMBER: _ClassVar[int]
    ASN_NUMBER_FIELD_NUMBER: _ClassVar[int]
    AS_NAME_FIELD_NUMBER: _ClassVar[int]
    CDN_NAME_FIELD_NUMBER: _ClassVar[int]
    CLOUD_NAME_FIELD_NUMBER: _ClassVar[int]
    WAF_NAME_FIELD_NUMBER: _ClassVar[int]
    CDN_FIELD_NUMBER: _ClassVar[int]
    CLOUD_FIELD_NUMBER: _ClassVar[int]
    WAF_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    ip: str
    country: str
    area: str
    asn_number: str
    as_name: str
    cdn_name: str
    cloud_name: str
    waf_name: str
    cdn: bool
    cloud: bool
    waf: bool
    extra: str
    def __init__(self, ip: _Optional[str] = ..., country: _Optional[str] = ..., area: _Optional[str] = ..., asn_number: _Optional[str] = ..., as_name: _Optional[str] = ..., cdn_name: _Optional[str] = ..., cloud_name: _Optional[str] = ..., waf_name: _Optional[str] = ..., cdn: _Optional[bool] = ..., cloud: _Optional[bool] = ..., waf: _Optional[bool] = ..., extra: _Optional[str] = ...) -> None: ...

class Cidr(_message.Message):
    __slots__ = ()
    CIDR_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    cidr: str
    extra: str
    def __init__(self, cidr: _Optional[str] = ..., extra: _Optional[str] = ...) -> None: ...

class Port(_message.Message):
    __slots__ = ()
    IP_FIELD_NUMBER: _ClassVar[int]
    PORT_FIELD_NUMBER: _ClassVar[int]
    PROTOCOL_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    ip: str
    port: str
    protocol: str
    extra: str
    def __init__(self, ip: _Optional[str] = ..., port: _Optional[str] = ..., protocol: _Optional[str] = ..., extra: _Optional[str] = ...) -> None: ...

class App(_message.Message):
    __slots__ = ()
    APP_ID_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    FRAMEWORKS_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    MIDWARE_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    STATUS_CODE_FIELD_NUMBER: _ClassVar[int]
    HOST_FIELD_NUMBER: _ClassVar[int]
    CONTENT_TYPE_FIELD_NUMBER: _ClassVar[int]
    BODY_LENGTH_FIELD_NUMBER: _ClassVar[int]
    HEADER_LENGTH_FIELD_NUMBER: _ClassVar[int]
    SCREENSHOT_ID_FIELD_NUMBER: _ClassVar[int]
    SCREENSHOT_PATH_FIELD_NUMBER: _ClassVar[int]
    IP_FIELD_NUMBER: _ClassVar[int]
    PORT_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    app_id: str
    url: str
    frameworks: _containers.RepeatedScalarFieldContainer[str]
    title: str
    midware: str
    status: str
    status_code: int
    host: str
    content_type: str
    body_length: int
    header_length: int
    screenshot_id: str
    screenshot_path: str
    ip: str
    port: str
    extra: str
    def __init__(self, app_id: _Optional[str] = ..., url: _Optional[str] = ..., frameworks: _Optional[_Iterable[str]] = ..., title: _Optional[str] = ..., midware: _Optional[str] = ..., status: _Optional[str] = ..., status_code: _Optional[int] = ..., host: _Optional[str] = ..., content_type: _Optional[str] = ..., body_length: _Optional[int] = ..., header_length: _Optional[int] = ..., screenshot_id: _Optional[str] = ..., screenshot_path: _Optional[str] = ..., ip: _Optional[str] = ..., port: _Optional[str] = ..., extra: _Optional[str] = ...) -> None: ...

class Url(_message.Message):
    __slots__ = ()
    URL_FIELD_NUMBER: _ClassVar[int]
    SCHEME_FIELD_NUMBER: _ClassVar[int]
    HOST_FIELD_NUMBER: _ClassVar[int]
    PORT_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    IP_FIELD_NUMBER: _ClassVar[int]
    STATUS_CODE_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    BODY_LENGTH_FIELD_NUMBER: _ClassVar[int]
    CONTENT_TYPE_FIELD_NUMBER: _ClassVar[int]
    REDIRECT_URL_FIELD_NUMBER: _ClassVar[int]
    FRAMEWORKS_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    url: str
    scheme: str
    host: str
    port: str
    path: str
    ip: str
    status_code: int
    title: str
    body_length: int
    content_type: str
    redirect_url: str
    frameworks: _containers.RepeatedScalarFieldContainer[str]
    extra: str
    def __init__(self, url: _Optional[str] = ..., scheme: _Optional[str] = ..., host: _Optional[str] = ..., port: _Optional[str] = ..., path: _Optional[str] = ..., ip: _Optional[str] = ..., status_code: _Optional[int] = ..., title: _Optional[str] = ..., body_length: _Optional[int] = ..., content_type: _Optional[str] = ..., redirect_url: _Optional[str] = ..., frameworks: _Optional[_Iterable[str]] = ..., extra: _Optional[str] = ...) -> None: ...

class Framework(_message.Message):
    __slots__ = ()
    NAME_FIELD_NUMBER: _ClassVar[int]
    PART_FIELD_NUMBER: _ClassVar[int]
    VENDOR_FIELD_NUMBER: _ClassVar[int]
    PRODUCT_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    TAGS_FIELD_NUMBER: _ClassVar[int]
    IS_FOCUS_FIELD_NUMBER: _ClassVar[int]
    SOURCES_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    name: str
    part: str
    vendor: str
    product: str
    version: str
    tags: _containers.RepeatedScalarFieldContainer[str]
    is_focus: bool
    sources: _containers.RepeatedScalarFieldContainer[str]
    extra: str
    def __init__(self, name: _Optional[str] = ..., part: _Optional[str] = ..., vendor: _Optional[str] = ..., product: _Optional[str] = ..., version: _Optional[str] = ..., tags: _Optional[_Iterable[str]] = ..., is_focus: _Optional[bool] = ..., sources: _Optional[_Iterable[str]] = ..., extra: _Optional[str] = ...) -> None: ...

class Vuln(_message.Message):
    __slots__ = ()
    VALUE_FIELD_NUMBER: _ClassVar[int]
    VULN_ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    ASSET_ID_FIELD_NUMBER: _ClassVar[int]
    SEVERITY_FIELD_NUMBER: _ClassVar[int]
    TAGS_FIELD_NUMBER: _ClassVar[int]
    IP_FIELD_NUMBER: _ClassVar[int]
    HOST_FIELD_NUMBER: _ClassVar[int]
    PORT_FIELD_NUMBER: _ClassVar[int]
    PROTOCOL_FIELD_NUMBER: _ClassVar[int]
    SCHEME_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    POCNAME_FIELD_NUMBER: _ClassVar[int]
    REQUEST_FIELD_NUMBER: _ClassVar[int]
    RESPONSE_FIELD_NUMBER: _ClassVar[int]
    USERNAME_FIELD_NUMBER: _ClassVar[int]
    PASSWORD_FIELD_NUMBER: _ClassVar[int]
    MATCHED_FIELD_NUMBER: _ClassVar[int]
    EXTRACTED_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    value: str
    vuln_id: str
    name: str
    asset_id: str
    severity: str
    tags: _containers.RepeatedScalarFieldContainer[str]
    ip: str
    host: str
    port: str
    protocol: str
    scheme: str
    url: str
    path: str
    pocname: str
    request: str
    response: str
    username: str
    password: str
    matched: bool
    extracted: bool
    extra: str
    def __init__(self, value: _Optional[str] = ..., vuln_id: _Optional[str] = ..., name: _Optional[str] = ..., asset_id: _Optional[str] = ..., severity: _Optional[str] = ..., tags: _Optional[_Iterable[str]] = ..., ip: _Optional[str] = ..., host: _Optional[str] = ..., port: _Optional[str] = ..., protocol: _Optional[str] = ..., scheme: _Optional[str] = ..., url: _Optional[str] = ..., path: _Optional[str] = ..., pocname: _Optional[str] = ..., request: _Optional[str] = ..., response: _Optional[str] = ..., username: _Optional[str] = ..., password: _Optional[str] = ..., matched: _Optional[bool] = ..., extracted: _Optional[bool] = ..., extra: _Optional[str] = ...) -> None: ...

class SarifVuln(_message.Message):
    __slots__ = ()
    VALUE_FIELD_NUMBER: _ClassVar[int]
    VULN_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    SOURCE_FIELD_NUMBER: _ClassVar[int]
    TARGET_FIELD_NUMBER: _ClassVar[int]
    TAGS_FIELD_NUMBER: _ClassVar[int]
    ASSET_CSTX_ID_FIELD_NUMBER: _ClassVar[int]
    KIND_FIELD_NUMBER: _ClassVar[int]
    LEVEL_FIELD_NUMBER: _ClassVar[int]
    BASELINE_STATE_FIELD_NUMBER: _ClassVar[int]
    RULE_ID_FIELD_NUMBER: _ClassVar[int]
    EVIDENCE_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    value: str
    vuln_id: str
    title: str
    description: str
    source: str
    target: str
    tags: _containers.RepeatedScalarFieldContainer[str]
    asset_cstx_id: str
    kind: str
    level: str
    baseline_state: str
    rule_id: str
    evidence: str
    extra: str
    def __init__(self, value: _Optional[str] = ..., vuln_id: _Optional[str] = ..., title: _Optional[str] = ..., description: _Optional[str] = ..., source: _Optional[str] = ..., target: _Optional[str] = ..., tags: _Optional[_Iterable[str]] = ..., asset_cstx_id: _Optional[str] = ..., kind: _Optional[str] = ..., level: _Optional[str] = ..., baseline_state: _Optional[str] = ..., rule_id: _Optional[str] = ..., evidence: _Optional[str] = ..., extra: _Optional[str] = ...) -> None: ...

class Certificate(_message.Message):
    __slots__ = ()
    FINGERPRINT_FIELD_NUMBER: _ClassVar[int]
    SERIAL_FIELD_NUMBER: _ClassVar[int]
    ISSUER_FIELD_NUMBER: _ClassVar[int]
    SUBJECT_FIELD_NUMBER: _ClassVar[int]
    NOT_BEFORE_FIELD_NUMBER: _ClassVar[int]
    NOT_AFTER_FIELD_NUMBER: _ClassVar[int]
    SAN_FIELD_NUMBER: _ClassVar[int]
    HOST_FIELD_NUMBER: _ClassVar[int]
    IP_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    fingerprint: str
    serial: str
    issuer: str
    subject: str
    not_before: str
    not_after: str
    san: _containers.RepeatedScalarFieldContainer[str]
    host: str
    ip: str
    extra: str
    def __init__(self, fingerprint: _Optional[str] = ..., serial: _Optional[str] = ..., issuer: _Optional[str] = ..., subject: _Optional[str] = ..., not_before: _Optional[str] = ..., not_after: _Optional[str] = ..., san: _Optional[_Iterable[str]] = ..., host: _Optional[str] = ..., ip: _Optional[str] = ..., extra: _Optional[str] = ...) -> None: ...

class Company(_message.Message):
    __slots__ = ()
    NAME_FIELD_NUMBER: _ClassVar[int]
    PERC_FIELD_NUMBER: _ClassVar[int]
    TYCID_FIELD_NUMBER: _ClassVar[int]
    ICP_FIELD_NUMBER: _ClassVar[int]
    PARENT_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    name: str
    perc: str
    tycid: str
    icp: str
    parent: str
    extra: str
    def __init__(self, name: _Optional[str] = ..., perc: _Optional[str] = ..., tycid: _Optional[str] = ..., icp: _Optional[str] = ..., parent: _Optional[str] = ..., extra: _Optional[str] = ...) -> None: ...

class Icp(_message.Message):
    __slots__ = ()
    ICP_FIELD_NUMBER: _ClassVar[int]
    SUB_FIELD_NUMBER: _ClassVar[int]
    DATE_FIELD_NUMBER: _ClassVar[int]
    COMPANY_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DOMAIN_FIELD_NUMBER: _ClassVar[int]
    IP_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    icp: str
    sub: str
    date: str
    company: str
    title: str
    domain: str
    ip: str
    extra: str
    def __init__(self, icp: _Optional[str] = ..., sub: _Optional[str] = ..., date: _Optional[str] = ..., company: _Optional[str] = ..., title: _Optional[str] = ..., domain: _Optional[str] = ..., ip: _Optional[str] = ..., extra: _Optional[str] = ...) -> None: ...

class Bucket(_message.Message):
    __slots__ = ()
    PROVIDER_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    REGION_FIELD_NUMBER: _ClassVar[int]
    ENDPOINT_FIELD_NUMBER: _ClassVar[int]
    ACL_FIELD_NUMBER: _ClassVar[int]
    OBJECT_COUNT_FIELD_NUMBER: _ClassVar[int]
    KNOWN_PATHS_FIELD_NUMBER: _ClassVar[int]
    SOURCE_URL_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    provider: str
    name: str
    region: str
    endpoint: str
    acl: str
    object_count: int
    known_paths: _containers.RepeatedScalarFieldContainer[str]
    source_url: str
    extra: str
    def __init__(self, provider: _Optional[str] = ..., name: _Optional[str] = ..., region: _Optional[str] = ..., endpoint: _Optional[str] = ..., acl: _Optional[str] = ..., object_count: _Optional[int] = ..., known_paths: _Optional[_Iterable[str]] = ..., source_url: _Optional[str] = ..., extra: _Optional[str] = ...) -> None: ...

class Endpoint(_message.Message):
    __slots__ = ()
    URL_FIELD_NUMBER: _ClassVar[int]
    METHOD_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    CONTENT_TYPE_FIELD_NUMBER: _ClassVar[int]
    STATUS_CODE_FIELD_NUMBER: _ClassVar[int]
    SOURCE_FIELD_NUMBER: _ClassVar[int]
    SOURCE_URL_FIELD_NUMBER: _ClassVar[int]
    PARAMETERS_FIELD_NUMBER: _ClassVar[int]
    TAGS_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    url: str
    method: str
    path: str
    content_type: str
    status_code: int
    source: str
    source_url: str
    parameters: _containers.RepeatedScalarFieldContainer[str]
    tags: _containers.RepeatedScalarFieldContainer[str]
    extra: str
    def __init__(self, url: _Optional[str] = ..., method: _Optional[str] = ..., path: _Optional[str] = ..., content_type: _Optional[str] = ..., status_code: _Optional[int] = ..., source: _Optional[str] = ..., source_url: _Optional[str] = ..., parameters: _Optional[_Iterable[str]] = ..., tags: _Optional[_Iterable[str]] = ..., extra: _Optional[str] = ...) -> None: ...

class Host(_message.Message):
    __slots__ = ()
    HOSTNAME_FIELD_NUMBER: _ClassVar[int]
    LOCAL_IPS_FIELD_NUMBER: _ClassVar[int]
    GATEWAY_IPS_FIELD_NUMBER: _ClassVar[int]
    DNS_SERVERS_FIELD_NUMBER: _ClassVar[int]
    DOMAIN_NAME_FIELD_NUMBER: _ClassVar[int]
    DOMAIN_ROLE_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    hostname: str
    local_ips: _containers.RepeatedScalarFieldContainer[str]
    gateway_ips: _containers.RepeatedScalarFieldContainer[str]
    dns_servers: _containers.RepeatedScalarFieldContainer[str]
    domain_name: str
    domain_role: str
    extra: str
    def __init__(self, hostname: _Optional[str] = ..., local_ips: _Optional[_Iterable[str]] = ..., gateway_ips: _Optional[_Iterable[str]] = ..., dns_servers: _Optional[_Iterable[str]] = ..., domain_name: _Optional[str] = ..., domain_role: _Optional[str] = ..., extra: _Optional[str] = ...) -> None: ...

class Repository(_message.Message):
    __slots__ = ()
    PROVIDER_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    OWNER_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    STARS_FIELD_NUMBER: _ClassVar[int]
    IS_FORK_FIELD_NUMBER: _ClassVar[int]
    MATCHED_DORKS_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    provider: str
    name: str
    url: str
    owner: str
    description: str
    stars: int
    is_fork: bool
    matched_dorks: _containers.RepeatedScalarFieldContainer[str]
    extra: str
    def __init__(self, provider: _Optional[str] = ..., name: _Optional[str] = ..., url: _Optional[str] = ..., owner: _Optional[str] = ..., description: _Optional[str] = ..., stars: _Optional[int] = ..., is_fork: _Optional[bool] = ..., matched_dorks: _Optional[_Iterable[str]] = ..., extra: _Optional[str] = ...) -> None: ...

class Secret(_message.Message):
    __slots__ = ()
    KIND_FIELD_NUMBER: _ClassVar[int]
    DETECTOR_FIELD_NUMBER: _ClassVar[int]
    REDACTED_FIELD_NUMBER: _ClassVar[int]
    FINGERPRINT_FIELD_NUMBER: _ClassVar[int]
    SOURCE_FIELD_NUMBER: _ClassVar[int]
    SOURCE_URL_FIELD_NUMBER: _ClassVar[int]
    FILE_PATH_FIELD_NUMBER: _ClassVar[int]
    LINE_FIELD_NUMBER: _ClassVar[int]
    COMMIT_FIELD_NUMBER: _ClassVar[int]
    VERIFIED_FIELD_NUMBER: _ClassVar[int]
    SEVERITY_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    kind: str
    detector: str
    redacted: str
    fingerprint: str
    source: str
    source_url: str
    file_path: str
    line: int
    commit: str
    verified: bool
    severity: str
    extra: str
    def __init__(self, kind: _Optional[str] = ..., detector: _Optional[str] = ..., redacted: _Optional[str] = ..., fingerprint: _Optional[str] = ..., source: _Optional[str] = ..., source_url: _Optional[str] = ..., file_path: _Optional[str] = ..., line: _Optional[int] = ..., commit: _Optional[str] = ..., verified: _Optional[bool] = ..., severity: _Optional[str] = ..., extra: _Optional[str] = ...) -> None: ...
