"""Node flags, read from the schemas the runtime has loaded.

A flag is a bit on `Node.flags_mask`, and which bits exist is declared by
extensions in their schema documents — `easm` declares `honeypot`, `noise` and
the rest, because they are its judgements about an asset, not the graph
store's. The runtime holds the mechanism (a 64-bit mask) and never the
vocabulary.

This class used to be eight constants compiled into the native binding, which
is why no extension but the built-in one could ever have a flag. The names and
values are unchanged; only the place they are declared moved.

Names resolve in either spelling, so `NodeFlags.HONEYPOT` and
`NodeFlags.bit("honeypot")` are the same question.
"""

from __future__ import annotations

from typing import Dict, Iterator, Tuple

from cstxpy.schema import registry

__all__ = ("NodeFlags",)

#: Bits 56-63 belong to the runtime; extensions declare 0-55.
EXTENSION_FLAG_BITS = 56


def _declared() -> Dict[str, Tuple[int, bool]]:
    """`name -> (bit, default_exclude)` across every loaded extension.

    Read through the registry on each call rather than cached: an extension
    registered at runtime gets its flags answered on the same terms as a
    bundled one, which is the whole point of declaring them.
    """
    flags: Dict[str, Tuple[int, bool]] = {}
    claimed: Dict[int, str] = {}
    for name in registry.extensions():
        schema = registry.extension(name)
        if schema is None:
            continue
        for flag, declaration in schema.flags.items():
            # First claimant of a bit keeps it: a bit is what a stored mask
            # means, so a second claim would make one stored value ambiguous.
            # The core refuses the same way at registration.
            if claimed.setdefault(declaration.bit, flag) != flag:
                continue
            flags[flag] = (declaration.bit, declaration.default_exclude)
    return flags


class _NodeFlagsMeta(type):
    """Resolves `NodeFlags.HONEYPOT` against what extensions declared."""

    def __getattr__(cls, name: str) -> int:
        if name.startswith("_"):
            raise AttributeError(name)
        bit = cls.bit(name)
        if bit is None:
            raise AttributeError(
                f"no loaded extension declares a node flag named {name.lower()!r}; "
                f"declared: {', '.join(sorted(_declared())) or 'none'}"
            )
        return 1 << bit

    def __dir__(cls) -> list:
        return [*type.__dir__(cls), *(name.upper() for name in _declared())]


class NodeFlags(metaclass=_NodeFlagsMeta):
    """Node flag bits, as declared by the loaded extensions.

    A namespace, not a wrapper: graph APIs take ordinary ``int`` masks. Names
    resolve in either spelling — ``NodeFlags.HONEYPOT`` and
    ``NodeFlags.bit("honeypot")`` ask the same question.
    """

    NONE = 0

    @staticmethod
    def bit(name: str) -> "int | None":
        """The bit one declared flag occupies, or ``None`` if undeclared."""
        declaration = _declared().get(str(name).lower())
        return None if declaration is None else declaration[0]

    @staticmethod
    def mask(name: str) -> int:
        """The single-bit mask for one declared flag name; 0 if undeclared."""
        bit = NodeFlags.bit(name)
        return 0 if bit is None else 1 << bit

    @staticmethod
    def all_mask() -> int:
        """Every bit any loaded extension declared."""
        return sum(1 << bit for bit, _ in _declared().values())

    @staticmethod
    def default_exclude_mask() -> int:
        """The bits extensions advise hiding from an ordinary view.

        Advice, not enforcement — nothing applies it on its own. A caller asks
        for it and passes it back as a filter, which is where the policy
        belongs: the same flag is noise on an inventory page and the whole
        point on a threat page.
        """
        return sum(1 << bit for bit, exclude in _declared().values() if exclude)

    @staticmethod
    def names() -> Iterator[Tuple[str, int]]:
        """Declared ``(name, bit)`` pairs, lowest bit first."""
        declared = _declared()
        for name in sorted(declared, key=lambda key: declared[key][0]):
            yield name, declared[name][0]
