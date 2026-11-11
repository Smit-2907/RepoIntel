#!/usr/bin/env bash
# Smoke test for the vertical slice: create an analysis, poll to completion, ask one question.
set -euo pipefail

API="${API:-http://localhost:8081}"
REPO="${1:-https://github.com/pallets/flask}"

json_get() { python3 -c "import sys,json;d=json.load(sys.stdin);print(d$1)"; }

echo "→ health"
curl -fsS "$API/health"; echo

echo "→ creating analysis for $REPO"
RESP=$(curl -fsS -X POST "$API/api/analyses" -H 'Content-Type: application/json' \
  -d "{\"repo_url\":\"$REPO\"}")
ID=$(printf '%s' "$RESP" | json_get "['analysis_id']")
echo "analysis_id=$ID"

STATUS=""
for i in $(seq 1 300); do
  BODY=$(curl -fsS "$API/api/analyses/$ID")
  STATUS=$(printf '%s' "$BODY" | json_get "['status']")
  DETAIL=$(printf '%s' "$BODY" | json_get ".get('status_detail','')")
  echo "[$i] $STATUS — $DETAIL"
  case "$STATUS" in ready|partially_ready|failed) break;; esac
  sleep 2
done

if [ "$STATUS" = "failed" ]; then
  echo "✗ analysis failed"
  printf '%s\n' "$BODY"
  exit 1
fi

echo "→ asking a question"
curl -fsS -X POST "$API/api/analyses/$ID/questions" -H 'Content-Type: application/json' \
  -d '{"question":"Where is routing handled?"}' | python3 -m json.tool
echo
echo "✓ smoke test finished (analysis_id=$ID)"
