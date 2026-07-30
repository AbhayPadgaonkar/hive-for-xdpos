"use client";

import { useEffect, useState } from "react";
import { api, type Run, type Stats, type GapMatrix } from "@/lib/api";

export default function Home() {
  const [runs, setRuns] = useState<Run[]>([]);
  const [stats, setStats] = useState<Stats | null>(null);
  const [matrices, setMatrices] = useState<GapMatrix[]>([]);
  const [expanded, setExpanded] = useState<number | null>(null);

  useEffect(() => {
    api.runs().then(setRuns).catch(console.error);
    api.stats().then(setStats).catch(console.error);
    api.gapMatrices().then(setMatrices).catch(console.error);
  }, []);

  return (
    <div className="space-y-8">
      <h1 className="text-2xl font-bold">Dashboard</h1>

      {stats && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <StatCard label="Runs" value={stats.total_runs} />
          <StatCard label="Probes" value={stats.total_probes} />
          <StatCard label="Comparisons" value={stats.total_comparisons} />
          <StatCard label="Gap Matrices" value={stats.total_gap_matrices} />
        </div>
      )}

      <section>
        <h2 className="text-lg font-semibold mb-3">Run Results</h2>
        <div className="overflow-x-auto rounded-lg border border-zinc-200">
          <table className="w-full text-sm">
            <thead className="bg-zinc-100 text-left">
              <tr>
                <th className="px-4 py-2 font-medium">Client</th>
                <th className="px-4 py-2 font-medium">Simulator</th>
                <th className="px-4 py-2 font-medium">Date</th>
                <th className="px-4 py-2 font-medium text-right">Total</th>
                <th className="px-4 py-2 font-medium text-right">Pass</th>
                <th className="px-4 py-2 font-medium text-right">Fail</th>
                <th className="px-4 py-2 font-medium text-right">Pass %</th>
              </tr>
            </thead>
            <tbody>
              {runs.map((run) => (
                <tr
                  key={run.id}
                  className="border-t border-zinc-100 hover:bg-zinc-50 cursor-pointer"
                  onClick={() => setExpanded(expanded === run.id ? null : run.id)}
                >
                  <td className="px-4 py-2 font-medium">{run.client_name}</td>
                  <td className="px-4 py-2 text-zinc-600">{run.sim_name}</td>
                  <td className="px-4 py-2 text-zinc-600">{run.date_run}</td>
                  <td className="px-4 py-2 text-right">{run.total}</td>
                  <td className="px-4 py-2 text-right text-green-600">{run.passed}</td>
                  <td className="px-4 py-2 text-right text-red-600">{run.failed}</td>
                  <td className="px-4 py-2 text-right font-medium">
                    {run.total > 0 ? ((run.passed / run.total) * 100).toFixed(1) : "-"}%
                  </td>
                </tr>
              ))}
              {runs.length === 0 && (
                <tr><td colSpan={7} className="px-4 py-8 text-center text-zinc-400">No runs yet</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </section>

      {matrices.length > 0 && (
        <section>
          <h2 className="text-lg font-semibold mb-3">Latest Feature Gap Matrix</h2>
          <div className="rounded-lg border border-zinc-200 p-4 bg-white">
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
              <div>
                <span className="text-zinc-500">Client A</span>
                <p className="font-medium">{matrices[0].client_a_name}</p>
              </div>
              <div>
                <span className="text-zinc-500">Client B</span>
                <p className="font-medium">{matrices[0].client_b_name}</p>
              </div>
              <div>
                <span className="text-zinc-500">Both Supported</span>
                <p className="font-medium text-green-600">{matrices[0].both_supported}</p>
              </div>
              <div>
                <span className="text-zinc-500">In A not B (gaps)</span>
                <p className="font-medium text-amber-600">{matrices[0].in_a_not_b}</p>
              </div>
            </div>
            <div className="mt-4">
              <a href="/gap-matrix" className="text-sm text-blue-600 hover:underline">
                View full matrix →
              </a>
            </div>
          </div>
        </section>
      )}
    </div>
  );
}

function StatCard({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-lg border border-zinc-200 bg-white p-4">
      <div className="text-xs text-zinc-500 uppercase tracking-wide">{label}</div>
      <div className="text-2xl font-bold mt-1">{value}</div>
    </div>
  );
}
