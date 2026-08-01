"use client";

import { useEffect, useMemo, useState } from "react";
import { api, type Run, type Stats, type GapMatrix, type Probe } from "@/lib/api";
import { Card, Badge, ProgressBar, SectionTitle, Skeleton, SparkLine } from "@/components/ui";
import {
  Activity, BarChart3, CheckCircle, ChevronRight, Cpu, GitCompare, Grid, Layers, Server, TrendingUp, XCircle, Zap,
} from "@/components/icons";
import Link from "next/link";

export default function Home() {
  const [runs, setRuns] = useState<Run[]>([]);
  const [stats, setStats] = useState<Stats | null>(null);
  const [matrices, setMatrices] = useState<GapMatrix[]>([]);
  const [probes, setProbes] = useState<Probe[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      api.runs(),
      api.stats(),
      api.gapMatrices(),
      api.probes(),
    ])
      .then(([r, s, m, p]) => {
        setRuns(r);
        setStats(s);
        setMatrices(m);
        setProbes(p);
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  const latestRuns = useMemo(() => runs.slice(0, 8), [runs]);
  const passRateHistory = useMemo(() =>
    runs.slice(0, 12).map((r) => (r.total > 0 ? (r.passed / r.total) * 100 : 0)).reverse(),
    [runs]
  );

  const matrix = matrices[0];

  return (
    <div className="space-y-8 animate-fade-in">
      <div className="relative overflow-hidden rounded-2xl border border-indigo-500/20 bg-gradient-to-br from-indigo-600/20 via-slate-900/80 to-slate-900/80 p-6 shadow-2xl lg:p-8">
        <div className="absolute -right-12 -top-12 h-64 w-64 rounded-full bg-indigo-500/10 blur-3xl" />
        <div className="relative">
          <div className="flex items-center gap-2">
            <Badge tone="primary"><Zap size={12} fill="currentColor" /> v0.1.0</Badge>
          </div>
          <h1 className="mt-4 text-3xl font-extrabold tracking-tight lg:text-4xl">
            Geth <span className="gradient-text">Comparison Framework</span>
          </h1>
          <p className="mt-2 max-w-2xl text-slate-400">
            Continuously probe, compare, and visualize how xdc-geth-audit diverges from upstream go-ethereum.
          </p>
          <div className="mt-6 flex flex-wrap gap-3">
            <Link
              href="/gap-matrix"
              className="inline-flex items-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow-lg shadow-indigo-500/25 transition hover:bg-indigo-500"
            >
              Explore Gap Matrix <ChevronRight size={16} />
            </Link>
            <Link
              href="/comparisons"
              className="inline-flex items-center gap-2 rounded-lg border border-slate-700 bg-slate-800/60 px-4 py-2 text-sm font-semibold text-slate-200 transition hover:bg-slate-700"
            >
              View Comparisons
            </Link>
          </div>
        </div>
      </div>

      {loading && !stats ? (
        <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
          <Skeleton className="h-28" />
          <Skeleton className="h-28" />
          <Skeleton className="h-28" />
          <Skeleton className="h-28" />
        </div>
      ) : (
        <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
          <StatCard
            label="Total Runs"
            value={stats?.total_runs ?? 0}
            icon={Activity}
            tone="primary"
            detail="Test executions"
          />
          <StatCard
            label="Method Probes"
            value={stats?.total_probes ?? 0}
            icon={Cpu}
            tone="info"
            detail="Supported methods"
          />
          <StatCard
            label="Comparisons"
            value={stats?.total_comparisons ?? 0}
            icon={GitCompare}
            tone="warning"
            detail="Cross-client diffs"
          />
          <StatCard
            label="Gap Matrices"
            value={stats?.total_gap_matrices ?? 0}
            icon={Grid}
            tone="success"
            detail="Feature coverage"
          />
        </div>
      )}

      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <SectionTitle title="Run History" subtitle="Recent simulator pass rates across clients" />
          {loading ? (
            <div className="space-y-3">
              <Skeleton className="h-12" />
              <Skeleton className="h-12" />
              <Skeleton className="h-12" />
              <Skeleton className="h-12" />
            </div>
          ) : latestRuns.length > 0 ? (
            <div className="overflow-hidden rounded-xl border border-slate-800">
              <table className="w-full text-sm">
                <thead className="bg-slate-900/80 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                  <tr>
                    <th className="px-4 py-3">Client</th>
                    <th className="px-4 py-3">Simulator</th>
                    <th className="px-4 py-3">Date</th>
                    <th className="px-4 py-3 text-right">Pass %</th>
                    <th className="px-4 py-3 text-right">Status</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800">
                  {latestRuns.map((run) => {
                    const pct = run.total > 0 ? (run.passed / run.total) * 100 : 0;
                    return (
                      <tr key={run.id} className="group transition hover:bg-slate-800/40">
                        <td className="px-4 py-3">
                          <div className="flex items-center gap-2">
                            <span className="flex h-6 w-6 items-center justify-center rounded-md bg-slate-800 text-slate-400">
                              <Server size={14} />
                            </span>
                            <span className="font-medium text-slate-200">{run.client_name}</span>
                          </div>
                        </td>
                        <td className="px-4 py-3 text-slate-400">{run.sim_name}</td>
                        <td className="px-4 py-3 text-slate-500">{run.date_run}</td>
                        <td className="px-4 py-3">
                          <div className="flex items-center gap-3">
                            <span className="w-10 text-right font-medium text-slate-200">{pct.toFixed(0)}%</span>
                            <div className="hidden w-24 sm:block">
                              <ProgressBar value={pct} tone={pct >= 90 ? "success" : pct >= 70 ? "warning" : "danger"} />
                            </div>
                          </div>
                        </td>
                        <td className="px-4 py-3 text-right">
                          {run.failed === 0 ? (
                            <Badge tone="success"><CheckCircle size={12} /> Clean</Badge>
                          ) : (
                            <Badge tone="danger"><XCircle size={12} /> {run.failed} failures</Badge>
                          )}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          ) : (
            <EmptyState icon={Layers} text="No runs yet. Start the API and run a Hive test." />
          )}
        </Card>

        <div className="space-y-6">
          <Card>
            <SectionTitle title="Pass Rate Trend" subtitle="Trailing run pass percentages" />
            <div className="rounded-xl bg-slate-900/60 p-4">
              <SparkLine data={passRateHistory} color="#22c55e" />
              <div className="mt-3 flex items-end justify-between">
                <div>
                  <p className="text-2xl font-bold text-slate-100">
                    {passRateHistory.length > 0 ? passRateHistory[passRateHistory.length - 1].toFixed(1) : "—"}%
                  </p>
                  <p className="text-xs text-slate-500">Latest run pass rate</p>
                </div>
                <TrendingUp size={20} className="text-emerald-400" />
              </div>
            </div>
          </Card>

          <Card>
            <SectionTitle title="Latest Gap Matrix" subtitle={matrix ? `${matrix.client_a_name} vs ${matrix.client_b_name}` : undefined} />
            {matrix ? (
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-3">
                  <div className="rounded-lg bg-slate-900/60 p-3">
                    <p className="text-xs text-slate-500">Both Supported</p>
                    <p className="text-xl font-bold text-emerald-400">{matrix.both_supported}</p>
                  </div>
                  <div className="rounded-lg bg-slate-900/60 p-3">
                    <p className="text-xs text-slate-500">Gaps</p>
                    <p className="text-xl font-bold text-amber-400">{matrix.in_a_not_b}</p>
                  </div>
                  <div className="rounded-lg bg-slate-900/60 p-3">
                    <p className="text-xs text-slate-500">Extras</p>
                    <p className="text-xl font-bold text-indigo-400">{matrix.in_b_not_a}</p>
                  </div>
                  <div className="rounded-lg bg-slate-900/60 p-3">
                    <p className="text-xs text-slate-500">Total Methods</p>
                    <p className="text-xl font-bold text-slate-200">{matrix.total_methods}</p>
                  </div>
                </div>
                <Link
                  href="/gap-matrix"
                  className="flex items-center justify-center gap-2 rounded-lg bg-slate-800 py-2 text-sm font-medium text-slate-300 transition hover:bg-slate-700 hover:text-slate-100"
                >
                  View full matrix <ChevronRight size={14} />
                </Link>
              </div>
            ) : (
              <EmptyState icon={Grid} text="No gap matrix available yet." />
            )}
          </Card>

          <Card>
            <SectionTitle title="Latest Probe Snapshot" subtitle="Method support coverage" />
            {probes[0] ? (
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-slate-400">{probes[0].client_name}</span>
                  <Badge tone="info">{probes[0].supported}/{probes[0].total} methods</Badge>
                </div>
                <ProgressBar
                  value={probes[0].total > 0 ? (probes[0].supported / probes[0].total) * 100 : 0}
                  tone="primary"
                />
                <div className="flex justify-between text-xs text-slate-500">
                  <span>Supported</span>
                  <span>Unsupported</span>
                </div>
              </div>
            ) : (
              <EmptyState icon={BarChart3} text="No probe data available yet." />
            )}
          </Card>
        </div>
      </div>
    </div>
  );
}

function StatCard({ label, value, icon: Icon, tone, detail }: {
  label: string; value: number; icon: React.ComponentType<{ size?: number; className?: string }>;
  tone: "primary" | "info" | "warning" | "success";
  detail: string;
}) {
  const tones = {
    primary: "text-indigo-400 bg-indigo-500/10",
    info: "text-sky-400 bg-sky-500/10",
    warning: "text-amber-400 bg-amber-500/10",
    success: "text-emerald-400 bg-emerald-500/10",
  }[tone];
  return (
    <Card className="relative overflow-hidden">
      <div className="flex items-start justify-between">
        <div className={`rounded-lg p-2 ${tones}`}>
          <Icon size={20} className={tones.split(" ")[0]} />
        </div>
      </div>
      <p className="mt-3 text-2xl font-bold text-slate-100">{value}</p>
      <p className="text-sm font-medium text-slate-300">{label}</p>
      <p className="mt-1 text-xs text-slate-500">{detail}</p>
    </Card>
  );
}

function EmptyState({ icon: Icon, text }: { icon: React.ComponentType<{ size?: number; className?: string }>; text: string }) {
  return (
    <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-slate-800 py-10 text-center">
      <Icon size={32} className="text-slate-600" />
      <p className="mt-3 text-sm text-slate-500">{text}</p>
    </div>
  );
}
