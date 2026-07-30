"use client";

import { useEffect, useState } from "react";
import { api, type Comparison } from "@/lib/api";

export default function ComparisonsPage() {
  const [comparisons, setComparisons] = useState<Comparison[]>([]);

  useEffect(() => {
    api.comparisons().then(setComparisons).catch(console.error);
  }, []);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Comparisons</h1>
      <p className="text-zinc-600">
        Cross-client test suite comparisons showing regressions and improvements.
      </p>

      {comparisons.length === 0 && (
        <p className="text-zinc-400">No comparisons available. Run <code>.\compare.ps1 diff</code> first.</p>
      )}

      <div className="overflow-x-auto rounded-lg border border-zinc-200">
        <table className="w-full text-sm">
          <thead className="bg-zinc-100 text-left">
            <tr>
              <th className="px-4 py-2 font-medium">Simulator</th>
              <th className="px-4 py-2 font-medium">Client A</th>
              <th className="px-4 py-2 font-medium">Client B</th>
              <th className="px-4 py-2 font-medium text-right">Both Pass</th>
              <th className="px-4 py-2 font-medium text-right">Regressions</th>
              <th className="px-4 py-2 font-medium text-right">Both Fail</th>
              <th className="px-4 py-2 font-medium">Date</th>
            </tr>
          </thead>
          <tbody>
            {comparisons.map((c) => (
              <tr key={c.id} className="border-t border-zinc-100 hover:bg-zinc-50">
                <td className="px-4 py-2 font-medium">{c.simulator}</td>
                <td className="px-4 py-2">{c.client_a_name}</td>
                <td className="px-4 py-2">{c.client_b_name}</td>
                <td className="px-4 py-2 text-right text-green-600 font-medium">{c.both_pass}</td>
                <td className="px-4 py-2 text-right text-red-600 font-medium">{c.a_only}</td>
                <td className="px-4 py-2 text-right text-zinc-600">{c.both_fail}</td>
                <td className="px-4 py-2 text-zinc-500">{c.date_compared}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
