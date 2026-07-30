const BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function fetchJSON<T>(url: string): Promise<T> {
  const res = await fetch(url);
  if (!res.ok) throw new Error(`${url}: ${res.status}`);
  return res.json();
}

export interface Run {
  id: number; client_name: string; sim_name: string;
  date_run: string; total: number; passed: number; failed: number;
  version: string; created_at: string;
}

export interface Probe {
  id: number; client_name: string; version: string;
  supported: number; unsupported: number; total: number;
  modules: string; date_run: string;
}

export interface GapMatrix {
  id: number; client_a_name: string; client_b_name: string;
  version_a: string; version_b: string;
  total_methods: number; both_supported: number;
  both_unsupported: number; in_a_not_b: number; in_b_not_a: number;
  modules_a: string; modules_b: string;
  date_created: string;
}

export interface GapMatrixMethod {
  method: string; a_supported: boolean; b_supported: boolean;
  a_error?: string; b_error?: string;
}

export interface GapMatrixDetail extends GapMatrix {
  methods: GapMatrixMethod[];
}

export interface Comparison {
  id: number; simulator: string;
  client_a_name: string; client_b_name: string;
  both_pass: number; a_only: number; b_only: number; both_fail: number;
  total_matched: number; date_compared: string; created_at: string;
}

export interface ComparisonDetail extends Comparison {
  tests: { test_name: string; a_pass: boolean | null; b_pass: boolean | null; status: string }[];
}

export interface Stats {
  total_runs: number; total_probes: number; total_comparisons: number;
  total_gap_matrices: number; total_tests: number;
  xdc_geth_pass_rate: number; go_eth_pass_rate: number;
}

export const api = {
  runs: () => fetchJSON<Run[]>(`${BASE}/api/runs`),
  run: (id: number) => fetchJSON<Run>(`${BASE}/api/runs/${id}`),
  probes: () => fetchJSON<Probe[]>(`${BASE}/api/probes`),
  gapMatrices: () => fetchJSON<GapMatrix[]>(`${BASE}/api/gap-matrices`),
  gapMatrix: (id: number) => fetchJSON<GapMatrixDetail>(`${BASE}/api/gap-matrices/${id}`),
  comparisons: () => fetchJSON<Comparison[]>(`${BASE}/api/comparisons`),
  stats: () => fetchJSON<Stats>(`${BASE}/api/stats`),
};
