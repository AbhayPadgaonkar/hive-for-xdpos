"use client";

import { useEffect, useMemo, useState } from "react";
import { api, type Comparison } from "@/lib/api";
import { Badge, Button, Card, ProgressBar, SectionTitle, Skeleton } from "@/components/ui";
import { Filter, GitCompare, Search, TrendingUp, XCircle, CheckCircle } from "@/components/icons";
import Link from "next/link";

export default function ComparisonsPage() {
  const [comparisons, setComparisons] = useState<Comparison[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState("");
  const [showOnlyRegressions, setShowOnlyRegressions] = useState(false);

  useEffect(() => {
    api.comparisons()
      .then(setComparisons)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  const filtered = useMemo(() => {
    let rows = comparisons.filter((c) =>
      [c.client_a_name, c.client_b_name, c.simulator].some((s) =>
        s.toLowerCase().includes(filter.toLowerCase())
      )
    );
    if (showOnlyRegressions) rows = rows.filter((c) => c.a_only > 0 || c.b_only > 0);
    return rows;
  }, [comparisons, filter, showOnlyRegressions]);

  const totalRegressions = comparisons.reduce((sum, c) => sum + c.a_only + c.b_only, 0);
  const totalImprovements = comparisons.reduce((sum, c) => sum + c.b_only, 0);

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-100">Comparisons</h1>
          <p className="mt-1 max-w-2xl text-slate-400">
            Cross-client test suite comparisons showing regressions and improvements between runs.
          </p>
        </div>
        <div className="flex gap-3">
          <Link href="/gap-matrix">
            <Button variant="secondary" className="gap-2">
              <GitCompare size={16} /> Gap Matrix
            </Button>
          </Link>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-indigo-500/10 text-indigo-400">
              <GitCompare size={20} />
            </div>
            <div>
              <p className="text-2xl font-bold text-slate-100">{comparisons.length}</p>
              <p className="text-xs text-slate-500">Total Comparisons</p>
            </div>
          </div>
        </Card>
        <Card>
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-rose-500/10 text-rose-400">
              <XCircle size={20} />
            </div>
            <div>
              <p className="text-2xl font-bold text-slate-100">{totalRegressions}</p>
              <p className="text-xs text-slate-500">Regressions Found</p>
            </div>
          </div>
        </Card>
        <Card>
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-400">
              <TrendingUp size={20} />
            </div>
            <div>
              <p className="text-2xl font-bold text-slate-100">{totalImprovements}</p>
              <p className="text-xs text-slate-500">Improvements</p>
            </div>
          </div>
        </Card>
      </div>

      <Card>
        <div className="mb-4 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <SectionTitle title="Comparison Results" subtitle="Filtered and sorted by latest first" />
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div className="relative">
              <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" />
              <input
                type="text"
                placeholder="Search clients or simulator..."
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
                className="w-full rounded-lg border border-slate-700 bg-slate-800/50 py-2 pl-9 pr-4 text-sm text-slate-100 placeholder-slate-500 outline-none transition focus:border-indigo-500/50 focus:bg-slate-800 sm:w-64"
              />
            </div>
            <button
              onClick={() => setShowOnlyRegressions(!showOnlyRegressions)}
              className={`inline-flex items-center gap-2 rounded-lg border px-3 py-2 text-xs font-medium transition ${
                showOnlyRegressions
                  ? "border-rose-500/30 bg-rose-500/10 text-rose-300"
                  : "border-slate-700 bg-slate-800/50 text-slate-400 hover:text-slate-200"
              }`}
            >
              <Filter size={14} /> Regressions only
            </button>
          </div>
        </div>

        {loading ? (
          <div className="space-y-3">
            <Skeleton className="h-12" />
            <Skeleton className="h-12" />
            <Skeleton className="h-12" />
          </div>
        ) : filtered.length > 0 ? (
          <div className="overflow-hidden rounded-xl border border-slate-800">
            <table className="w-full text-sm">
              <thead className="bg-slate-900/80 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                <tr>
                  <th className="px-4 py-3">Simulator</th>
                  <th className="px-4 py-3">Client A</th>
                  <th className="px-4 py-3">Client B</th>
                  <th className="px-4 py-3 text-right">Both Pass</th>
                  <th className="px-4 py-3 text-right">Regressions</th>
                  <th className="px-4 py-3 text-right">Both Fail</th>
                  <th className="px-4 py-3 text-center">Health</th>
                  <th className="px-4 py-3">Date</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-800">
                {filtered.map((c) => {
                  const hasRegression = c.a_only > 0 || c.b_only > 0;
                  const total = c.total_matched || 1;
                  const passPct = ((c.both_pass / total) * 100);
                  return (
                    <tr key={c.id} className="transition hover:bg-slate-800/40">
                      <td className="px-4 py-3 font-medium text-slate-200">{c.simulator}</td>
                      <td className="px-4 py-3 text-slate-400">{c.client_a_name}</td>
                      <td className="px-4 py-3 text-slate-400">{c.client_b_name}</td>
                      <td className="px-4 py-3 text-right font-medium text-emerald-400">{c.both_pass}</td>
                      <td className="px-4 py-3 text-right">
                        {hasRegression ? (
                          <span className="font-medium text-rose-400">{c.a_only + c.b_only}</span>
                        ) : (
                          <span className="text-slate-500">0</span>
                        )}
                      </td>
                      <td className="px-4 py-3 text-right text-slate-400">{c.both_fail}</td>
                      <td className="px-4 py-3">
                        <div className="flex items-center justify-center">
                          {hasRegression ? (
                            <Badge tone="danger"><XCircle size={12} /> Regression</Badge>
                          ) : (
                            <Badge tone="success"><CheckCircle size={12} /> Aligned</Badge>
                          )}
                        </div>
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-3">
                          <div className="w-24">
                            <ProgressBar value={passPct} tone={hasRegression ? "warning" : "success"} />
                          </div>
                          <span className="w-10 text-right text-xs text-slate-400">{passPct.toFixed(0)}%</span>
                        </div>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-slate-800 py-14 text-center">
            <GitCompare size={40} className="text-slate-600" />
            <p className="mt-4 text-sm font-medium text-slate-400">No comparisons match your filters.</p>
            <p className="text-xs text-slate-500">Run \u003ccode\u003ecompare.ps1 diff\u003c/code\u003e to generate comparison data.</p>
          </div>
        )}
      </Card>
    </div>
  );
}
