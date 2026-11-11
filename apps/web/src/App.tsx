import { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { CreateAnalysisRequest } from "@codebase-intel/core";
import { createAnalysis, getAnalysis } from "./api";
import { StatusCard } from "./components";
import { Workspace } from "./Workspace";

const TERMINAL = new Set(["ready", "partially_ready", "failed"]);

export default function App() {
  const [url, setUrl] = useState("");
  const [analysisId, setAnalysisId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const create = useMutation({
    mutationFn: createAnalysis,
    onSuccess: (data) => {
      setAnalysisId(data.analysis_id);
      setError(null);
    },
    onError: (e: Error) => setError(e.message),
  });

  const analysis = useQuery({
    queryKey: ["analysis", analysisId],
    queryFn: () => getAnalysis(analysisId!),
    enabled: !!analysisId,
    refetchInterval: (q) => (q.state.data && TERMINAL.has(q.state.data.status) ? false : 2000),
  });

  function submitUrl(e: React.FormEvent) {
    e.preventDefault();
    const parsed = CreateAnalysisRequest.safeParse({ repo_url: url.trim() });
    if (!parsed.success) {
      setError(parsed.error.issues[0]?.message ?? "Invalid URL");
      return;
    }
    create.mutate(parsed.data.repo_url);
  }

  const status = analysis.data?.status;
  const isReady = status === "ready" || status === "partially_ready";

  return (
    <div className="min-h-screen bg-[#f6f7fb]">
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex max-w-5xl items-center gap-3 px-6 py-4">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-slate-900 text-xs font-bold text-white">
            CI
          </div>
          <div className="min-w-0">
            <p className="text-sm font-semibold text-slate-900">Codebase Intelligence</p>
            <p className="truncate text-xs text-slate-500">Evidence-grounded repository analysis</p>
          </div>
          <span className="ml-auto hidden rounded-full bg-emerald-50 px-3 py-1 text-xs font-medium text-emerald-700 sm:block">
            MVP · v1.6
          </span>
        </div>
      </header>

      <main className="mx-auto max-w-5xl px-6 py-10">
        <section className="mb-8">
          <h1 className="text-3xl font-semibold tracking-tight text-slate-900">
            Understand any public repository
          </h1>
          <p className="mt-2 max-w-2xl text-sm leading-relaxed text-slate-600">
            Paste a GitHub URL. We index the code with Tree-sitter, build a dependency graph, and answer
            questions with exact file and line citations — no guessing.
          </p>
        </section>

        <form onSubmit={submitUrl} className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <label
            htmlFor="repo-url"
            className="mb-2 block text-xs font-semibold uppercase tracking-wide text-slate-500"
          >
            Repository URL
          </label>
          <div className="flex flex-col gap-2 sm:flex-row">
            <input
              id="repo-url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="https://github.com/pallets/flask"
              className="w-full rounded-xl border border-slate-300 bg-white px-4 py-3 text-sm text-slate-900 placeholder:text-slate-400 focus:border-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-900/10"
            />
            <button
              type="submit"
              disabled={create.isPending}
              className="rounded-xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white transition hover:bg-slate-800 disabled:opacity-50"
            >
              {create.isPending ? "Starting…" : "Analyze repository"}
            </button>
          </div>
          <p className="mt-2 text-xs text-slate-500">
            Public repositories only · pinned to the default-branch commit
          </p>
        </form>

        {error && (
          <div className="mt-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        )}

        {analysis.data && <StatusCard analysis={analysis.data} />}
        {analysisId && isReady && <Workspace analysisId={analysisId} />}
      </main>
    </div>
  );
}

