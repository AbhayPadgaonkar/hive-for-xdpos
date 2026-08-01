"use client";

import { useEffect, useMemo, useState } from "react";
import { api, type GapMatrix, type GapMatrixDetail } from "@/lib/api";
import { Badge, Button, Card, ProgressBar, SectionTitle, Skeleton } from "@/components/ui";
import { ChevronLeft, Grid, Search, CheckCircle, XCircle, Filter } from "@/components/icons";

function GapStatus({ supported }: { supported: boolean }) {
  return supported
    ? <span className="inline-flex items-center gap-1 text-xs font-medium text-emerald-400"><CheckCircle size={12} /> Yes</span>
    : <span className="inline-flex items-center gap-1 text-xs font-medium text-rose-400"><XCircle size={12} /> No</span>;
}

export default function GapMatrixPage() {
  const [matrices, setMatrices] = useState<GapMatrix[]>([]);
  const [selected, setSelected] = useState<GapMatrixDetail | null>(null);
  const [detailId, setDetailId] = useState<number | null>(null);
  const [search, setSearch] = useState("");
  const [tab, setTab] = useState<"all" | "gaps" | "extras" | "both" | "none">("all");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.gapMatrices()
      .then(setMatrices)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (detailId === null) { setSelected(null); return; }
    api.gapMatrix(detailId).then(setSelected).catch(console.error);
  }, [detailId]);

  const currentMatrix = matrices[0];

  const filteredMethods = useMemo(() => {
    if (!selected) return [];
    const q = search.toLowerCase();
    return selected.methods.filter((m) =>
      m.method.toLowerCase().includes(q) ||
      (m.a_error || "").toLowerCase().includes(q) ||
      (m.b_error || "").toLowerCase().includes(q)
    ).filter((m) => {
      if (tab === "all") return true;
      if (tab === "gaps") return m.a_supported && !m.b_supported;
      if (tab === "extras") return !m.a_supported && m.b_supported;
      if (tab === "both") return m.a_supported && m.b_supported;
      if (tab === "none") return !m.a_supported && !m.b_supported;
      return true;
    });
  }, [selected, search, tab]);

  const counts = useMemo(() => {
    if (!selected) return { gaps: 0, extras: 0, both: 0, none: 0, total: 0 };
    const gaps = selected.methods.filter((m) => m.a_supported && !m.b_supported).length;
    const extras = selected.methods.filter((m) => !m.a_supported && m.b_supported).length;
    const both = selected.methods.filter((m) => m.a_supported && m.b_supported).length;
    const none = selected.methods.filter((m) => !m.a_supported && !m.b_supported).length;
    return { gaps, extras, both, none, total: selected.methods.length };
  }, [selected]);

  return (
    <div className="space-y-6 animate-fade-in">
      <div className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-100">Feature Gap Matrix</h1>
          <p className="mt-1 max-w-2xl text-slate-400">
            Compares RPC method support between two clients via method probing. Gaps are methods supported by client A but not client B.
          </p>
        </div>
      </div>

      {loading ? (
        <div className="grid gap-4 md:grid-cols-2">
          <Skeleton className="h-64" />
          <Skeleton className="h-64" />
        </div>
      ) : currentMatrix ? (
        <>
          {!selected ? (
            <div className="grid gap-6 lg:grid-cols-3">
              <Card className="lg:col-span-2">
                <SectionTitle
                  title="Matrix Summary"
                  subtitle={`${currentMatrix.client_a_name} vs ${currentMatrix.client_b_name} • ${currentMatrix.date_created}`}
                />

                <div className="grid grid-cols-2 gap-3 sm:grid-cols-5">
                  <Metric label="Total" value={currentMatrix.total_methods} color="text-slate-200" />
                  <Metric label="Both ✓" value={currentMatrix.both_supported} color="text-emerald-400" />
                  <Metric label="Both ✗" value={currentMatrix.both_unsupported} color="text-slate-400" />
                  <Metric label="Gaps" value={currentMatrix.in_a_not_b} color="text-amber-400" />
                  <Metric label="Extras" value={currentMatrix.in_b_not_a} color="text-indigo-400" />
                </div>

                <div className="mt-6 grid gap-4 md:grid-cols-2">
                  <div className="rounded-xl border border-slate-800 bg-slate-900/50 p-4">
                    <p className="mb-2 text-xs font-semibold uppercase tracking-wider text-slate-500">{currentMatrix.client_a_name}</p>
                    <pre className="max-h-48 overflow-auto rounded-lg bg-slate-950 p-3 text-xs text-slate-400">
                      {currentMatrix.modules_a ? JSON.stringify(JSON.parse(currentMatrix.modules_a), null, 2) : "—"}
                    </pre>
                  </div>
                  <div className="rounded-xl border border-slate-800 bg-slate-900/50 p-4">
                    <p className="mb-2 text-xs font-semibold uppercase tracking-wider text-slate-500">{currentMatrix.client_b_name}</p>
                    <pre className="max-h-48 overflow-auto rounded-lg bg-slate-950 p-3 text-xs text-slate-400">
                      {currentMatrix.modules_b ? JSON.stringify(JSON.parse(currentMatrix.modules_b), null, 2) : "—"}
                    </pre>
                  </div>
                </div>

                <div className="mt-6 flex justify-end">
                  <Button onClick={() => setDetailId(currentMatrix.id)} className="gap-2">
                    View Method Details <ChevronLeft className="rotate-180" size={16} />
                  </Button>
                </div>
              </Card>

              <div className="space-y-6">
                <Card>
                  <SectionTitle title="Coverage Breakdown" />
                  <div className="space-y-4">
                    <Breakdown label="Both Supported" count={currentMatrix.both_supported} total={currentMatrix.total_methods} color="emerald" />
                    <Breakdown label="Gaps (A only)" count={currentMatrix.in_a_not_b} total={currentMatrix.total_methods} color="amber" />
                    <Breakdown label="Extras (B only)" count={currentMatrix.in_b_not_a} total={currentMatrix.total_methods} color="indigo" />
                    <Breakdown label="Both Unsupported" count={currentMatrix.both_unsupported} total={currentMatrix.total_methods} color="slate" />
                  </div>
                </Card>

                <Card>
                  <SectionTitle title="Interpretation" />
                  <div className="space-y-3 text-sm text-slate-400">
                    <p>
                      <span className="font-semibold text-amber-400">Gaps:</span> Methods supported by upstream geth but missing in xdc-geth-audit. These are compatibility risks.
                    </p>
                    <p>
                      <span className="font-semibold text-indigo-400">Extras:</span> Methods added by XDPoS that are not present in upstream geth.
                    </p>
                    <p>
                      <span className="font-semibold text-emerald-400">Both Supported:</span> Shared compatibility surface.
                    </p>
                  </div>
                </Card>
              </div>
            </div>
          ) : (
            <Card>
              <div className="mb-5 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                <div className="flex items-center gap-3">
                  <Button variant="ghost" onClick={() => setDetailId(null)} className="gap-2 px-2">
                    <ChevronLeft size={16} /> Back
                  </Button>
                  <SectionTitle
                    title="Method-Level Details"
                    subtitle={`${selected.client_a_name} vs ${selected.client_b_name}`}
                  />
                </div>

                <div className="relative">
                  <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" />
                  <input
                    type="text"
                    placeholder="Search methods..."
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    className="w-full rounded-lg border border-slate-700 bg-slate-800/50 py-2 pl-9 pr-4 text-sm text-slate-100 placeholder-slate-500 outline-none transition focus:border-indigo-500/50 focus:bg-slate-800 md:w-72"
                  />
                </div>
              </div>

              <div className="mb-5 flex flex-wrap gap-2">
                <TabBadge active={tab === "all"} count={counts.total} onClick={() => setTab("all")}>All</TabBadge>
                <TabBadge active={tab === "gaps"} count={counts.gaps} color="amber" onClick={() => setTab("gaps")}>Gaps</TabBadge>
                <TabBadge active={tab === "extras"} count={counts.extras} color="indigo" onClick={() => setTab("extras")}>Extras</TabBadge>
                <TabBadge active={tab === "both"} count={counts.both} color="emerald" onClick={() => setTab("both")}>Both Supported</TabBadge>
                <TabBadge active={tab === "none"} count={counts.none} color="slate" onClick={() => setTab("none")}>Both Unsupported</TabBadge>
              </div>

              <div className="overflow-hidden rounded-xl border border-slate-800">
                <table className="w-full text-sm">
                  <thead className="bg-slate-900/80 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">
                    <tr>
                      <th className="px-4 py-3">Method</th>
                      <th className="px-4 py-3 text-center">{selected.client_a_name}</th>
                      <th className="px-4 py-3 text-center">{selected.client_b_name}</th>
                      <th className="px-4 py-3">Error / Note</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-800">
                    {filteredMethods.map((m) => (
                      <tr
                        key={m.method}
                        className={`transition hover:bg-slate-800/40 ${
                          m.a_supported && !m.b_supported ? "bg-amber-500/5" : ""
                        } ${!m.a_supported && m.b_supported ? "bg-indigo-500/5" : ""}`}
                      >
                        <td className="px-4 py-3 font-mono text-xs text-slate-200">{m.method}</td>
                        <td className="px-4 py-3 text-center"><GapStatus supported={m.a_supported} /></td>
                        <td className="px-4 py-3 text-center"><GapStatus supported={m.b_supported} /></td>
                        <td className="px-4 py-3 text-xs text-slate-500">
                          {m.b_error || m.a_error || "—"}
                        </td>
                      </tr>
                    ))}
                    {filteredMethods.length === 0 && (
                      <tr>
                        <td colSpan={4} className="px-4 py-12 text-center text-slate-500">
                          <div className="flex flex-col items-center gap-2">
                            <Filter size={24} className="text-slate-600" />
                            <p>No methods match the current filter.</p>
                          </div>
                        </td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            </Card>
          )}
        </>
      ) : (
        <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-slate-800 py-20 text-center">
          <Grid size={48} className="text-slate-600" />
          <p className="mt-4 text-sm font-medium text-slate-400">No gap matrices available.</p>
          <p className="text-xs text-slate-500">Run \u003ccode\u003ecompare.ps1 gap-matrix\u003c/code\u003e to generate one.</p>
        </div>
      )}
    </div>
  );
}

function Metric({ label, value, color }: { label: string; value: number; color: string }) {
  return (
    <div className="rounded-xl border border-slate-800 bg-slate-900/50 p-3 text-center">
      <p className={`text-2xl font-bold ${color}`}>{value}</p>
      <p className="mt-1 text-[10px] font-semibold uppercase tracking-wider text-slate-500">{label}</p>
    </div>
  );
}

function Breakdown({ label, count, total, color }: { label: string; count: number; total: number; color: "emerald" | "amber" | "indigo" | "slate" }) {
  const pct = total > 0 ? (count / total) * 100 : 0;
  const tones = {
    emerald: { bar: "bg-emerald-500", text: "text-emerald-400" },
    amber: { bar: "bg-amber-500", text: "text-amber-400" },
    indigo: { bar: "bg-indigo-500", text: "text-indigo-400" },
    slate: { bar: "bg-slate-500", text: "text-slate-400" },
  }[color];
  return (
    <div>
      <div className="mb-1 flex items-center justify-between text-sm">
        <span className="text-slate-400">{label}</span>
        <span className={`font-semibold ${tones.text}`}>{count} ({pct.toFixed(1)}%)</span>
      </div>
      <div className="h-2 overflow-hidden rounded-full bg-slate-800">
        <div className={`h-full rounded-full transition-all ${tones.bar}`} style={{ width: `${pct}%` }} />
      </div>
    </div>
  );
}

function TabBadge({ children, active, count, color = "slate", onClick }: {
  children: React.ReactNode; active: boolean; count: number; color?: "slate" | "amber" | "indigo" | "emerald"; onClick: () => void;
}) {
  const styles = {
    slate: active ? "bg-slate-700 text-slate-100 border-slate-600" : "text-slate-400 hover:bg-slate-800/60 border-slate-800",
    amber: active ? "bg-amber-500/15 text-amber-300 border-amber-500/30" : "text-slate-400 hover:bg-slate-800/60 border-slate-800",
    indigo: active ? "bg-indigo-500/15 text-indigo-300 border-indigo-500/30" : "text-slate-400 hover:bg-slate-800/60 border-slate-800",
    emerald: active ? "bg-emerald-500/15 text-emerald-300 border-emerald-500/30" : "text-slate-400 hover:bg-slate-800/60 border-slate-800",
  }[color];
  return (
    <button
      onClick={onClick}
      className={`inline-flex items-center gap-2 rounded-lg border px-3 py-1.5 text-xs font-medium transition ${styles}`}
    >
      {children}
      <span className={`rounded-md px-1.5 py-0.5 text-[10px] ${active ? "bg-slate-900/40 text-slate-300" : "bg-slate-900/50 text-slate-500"}`}>{count}</span>
    </button>
  );
}
