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
    <div className="mx-auto max-w-4xl px-4 py-10 text-slate-800">
      <header className="mb-8">
        <h1 className="text-2xl font-semibold">AI Codebase Intelligence</h1>
        <p className="mt-1 text-sm text-slate-500">
          Paste a public GitHub URL → deterministic Tree-sitter index → evidence-grounded answers with
          file:line citations.
        </p>
      </header>

      <form onSubmit={submitUrl} className="flex gap-2">
        <input
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://github.com/pallets/flask"
          className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm"
        />
        <button
          type="submit"
          disabled={create.isPending}
          className="rounded-lg bg-slate-900 px-4 py-2 text-sm font-medium text-white disabled:opacity-50"
        >
          {create.isPending ? "Starting…" : "Analyze"}
        </button>
      </form>

      {error && <p className="mt-3 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>}

      {analysis.data && <StatusCard analysis={analysis.data} />}

      {analysisId && isReady && <Workspace analysisId={analysisId} />}
    </div>
  );
}
