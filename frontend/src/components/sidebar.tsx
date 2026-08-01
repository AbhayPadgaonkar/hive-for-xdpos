"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState } from "react";
import {
  LayoutDashboard, GitCompare, Grid, Menu, X, Zap,
} from "@/components/icons";

const nav = [
  { href: "/", label: "Dashboard", icon: LayoutDashboard },
  { href: "/comparisons", label: "Comparisons", icon: GitCompare },
  { href: "/gap-matrix", label: "Gap Matrix", icon: Grid },
];

export function Sidebar() {
  const [open, setOpen] = useState(false);
  const pathname = usePathname();

  return (
    <>
      <button
        onClick={() => setOpen(!open)}
        className="fixed left-4 top-3 z-50 rounded-lg bg-slate-800/80 p-2 text-slate-300 backdrop-blur-md transition hover:bg-slate-700 lg:hidden"
        aria-label="Toggle menu"
      >
        {open ? <X size={20} /> : <Menu size={20} />}
      </button>

      <aside
        className={`fixed inset-y-0 left-0 z-40 w-64 transform border-r border-slate-800 bg-slate-950/95 backdrop-blur-xl transition-transform duration-300 lg:translate-x-0 ${open ? "translate-x-0" : "-translate-x-full"}`}
      >
        <div className="flex h-16 items-center gap-3 border-b border-slate-800 px-6">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-indigo-600 text-white shadow-lg shadow-indigo-500/25">
            <Zap size={18} fill="currentColor" />
          </div>
          <div>
            <Link href="/" className="text-base font-bold tracking-tight text-slate-100">
              Hive Dashboard
            </Link>
            <p className="text-[10px] font-medium uppercase tracking-widest text-slate-500">Test Ops Center</p>
          </div>
        </div>

        <nav className="flex flex-col gap-1 p-4">
          {nav.map((item) => {
            const active = pathname === item.href || pathname.startsWith(`${item.href}/`);
            const Icon = item.icon;
            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={() => setOpen(false)}
                className={`flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-all ${
                  active
                    ? "bg-indigo-500/10 text-indigo-300 shadow-[inset_0_1px_0_rgba(99,102,241,0.1)]"
                    : "text-slate-400 hover:bg-slate-800/60 hover:text-slate-100"
                }`}
              >
                <Icon size={18} className={active ? "text-indigo-400" : "text-slate-500"} />
                {item.label}
                {active && <span className="ml-auto h-1.5 w-1.5 rounded-full bg-indigo-400" />}
              </Link>
            );
          })}
        </nav>

        <div className="absolute bottom-0 left-0 right-0 border-t border-slate-800 p-4">
          <div className="rounded-lg bg-slate-900/80 p-3">
            <p className="text-xs font-medium text-slate-300">Local framework</p>
            <p className="mt-0.5 text-[10px] text-slate-500">Go backend + SQLite + Next.js</p>
          </div>
        </div>
      </aside>

      {open && (
        <div
          className="fixed inset-0 z-30 bg-black/60 backdrop-blur-sm lg:hidden"
          onClick={() => setOpen(false)}
        />
      )}
    </>
  );
}

export function TopBar() {
  const pathname = usePathname();
  const label = nav.find((n) => pathname === n.href || pathname.startsWith(`${n.href}/`))?.label || "Dashboard";
  return (
    <header className="glass sticky top-0 z-20 flex h-16 items-center justify-between px-6">
      <h1 className="hidden text-sm font-medium text-slate-400 lg:block">{label}</h1>
      <div className="flex items-center gap-3">
        <span className="inline-flex h-2 w-2 rounded-full bg-emerald-400 animate-pulse" />
        <span className="text-xs font-medium text-slate-400">API Online</span>
      </div>
    </header>
  );
}
