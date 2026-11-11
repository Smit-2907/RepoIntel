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
    onSuccess: (d) => {
      setAnswer(d);
      setInvestigation(null);
      setError(null);
    },
    onError: (e: Error) => setError(e.message),
  });
  const inv = useMutation({
    mutationFn: (t: string) => investigate(analysisId, t),
    onSuccess: (d) => {
      setInvestigation(d);
      setAnswer(null);
      setError(null);
    },
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
    if (!parsed.success) {
      setError(parsed.error.issues[0]?.message ?? "Invalid question");
      return;
    }
    ask.mutate(parsed.data.question);
  }

  function submitTrace(e: React.FormEvent) {
    e.preventDefault();
    const parsed = InvestigationRequest.safeParse({ stack_trace: trace });
    if (!parsed.success) {
      setError(parsed.error.issues[0]?.message ?? "Paste a stack trace");
      return;
    }
    inv.mutate(parsed.data.stack_trace);
  }

  const seg = (active: boolean) =>
    `rounded-lg px-3 py-1.5 text-xs font-semibold transition ${
      active ? "bg-white text-slate-900 shadow-sm" : "text-slate-600 hover:text-slate-900"
    }`;

  return (
    <section className="mt-6 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
      <div className="flex flex-wrap items-center gap-2 border-b border-slate-100 pb-3">
        <div className="flex rounded-xl bg-slate-100 p-1">
          <button type="button" onClick={() => setMode("ask")} className={seg(mode === "ask")}>
            Ask a question
          </button>
          <button type="button" onClick={() => setMode("investigate")} className={seg(mode === "investigate")}>
            Investigate a stack trace
          </button>
        </div>
        <button
          type="button"
          onClick={onExport}
          disabled={exporting}
          className="ml-auto rounded-xl border border-slate-300 bg-white px-3 py-2 text-xs font-semibold text-slate-700 transition hover:bg-slate-50 disabled:opacity-50"
        >
          {exporting ? "Building…" : "Download Briefing Pack"}
        </button>
      </div>

      <div className="mt-4">
        {mode === "ask" ? (
          <form onSubmit={submitQuestion} className="flex flex-col gap-2 sm:flex-row">
            <input
              value={question}
              onChange={(e) => setQuestion(e.target.value)}
              placeholder="Where is routing handled?"
              className="w-full rounded-xl border border-slate-300 bg-white px-4 py-3 text-sm text-slate-900 placeholder:text-slate-400 focus:border-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-900/10"
            />
            <button
              type="submit"
              disabled={ask.isPending}
              className="rounded-xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white transition hover:bg-slate-800 disabled:opacity-50"
            >
              {ask.isPending ? "Thinking…" : "Ask"}
            </button>
          </form>
        ) : (
          <form onSubmit={submitTrace}>
            <textarea
              value={trace}
              onChange={(e) => setTrace(e.target.value)}
              rows={8}
              placeholder={'Traceback (most recent call last):\n  File "src/flask/app.py", line 123, in dispatch_request\nKeyError: None'}
              className="w-full rounded-xl border border-slate-300 bg-white px-4 py-3 font-mono text-xs text-slate-900 placeholder:text-slate-400 focus:border-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-900/10"
            />
            <div className="mt-2 flex items-center gap-3">
              <button
                type="submit"
                disabled={inv.isPending}
                className="rounded-xl bg-slate-900 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-slate-800 disabled:opacity-50"
              >
                {inv.isPending ? "Investigating…" : "Investigate"}
              </button>
              <p className="text-xs text-slate-500">Python, JS, or Go traces.</p>
            </div>
          </form>
        )}
      </div>

      {error && (
        <div className="mt-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>
      )}

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

