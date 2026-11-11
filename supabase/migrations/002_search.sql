-- 002_search.sql — single auditable hybrid retrieval function
-- Design: docs/Architecture_and_Data_Design_v1.1.md §8.3
-- Ablations: vector-only = p_full_text_weight 0; hybrid = 1,1; hybrid+graph = post-processing.

create or replace function hybrid_search(
  p_analysis_id uuid,
  p_query_text text,
  p_query_embedding vector(1536),
  p_match_count int default 12,
  p_full_text_weight float default 1.0,
  p_semantic_weight float default 1.0,
  p_rrf_k int default 50
) returns table (
  chunk_id uuid,
  file_path text,
  symbol_name text,
  start_line int,
  end_line int,
  content text,
  ft_rank int,
  vec_rank int,
  rrf_score float
) language sql stable as $$
with ft as (
  select c.id,
         row_number() over (
           order by ts_rank_cd(c.tsv, websearch_to_tsquery('simple', p_query_text)) desc
         ) as rank
  from code_chunks c
  where c.analysis_id = p_analysis_id
    and c.tsv @@ websearch_to_tsquery('simple', p_query_text)
  limit p_match_count * 2
),
vec as (
  select c.id,
         row_number() over (order by c.embedding <=> p_query_embedding) as rank
  from code_chunks c
  where c.analysis_id = p_analysis_id
    and c.embedding is not null
  order by c.embedding <=> p_query_embedding
  limit p_match_count * 2
)
select
  c.id as chunk_id,
  f.path as file_path,
  s.name as symbol_name,
  c.start_line,
  c.end_line,
  c.content,
  coalesce(ft.rank, 0)::int as ft_rank,
  coalesce(vec.rank, 0)::int as vec_rank,
  (
    coalesce(p_full_text_weight / (p_rrf_k + ft.rank), 0) +
    coalesce(p_semantic_weight  / (p_rrf_k + vec.rank), 0)
  )::float as rrf_score
from code_chunks c
join analysis_files f on f.id = c.file_id
left join code_symbols s on s.id = c.symbol_id
left join ft on ft.id = c.id
left join vec on vec.id = c.id
where c.analysis_id = p_analysis_id
  and (ft.id is not null or vec.id is not null)
order by rrf_score desc
limit p_match_count;
$$;
