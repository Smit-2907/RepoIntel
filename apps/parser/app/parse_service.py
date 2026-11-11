"""Deterministic extraction: files, symbols, imports, symbol-bounded chunks.

Never executes repository code (design v1.1 §3 step 5, §12.1).
"""
from __future__ import annotations

import hashlib
import time
from typing import Any

from .languages import (
    PARSER_VERSION,
    QUERY_VERSION,
    SUPPORTED,
    count_errors,
    is_exported,
    iter_matches,
    node_text,
    parent_symbol_name,
    parser,
)


def _sha(text: str) -> str:
    return hashlib.sha256(text.encode("utf-8", "replace")).hexdigest()


def _line_count(text: str) -> int:
    return text.count("\n") + 1 if text else 0


def _chunk_for(path: str, name: str | None, content: str, start: int, end: int,
               max_lines: int) -> list[dict[str, Any]]:
    lines = content.split("\n")
    header = f"file: {path}\nsymbol: {name or '(file)'}\nlines: {start}-{end}\n---\n"
    chunks: list[dict[str, Any]] = []
    cursor = start
    while cursor <= end:
        stop = min(end, cursor + max_lines - 1)
        text = header + "\n".join(lines[cursor - 1:stop])
        chunks.append({
            "path": path,
            "symbol_name": name,
            "start_line": cursor,
            "end_line": stop,
            "content": text,
            "content_hash": _sha(text),
            "token_estimate": max(1, len(text) // 4),
        })
        cursor = stop + 1
    return chunks


def parse_file(path: str, language: str, content: str, max_chunk_lines: int) -> dict[str, Any]:
    entry = {
        "path": path,
        "language": language,
        "size_bytes": len(content.encode("utf-8", "replace")),
        "line_count": _line_count(content),
        "parse_error_count": 0,
        "is_indexed": language in SUPPORTED,
    }
    symbols: list[dict[str, Any]] = []
    imports: list[dict[str, Any]] = []
    chunks: list[dict[str, Any]] = []
    if language not in SUPPORTED:
        return {"file": entry, "symbols": symbols, "imports": imports, "chunks": chunks}

    tree = parser(language).parse(content.encode("utf-8", "replace"))
    entry["parse_error_count"] = count_errors(tree.root_node)

    found_symbol = False
    for _, caps in iter_matches(language, tree.root_node):
        sym_kind = None
        sym_node = None
        for key, nodes in caps.items():
            if key.startswith("symbol.") and nodes:
                sym_kind = key.split(".", 1)[1]
                sym_node = nodes[0]
                break
        if sym_node is not None and sym_kind:
            name_nodes = caps.get("name") or []
            name = node_text(name_nodes[0]) if name_nodes else ""
            start = sym_node.start_point[0] + 1
            end = sym_node.end_point[0] + 1
            parent = parent_symbol_name(sym_node)
            qualified = f"{parent}.{name}" if parent and name else name
            symbols.append({
                "path": path,
                "name": name,
                "qualified_name": qualified,
                "kind": sym_kind,
                "start_line": start,
                "end_line": end,
                "parent_name": parent,
                "is_exported": is_exported(sym_node, language, name),
            })
            chunks.extend(_chunk_for(path, qualified or name, content, start, end, max_chunk_lines))
            found_symbol = True
            continue
        imp_nodes = caps.get("import") or []
        if imp_nodes:
            path_nodes = caps.get("import.path") or []
            raw = node_text(path_nodes[0]) if path_nodes else node_text(imp_nodes[0])
            raw = raw.strip().strip('"').strip("'").strip()
            if raw:
                imports.append({"from_path": path, "raw_specifier": raw[:500], "kind": "unknown"})

    if not found_symbol:
        chunks.extend(_chunk_for(path, None, content, 1, max(1, _line_count(content)), max_chunk_lines))

    return {"file": entry, "symbols": symbols, "imports": imports, "chunks": chunks}


def parse_files(files: list[dict[str, Any]], max_files: int, max_chunks: int,
                max_chunk_lines: int) -> dict[str, Any]:
    started = time.time()
    result: dict[str, Any] = {"files": [], "symbols": [], "imports": [], "chunks": []}
    for f in files[:max_files]:
        parsed = parse_file(f["path"], f["language"], f["content"], max_chunk_lines)
        result["files"].append(parsed["file"])
        result["symbols"].extend(parsed["symbols"])
        result["imports"].extend(parsed["imports"])
        result["chunks"].extend(parsed["chunks"])
    result["chunks"] = result["chunks"][:max_chunks]
    result["meta"] = {
        "parser_version": PARSER_VERSION,
        "query_version": QUERY_VERSION,
        "duration_ms": int((time.time() - started) * 1000),
    }
    return result
