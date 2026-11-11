"""HTTP surface for the parser service (internal, not public).

Contract: POST /parse  (design v1.1 §3 step 5)
"""
from __future__ import annotations

from fastapi import FastAPI
from pydantic import BaseModel, Field

from .parse_service import PARSER_VERSION, QUERY_VERSION, parse_files

app = FastAPI(title="codebase-intel-parser", version="0.1.0")


class FileIn(BaseModel):
    path: str = Field(min_length=1, max_length=500)
    language: str = Field(min_length=1, max_length=20)
    content: str


class ParseRequest(BaseModel):
    analysis_id: str
    files: list[FileIn]
    max_files: int = 3000
    max_chunks: int = 15000
    max_chunk_lines: int = 120


@app.get("/health")
def health() -> dict:
    return {
        "status": "ok",
        "parser_version": PARSER_VERSION,
        "query_version": QUERY_VERSION,
    }


@app.post("/parse")
def parse(req: ParseRequest) -> dict:
    return parse_files(
        files=[f.model_dump() for f in req.files],
        max_files=req.max_files,
        max_chunks=req.max_chunks,
        max_chunk_lines=req.max_chunk_lines,
    )
