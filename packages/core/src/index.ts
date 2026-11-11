import { z } from "zod";

// Client-side validation (UX only); the Go API is authoritative.

export const RepoUrl = z
  .string()
  .max(300)
  .regex(
    /^https:\/\/github\.com\/[A-Za-z0-9-]{1,39}\/[A-Za-z0-9._-]{1,100}(?:\.git)?\/?$/,
    "Must be a public GitHub URL: https://github.com/owner/repo",
  );

export const CreateAnalysisRequest = z.object({ repo_url: RepoUrl });
export const QuestionRequest = z.object({
  question: z.string().trim().min(1, "Ask a question").max(2000, "Max 2,000 characters"),
});

export const AnalysisStatus = z.enum([
  "queued",
  "fetching",
  "downloading",
  "extracting",
  "parsing",
  "embedding",
  "summarizing",
  "ready",
  "partially_ready",
  "failed",
]);

export const Citation = z
  .object({
    file_path: z.string(),
    start_line: z.number().int().min(1),
    end_line: z.number().int().min(1),
    symbol: z.string().optional(),
    kind: z.enum(["fact", "hypothesis"]),
  })
  .refine((c) => c.end_line >= c.start_line, { message: "end_line must be >= start_line" });

export const InvestigationRequest = z.object({
  stack_trace: z.string().trim().min(10, "Paste a stack trace").max(20000, "Max 20,000 characters"),
});

export const AnalysisResponse = z.object({
  id: z.string(),
  repository: z.string(),
  default_branch: z.string(),
  commit_sha: z.string(),
  status: AnalysisStatus,
  status_detail: z.string(),
  size_kb: z.number(),
  counts: z.object({
    files: z.number(),
    indexed_files: z.number(),
    symbols: z.number(),
    chunks: z.number(),
  }),
  parse_error_rate: z.number(),
});

export const AnswerResponse = z.object({
  message_id: z.string(),
  answer_markdown: z.string(),
  citations: z.array(Citation),
  facts: z.array(z.string()).nullable(),
  hypotheses: z.array(z.string()).nullable(),
  unknowns: z.array(z.string()).nullable(),
  usage: z.object({
    input_tokens: z.number(),
    output_tokens: z.number(),
    latency_ms: z.number(),
  }),
  stub_mode: z.boolean(),
});

export type CreateAnalysisRequest = z.infer<typeof CreateAnalysisRequest>;
export type InvestigationRequest = z.infer<typeof InvestigationRequest>;
export type QuestionRequest = z.infer<typeof QuestionRequest>;
export type AnalysisStatus = z.infer<typeof AnalysisStatus>;
export type Citation = z.infer<typeof Citation>;
export type AnalysisResponse = z.infer<typeof AnalysisResponse>;
export type AnswerResponse = z.infer<typeof AnswerResponse>;
