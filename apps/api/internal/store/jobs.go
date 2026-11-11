package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Job struct {
	ID          string
	AnalysisID  string
	Step        string
	Attempts    int
	MaxAttempts int
	Payload     []byte
}

func (s *Store) EnqueueJob(ctx context.Context, analysisID, step string, payload any, maxAttempts int) error {
	var raw []byte
	if payload != nil {
		var err error
		if raw, err = json.Marshal(payload); err != nil {
			return err
		}
	}
	_, err := s.Pool.Exec(ctx,
		`insert into jobs(analysis_id, step, payload, max_attempts) values($1,$2,$3,$4)
		 on conflict (analysis_id, step) do nothing`, analysisID, step, raw, maxAttempts)
	return err
}

// ClaimJob atomically leases the next due job (or reclaims an expired lease).
func (s *Store) ClaimJob(ctx context.Context, workerID string) (*Job, error) {
	var j Job
	err := s.Pool.QueryRow(ctx, `
		with candidate as (
		  select id from jobs
		  where (status='pending' and run_at <= now())
		     or (status='running' and locked_at < now() - interval '10 minutes')
		  order by run_at
		  limit 1
		  for update skip locked
		)
		update jobs j
		set status='running', locked_by=$1, locked_at=now(), attempts=j.attempts+1
		from candidate c where j.id = c.id
		returning j.id, j.analysis_id, j.step, j.attempts, j.max_attempts, coalesce(j.payload,'{}'::jsonb)`,
		workerID).Scan(&j.ID, &j.AnalysisID, &j.Step, &j.Attempts, &j.MaxAttempts, &j.Payload)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (s *Store) CompleteJob(ctx context.Context, id string) error {
	_, err := s.Pool.Exec(ctx,
		`update jobs set status='succeeded', locked_by=null, locked_at=null where id=$1`, id)
	return err
}

// MarkJobFailed marks a job permanently failed (permanent errors, no retries).
func (s *Store) MarkJobFailed(ctx context.Context, id, msg string) error {
	_, err := s.Pool.Exec(ctx,
		`update jobs set status='failed', last_error=$2, locked_by=null, locked_at=null where id=$1`, id, msg)
	return err
}

// FailJob records a failure; returns true when the retry budget is exhausted.
func (s *Store) FailJob(ctx context.Context, id, msg string, delay time.Duration) (bool, error) {
	var status string
	err := s.Pool.QueryRow(ctx, `
		update jobs set
		  status = case when attempts >= max_attempts then 'failed' else 'pending' end,
		  run_at = case when attempts >= max_attempts then run_at else now() + $3::interval end,
		  last_error = $2, locked_by = null, locked_at = null
		where id=$1 returning status`,
		id, msg, fmt.Sprintf("%d seconds", int(delay.Seconds()))).Scan(&status)
	return status == "failed", err
}

func (s *Store) RecordEvent(ctx context.Context, analysisID, step, event string, detail any) {
	var raw []byte
	if detail != nil {
		raw, _ = json.Marshal(detail)
	}
	_, _ = s.Pool.Exec(ctx,
		`insert into job_events(analysis_id, step, event, detail) values($1,$2,$3,$4)`,
		analysisID, step, event, raw)
}
