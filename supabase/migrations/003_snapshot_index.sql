-- 003_snapshot_index.sql — allow multiple analyses of the same snapshot.
-- The partial unique index blocked finalize for repeated analyses of the same
-- repo+commit (common during demos/testing). Idempotent reuse will be handled
-- at the API level later; the DB must not reject a valid analysis.

drop index if exists analyses_snapshot_ready_uniq;
create index if not exists analyses_snapshot_idx on analyses(repository_id, commit_sha);
