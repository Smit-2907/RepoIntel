import type { AnalysisResponse, AnswerResponse } from "@codebase-intel/core";

const API = (import.meta.env.VITE_API_URL as string | undefined) ?? "http://localhost:8081";

async function errorMessage(res: Response): Promise<string> {
  try {
    const body = await res.json();
    return body?.error?.message ?? res.statusText;
  } catch {
    return res.statusText;
  }
}

export async function createAnalysis(repoUrl: string): Promise<{ analysis_id: string; status: string; repository: string }> {
  const res = await fetch(`${API}/api/analyses`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ repo_url: repoUrl }),
  });
  if (!res.ok) throw new Error(await errorMessage(res));
  return res.json();
}

export async function getAnalysis(id: string): Promise<AnalysisResponse> {
  const res = await fetch(`${API}/api/analyses/${id}`);
  if (!res.ok) throw new Error(await errorMessage(res));
  return res.json();
}

export async function askQuestion(id: string, question: string): Promise<AnswerResponse> {
  const res = await fetch(`${API}/api/analyses/${id}/questions`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ question }),
  });
  if (!res.ok) throw new Error(await errorMessage(res));
  return res.json();
}
