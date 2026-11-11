import type { AnalysisResponse, AnswerResponse } from "@codebase-intel/core";

const STEP_LABEL: Record<string, string> = {
  queued: "Queued",
  fetching: "Fetching metadata",
  downloading: "Downloading archive",
  extracting: "Extracting files",
  parsing: "Parsing with Tree-sitter",
  embedding: "Embedding chunks",
  summarizing: "Building graph + summary",
  ready: "Ready",
  partially_ready: "Partially ready",
  failed: "Failed",
};

export function StatusCard({ analysis }: { analysis: AnalysisResponse }) {
  return (
    <section className="mt-6 rounded-xl border border-slate-200 p-4">
      <div className="flex items-center justify-between">
        <div>
          <p className="font-medium">{analysis.repository}</p>
          <p className="text-xs text-slate-500">
            {analysis.default_branch} @ {analysis.commit_sha.slice(0, 12)} · {analysis.size_kb} KB
          </p>
        </div>
        <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium">
          {STEP_LABEL[analysis.status] ?? analysis.status}
        </span>
      </div>
      <p className="mt-2 text-sm text-slate-600">{analysis.status_detail}</p>
      <div className="mt-3 grid grid-cols-4 gap-2 text-center text-xs">
        <Stat label="Files" value={analysis.counts.files} />
        <Stat label="Indexed" value={analysis.counts.indexed_files} />
        <Stat label="Symbols" value={analysis.counts.symbols} />
        <Stat label="Chunks" value={analysis.counts.chunks} />
      </div>
      {analysis.status === "partially_ready" && (
        <p className="mt-3 rounded-lg bg-amber-50 px-3 py-2 text-xs text-amber-800">
          Some files could not be parsed; answers may be incomplete. Parse-error rate:{" "}
          {(analysis.parse_error_rate * 100).toFixed(0)}%.
        </p>
      )}
    </section>
  );
}

export function AnswerCard({ answer }: { answer: AnswerResponse }) {
  return (
    <article className="mt-4 rounded-xl border border-slate-200 p-4">
      {answer.stub_mode && (
        <p className="mb-2 rounded-lg bg-amber-50 px-3 py-1 text-xs text-amber-800">
          Stub mode (no GEMINI_API_KEY): extractive answer only.
        </p>
      )}
      <pre className="whitespace-pre-wrap text-sm">{answer.answer_markdown}</pre>
      {answer.citations.length > 0 && (
        <div className="mt-4">
          <p className="text-xs font-semibold uppercase tracking-wide text-slate-500">Citations</p>
          <ul className="mt-1 space-y-1 text-sm">
            {answer.citations.map((c, i) => (
              <li key={i} className="rounded bg-slate-50 px-2 py-1">
                <code>{c.file_path}</code>:{c.start_line}-{c.end_line}
                {c.symbol ? ` · ${c.symbol}` : ""}{" "}
                <span className="text-xs text-slate-500">({c.kind})</span>
              </li>
            ))}
          </ul>
        </div>
      )}
      {answer.unknowns && answer.unknowns.length > 0 && (
        <div className="mt-3 text-xs text-slate-500">
          <p className="font-semibold uppercase tracking-wide">Unknowns</p>
          <ul className="list-disc pl-4">
            {answer.unknowns.map((u, i) => (
              <li key={i}>{u}</li>
            ))}
          </ul>
        </div>
      )}
    </article>
  );
}

function Stat({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-lg bg-slate-50 px-2 py-2">
      <p className="text-lg font-semibold">{value}</p>
      <p className="text-slate-500">{label}</p>
    </div>
  );
}
