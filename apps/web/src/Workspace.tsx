import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { InvestigationRequest, QuestionRequest, type AnswerResponse } from "@codebase-intel/core";
import { askQuestion, downloadBriefing, investigate, type InvestigationResponse } from "./api";
import { AnswerCard } from "./components";

export function Workspace({ analysisId }: { analysisId: string }) {
  const [mode, setMode] = useState<"ask" | "investigate">("ask");
  const [question, setQuestion] = useState("");
  const [trace, setTrace] = useState("");
  const [answer, setAnswer] = useState<AnswerResponse | null>(null);
  const [investigation, setInvestigation] = useState<InvestigationResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [exporting, setExporting] = useState(false);

  const ask = useMutation({
    mutationFn: (q: string) => askQuestion(analysisId, q),
    onSuccess: (d) => { setAnswer(d); setInvestigation(null); setError(null); },
    onError: (e: Error) => setError(e.message),
  });
  const inv = useMutation({
    mutationFn: (t: string) => investigate(analysisId, t),
    onSuccess: (d) => { setInvestigation(d); setAnswer(null); setError(null); },
    onError: (e: Error) => setError(e.message),
  });

  async function onExport() {
    try {
      setExporting(true);
      const blob = await downloadBriefing(analysisId);
      const a = document.createElement("a");
      a.href = URL.createObjectURL(blob);
      a.download = "briefing.zip";
      a.click();
      URL.revokeObjectURL(a.href);
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setExporting(false);
    }
  }

  function submitQuestion(e: React.FormEvent) {
    e.preventDefault();
    const parsed = QuestionRequest.safeParse({ question });
    if (!parsed.success) { setError(parsed.error.issues[0]?.message ?? "Invalid question"); return; }
    ask.mutate(parsed.data.question);
  }

  function submitTrace(e: React.FormEvent) {
    e.preventDefault();
    const parsed = InvestigationRequest.safeParse({ stack_trace: trace });
    if (!parsed.success) { setError(parsed.error.issues[0]?.message ?? "Paste a stack trace"); return; }
    inv.mutate(parsed.data.stack_trace);
  }

  const tab = (active: boolean) =>
    `rounded-lg px-3 py-1.5 text-sm font-medium ${active ? "bg-slate-900 text-white" : "bg-slate-100 text-slate-700"}`;

  return (
    <section className="mt-6">
      <div className="mb-3 flex flex-wrap items-center gap-2">
        <button type="button" onClick={() => setMode("ask")} className={tab(mode === "ask")}>Ask a question</button>
        <button type="button" onClick={() => setMode("investigate")} className={tab(mode === "investigate")}>
          Investigate a stack trace
        </button>
        <button
          type="button"
          onClick={onExport}
          disabled={exporting}
          className="ml-auto rounded-lg border border-slate-300 px-3 py-1.5 text-sm font-medium disabled:opacity-50"
        >
          {exporting ? "Building…" : "Download Briefing Pack"}
        </button>
      </div>

      {mode === "ask" ? (
        <form onSubmit={submitQuestion} className="flex gap-2">
          <input
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
            placeholder="Where is routing handled?"
            className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm"
          />
          <button type="submit" disabled={ask.isPending}
            className="rounded-lg bg-slate-900 px-4 py-2 text-sm font-medium text-white disabled:opacity-50">
            {ask.isPending ? "Thinking…" : "Ask"}
          </button>
        </form>
      ) : (
        <form onSubmit={submitTrace}>
          <textarea
            value={trace}
            onChange={(e) => setTrace(e.target.value)}
            rows={8}
            placeholder={'Traceback (most recent call last):\n  File "src/flask/app.py", line 123, in ...'}
            className="w-full rounded-lg border border-slate-300 px-3 py-2 font-mono text-xs"
          />
          <button type="submit" disabled={inv.isPending}
            className="mt-2 rounded-lg bg-slate-900 px-4 py-2 text-sm font-medium text-white disabled:opacity-50">
            {inv.isPending ? "Investigating…" : "Investigate"}
          </button>
        </form>
      )}

      {error && <p className="mt-3 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>}
      {answer && <AnswerCard answer={answer} />}
      {investigation && (
        <div className="mt-4">
          <p className="text-xs text-slate-500">
            Parsed {investigation.parsed_frames} frames · matched {investigation.matched_frames} to indexed files
          </p>
          <AnswerCard answer={investigation} />
        </div>
      )}
    </section>
  );
}
