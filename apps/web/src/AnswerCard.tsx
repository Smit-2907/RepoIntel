import type { AnswerResponse } from "@codebase-intel/core";

export function AnswerCard({ answer }: { answer: AnswerResponse }) {
  return (
    <article className="mt-4 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
      {answer.stub_mode && (
        <div className="mb-3 flex items-center gap-2 rounded-xl bg-amber-50 px-3 py-2 text-xs text-amber-800">
          <span className="h-1.5 w-1.5 rounded-full bg-amber-500" />
          Extractive fallback — no LLM reasoning (provider unavailable or key missing).
        </div>
      )}
      <div className="whitespace-pre-wrap text-sm leading-relaxed text-slate-800">
        {answer.answer_markdown}
      </div>

      {answer.citations.length > 0 && (
        <div className="mt-4 border-t border-slate-100 pt-3">
          <p className="text-xs font-semibold uppercase tracking-wide text-slate-500">Citations</p>
          <div className="mt-2 flex flex-wrap gap-2">
            {answer.citations.map((c, i) => (
              <span
                key={i}
                className="inline-flex items-center gap-1 rounded-lg border border-slate-200 bg-slate-50 px-2 py-1 font-mono text-xs text-slate-700"
              >
                {c.file_path}:{c.start_line}-{c.end_line}
              </span>
            ))}
          </div>
        </div>
      )}

      <Section title="Facts" items={answer.facts} tone="text-emerald-700" />
      <Section title="Hypotheses" items={answer.hypotheses} tone="text-amber-700" />
      <Section title="Unknowns" items={answer.unknowns} tone="text-slate-500" />
    </article>
  );
}

function Section({ title, items, tone }: { title: string; items: string[] | null; tone: string }) {
  if (!items || items.length === 0) return null;
  return (
    <div className="mt-3 border-t border-slate-100 pt-3">
      <p className={`text-xs font-semibold uppercase tracking-wide ${tone}`}>{title}</p>
      <ul className="mt-1 list-disc space-y-1 pl-5 text-xs text-slate-600">
        {items.map((item, i) => (
          <li key={i}>{item}</li>
        ))}
      </ul>
    </div>
  );
}
