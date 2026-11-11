"""Tree-sitter language registry and per-language extraction queries."""
from __future__ import annotations

import importlib.metadata
from typing import Any, Iterator

from tree_sitter import Language, Node, Parser, Query

import tree_sitter_go as tsgo
import tree_sitter_javascript as tsjs
import tree_sitter_python as tspython
import tree_sitter_typescript as tsts

try:
    LIB_VERSION = importlib.metadata.version("tree-sitter")
except Exception:  # pragma: no cover
    LIB_VERSION = "unknown"

PARSER_VERSION = f"tree-sitter-py-{LIB_VERSION}"
QUERY_VERSION = "queries-1"
SUPPORTED = {"python", "go", "javascript", "typescript"}

QUERIES = {
    "python": """
(function_definition name: (identifier) @name) @symbol.function
(class_definition name: (identifier) @name) @symbol.class
(import_statement) @import
(import_from_statement) @import
""",
    "go": """
(function_declaration name: (identifier) @name) @symbol.function
(method_declaration name: (field_identifier) @name) @symbol.method
(type_declaration (type_spec name: (type_identifier) @name)) @symbol.type
(import_spec path: (interpreted_string_literal) @import.path) @import
""",
    "javascript": """
(function_declaration name: (identifier) @name) @symbol.function
(class_declaration name: (identifier) @name) @symbol.class
(method_definition name: (property_identifier) @name) @symbol.method
(import_statement source: (string) @import.path) @import
""",
    "typescript": """
(function_declaration name: (identifier) @name) @symbol.function
(class_declaration name: (identifier) @name) @symbol.class
(method_definition name: (property_identifier) @name) @symbol.method
(interface_declaration name: (type_identifier) @name) @symbol.interface
(type_alias_declaration name: (type_identifier) @name) @symbol.type
(import_statement source: (string) @import.path) @import
""",
}

_LANG: dict[str, Language] = {}
_PARSER: dict[str, Parser] = {}
_QUERY: dict[str, Query] = {}


def language(name: str) -> Language:
    if name not in _LANG:
        if name == "python":
            _LANG[name] = Language(tspython.language())
        elif name == "go":
            _LANG[name] = Language(tsgo.language())
        elif name == "javascript":
            _LANG[name] = Language(tsjs.language())
        elif name == "typescript":
            _LANG[name] = Language(tsts.language_typescript())
        else:
            raise ValueError(f"unsupported language: {name}")
    return _LANG[name]


def parser(name: str) -> Parser:
    if name not in _PARSER:
        lang = language(name)
        try:
            _PARSER[name] = Parser(lang)
        except TypeError:  # older py-tree-sitter API
            p = Parser()
            p.set_language(lang)
            _PARSER[name] = p
    return _PARSER[name]


def query(name: str) -> Query:
    if name not in _QUERY:
        lang = language(name)
        try:
            _QUERY[name] = Query(lang, QUERIES[name])
        except TypeError:  # older API
            _QUERY[name] = lang.query(QUERIES[name])
    return _QUERY[name]


def iter_matches(name: str, node: Node) -> Iterator[Any]:
    """Yield (pattern_index, captures) pairs across py-tree-sitter API versions.

    Newer py-tree-sitter (0.25+) requires QueryCursor; older versions expose
    Query.matches(). Both return the same (pattern_index, {capture: [nodes]})
    shape, but QueryCursor.captures() returns a dict, which we normalise too.
    """
    q = query(name)
    results = None
    try:  # new API
        from tree_sitter import QueryCursor  # type: ignore

        results = QueryCursor(q).matches(node)
    except Exception:  # old API (or QueryCursor absent)
        results = q.matches(node)

    for item in results:
        if isinstance(item, tuple) and len(item) == 2 and isinstance(item[1], dict):
            yield item
        elif isinstance(item, dict):
            yield (0, item)
        else:  # pragma: no cover - unexpected shape
            continue


def count_errors(node: Node) -> int:
    n = 1 if node.type == "ERROR" else 0
    for child in node.children:
        n += count_errors(child)
    return n


def node_text(node: Node) -> str:
    return node.text.decode("utf-8", "replace") if node.text is not None else ""


def is_exported(node: Node, lang: str, name: str) -> bool:
    if lang == "go":
        return bool(name) and name[0].isupper()
    parent = node.parent
    while parent is not None:
        if parent.type == "export_statement":
            return True
        parent = parent.parent
    return False


def parent_symbol_name(node: Node) -> str | None:
    parent = node.parent
    while parent is not None:
        if parent.type in ("class_definition", "class_declaration", "class_body"):
            for child in parent.children:
                if child.type in ("identifier", "type_identifier", "property_identifier"):
                    return node_text(child)
        parent = parent.parent
    return None
