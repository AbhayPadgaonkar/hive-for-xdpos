"use client";

import { useEffect, useState } from "react";
import { api, type GapMatrix, type GapMatrixDetail } from "@/lib/api";

function GapStatus({ supported }: { supported: boolean }) {
  return supported
    ? <span className="text-green-600 font-medium">Y</span>
    : <span className="text-red-400">N</span>;
}

export default function GapMatrixPage() {
  const [matrices, setMatrices] = useState<GapMatrix[]>([]);
  const [selected, setSelected] = useState<GapMatrixDetail | null>(null);
  const [detailId, setDetailId] = useState<number | null>(null);
  const [search, setSearch] = useState("");

  useEffect(() => {
    api.gapMatrices().then(setMatrices).catch(console.error);
  }, []);

  useEffect(() => {
    if (detailId === null) { setSelected(null); return }
    api.gapMatrix(detailId).then(setSelected).catch(console.error);
  }, [detailId]);

  const currentMatrix = matrices[0];

  const filteredMethods = selected?.methods.filter(m =>
    m.method.toLowerCase().includes(search.toLowerCase())
  ) ?? [];

  const gaps = filteredMethods.filter(m => m.a_supported && !m.b_supported);
  const extras = filteredMethods.filter(m => !m.a_supported && m.b_supported);
  const bothYes = filteredMethods.filter(m => m.a_supported && m.b_supported);
  const bothNo = filteredMethods.filter(m => !m.a_supported && !m.b_supported);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Feature Gap Matrix</h1>
      <p className="text-zinc-600">
        Compares RPC method support between two clients via method probing.
      </p>

      {matrices.length === 0 && (
        <p className="text-zinc-400">No gap matrices available.</p>
      )}

      {currentMatrix && !selected && (
        <div className="rounded-lg border border-zinc-200 bg-white">
          <div className="p-4 border-b border-zinc-100 bg-zinc-50 flex items-center justify-between">
            <h2 className="font-semibold">{currentMatrix.client_a_name} vs {currentMatrix.client_b_name}</h2>
            <span className="text-xs text-zinc-400">{currentMatrix.date_created}</span>
          </div>
          <div className="p-4">
            <div className="grid grid-cols-2 md:grid-cols-5 gap-4 text-sm mb-4">
              <div><span className="text-zinc-500">Total</span><p className="font-medium">{currentMatrix.total_methods}</p></div>
              <div><span className="text-green-600">Both ✓</span><p className="font-medium">{currentMatrix.both_supported}</p></div>
              <div><span className="text-red-600">Both ✗</span><p className="font-medium">{currentMatrix.both_unsupported}</p></div>
              <div><span className="text-amber-600">Gaps</span><p className="font-medium">{currentMatrix.in_a_not_b}</p></div>
              <div><span className="text-blue-600">Extras</span><p className="font-medium">{currentMatrix.in_b_not_a}</p></div>
            </div>
            <div className="grid grid-cols-2 gap-4 text-xs font-mono mb-4">
              <div><span className="font-semibold">{currentMatrix.client_a_name}</span><pre className="bg-zinc-50 p-2 rounded mt-1 overflow-auto max-h-24">{currentMatrix.modules_a ? JSON.stringify(JSON.parse(currentMatrix.modules_a), null, 2) : "—"}</pre></div>
              <div><span className="font-semibold">{currentMatrix.client_b_name}</span><pre className="bg-zinc-50 p-2 rounded mt-1 overflow-auto max-h-24">{currentMatrix.modules_b ? JSON.stringify(JSON.parse(currentMatrix.modules_b), null, 2) : "—"}</pre></div>
            </div>
            <button
              onClick={() => setDetailId(currentMatrix.id)}
              className="text-sm text-blue-600 hover:underline"
            >
              View method-level details →
            </button>
          </div>
        </div>
      )}

      {selected && (
        <div className="space-y-4">
          <button onClick={() => setDetailId(null)} className="text-sm text-blue-600 hover:underline">
            ← Back to summary
          </button>

          <div className="flex items-center gap-3">
            <input
              type="text"
              placeholder="Search methods..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="flex-1 border border-zinc-300 rounded px-3 py-1.5 text-sm"
            />
            <span className="text-xs text-zinc-400">{filteredMethods.length} methods</span>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-xs">
            <SummaryBlock title="Both Supported" color="text-green-600" count={bothYes.length} />
            <SummaryBlock title="Gaps (A only)" color="text-amber-600" count={gaps.length} />
            <SummaryBlock title="Both Unsupported" color="text-zinc-400" count={bothNo.length} />
          </div>

          <div className="overflow-x-auto rounded-lg border border-zinc-200 max-h-[600px] overflow-y-auto">
            <table className="w-full text-xs">
              <thead className="bg-zinc-100 sticky top-0">
                <tr>
                  <th className="px-3 py-1.5 text-left font-medium">Method</th>
                  <th className="px-3 py-1.5 text-center font-medium">{selected.client_a_name}</th>
                  <th className="px-3 py-1.5 text-center font-medium">{selected.client_b_name}</th>
                  <th className="px-3 py-1.5 text-left font-medium">Error</th>
                </tr>
              </thead>
              <tbody>
                {gaps.map(m => (
                  <tr key={m.method} className="border-t border-zinc-100 bg-amber-50/50">
                    <td className="px-3 py-1.5 font-mono">{m.method}</td>
                    <td className="px-3 py-1.5 text-center"><GapStatus supported={m.a_supported} /></td>
                    <td className="px-3 py-1.5 text-center"><GapStatus supported={m.b_supported} /></td>
                    <td className="px-3 py-1.5 text-zinc-500">{m.b_error || m.a_error || ""}</td>
                  </tr>
                ))}
                {bothYes.map(m => (
                  <tr key={m.method} className="border-t border-zinc-100">
                    <td className="px-3 py-1.5 font-mono">{m.method}</td>
                    <td className="px-3 py-1.5 text-center"><GapStatus supported={m.a_supported} /></td>
                    <td className="px-3 py-1.5 text-center"><GapStatus supported={m.b_supported} /></td>
                    <td className="px-3 py-1.5 text-zinc-500"></td>
                  </tr>
                ))}
                {bothNo.map(m => (
                  <tr key={m.method} className="border-t border-zinc-100 text-zinc-400">
                    <td className="px-3 py-1.5 font-mono">{m.method}</td>
                    <td className="px-3 py-1.5 text-center"><GapStatus supported={m.a_supported} /></td>
                    <td className="px-3 py-1.5 text-center"><GapStatus supported={m.b_supported} /></td>
                    <td className="px-3 py-1.5 text-zinc-400">{m.a_error || m.b_error || ""}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
}

function SummaryBlock({ title, color, count }: { title: string; color: string; count: number }) {
  return (
    <div className="rounded border border-zinc-200 p-3 bg-white">
      <div className={`font-semibold ${color}`}>{title}</div>
      <div className="text-2xl font-bold mt-1">{count}</div>
    </div>
  );
}
