from app.parse_service import parse_file

PY_SOURCE = (
    "import os\n"
    "from helpers import util\n\n"
    "def add(a, b):\n"
    "    return a + b\n\n"
    "class Calc:\n"
    "    def mul(self, a, b):\n"
    "        return a * b\n"
)


def test_python_symbols_and_chunks():
    result = parse_file("calc.py", "python", PY_SOURCE, 120)
    names = {s["name"] for s in result["symbols"]}
    assert "add" in names
    assert "Calc" in names
    assert result["file"]["is_indexed"] is True
    assert result["file"]["parse_error_count"] == 0
    assert len(result["chunks"]) >= 2
    assert all(c["start_line"] >= 1 for c in result["chunks"])
    assert all(c["end_line"] >= c["start_line"] for c in result["chunks"])


def test_python_imports_extracted():
    result = parse_file("m.py", "python", PY_SOURCE, 120)
    assert len(result["imports"]) >= 1


def test_broken_file_is_still_indexed():
    result = parse_file("bad.py", "python", "def broken(:\n    pass\n", 120)
    assert result["file"]["is_indexed"] is True
    assert result["file"]["parse_error_count"] >= 1


def test_unsupported_language_not_indexed():
    result = parse_file("README.md", "other", "# hello\n", 120)
    assert result["file"]["is_indexed"] is False
    assert result["symbols"] == []
    assert result["chunks"] == []
