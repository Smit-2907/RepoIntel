import type { AnalysisResponse } from "@codebase-intel/core";

const STEPS = [
  "queued",
  "fetching",
  "downloading",
  "extracting",
  "parsing",
  "embedding",
  "summarizing",
  "ready",
];

const LABEL: Record<string, string> = {
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

function statusBadge(status: string): string {
  const base = "rounded-full px-3 py-1 text-xs font-semibold";
  switch (status) {
    case "ready":
      return `${base} bg-emerald-50 text-emerald-700`;
    case "partially_ready":
      return `${base} bg-amber-50 text-amber-700`;
    case "failed":
      return `${base} bg-red-50 text-red-700`;
    default:
      return `${base} bg-blue-50 text-blue-700`;
  }
}

function progress(status: string): { pct: number; bar: string } {
  if (status === "failed") return { pct: 100, bar: "bg-red-500" };
  const step = status === "partially_ready" ? "ready" : status;
  const idx = STEPS.indexOf(step);
  const pct = idx < 0 ? 8 : Math.round(((idx + 1) / STEPS.length) * 100);
  const bar =
    status === "ready" ? "bg-emerald-500" : status === "partially_ready" ? "bg-amber-500" : "bg-blue-500";
  return { pct, bar };
}

export function StatusCard({ analysis }: { analysis: AnalysisResponse }) {
  const { pct, bar } = progress(analysis.status);
  return (
    <section className="mt-6 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <p className="text-sm font-semibold text-slate-900">{analysis.repository}</p>
          <p className="mt-0.5 text-xs text-slate-500">
            {analysis.default_branch} @ {analysis.commit_sha.slice(0, 12)} ·{" "}
            {analysis.size_kb.toLocaleString()} KB
          </p>
        </div>
        <span className={statusBadge(analysis.status)}>{LABEL[analysis.status] ?? analysis.status}</span>
      </div>

      <div className="mt-4">
        <div className="h-2 w-full overflow-hidden rounded-full bg-slate-100">
          <div className={`h-full rounded-full transition-all ${bar}`} style={{ width: `${pct}%` }} />
        </div>
        <p className="mt-2 text-xs text-slate-600">{analysis.status_detail || LABEL[analysis.status]}</p>
      </div>

      <div className="mt-4 grid grid-cols-2 gap-2 sm:grid-cols-4">
        <Stat label="Files" value={analysis.counts.files} />
        <Stat label="Indexed" value={analysis.counts.indexed_files} />
        <Stat label="Symbols" value={analysis.counts.symbols} />
        <Stat label="Chunks" value={analysis.counts.chunks} />
      </div>

      {analysis.status === "partially_ready" && (
        <p className="mt-4 rounded-xl border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800">
          Some files could not be parsed; answers may be incomplete. Parse-error rate:{" "}
          {(analysis.parse_error_rate * 100).toFixed(0)}%.
        </p>
      )}
    </section>
  );
}

function Stat({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-xl border border-slate-100 bg-slate-50 px-3 py-2">
      <p className="text-lg font-semibold text-slate-900">{value.toLocaleString()}</p>
      <p className="text-xs text-slate-500">{label}</p>
    </div>
  );
}
