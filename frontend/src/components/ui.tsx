import React from "react";

export function Card({ children, className = "", padding = "lg" }: {
  children: React.ReactNode;
  className?: string;
  padding?: "none" | "sm" | "md" | "lg";
}) {
  const p = {
    none: "",
    sm: "p-3",
    md: "p-4",
    lg: "p-5",
  }[padding];
  return (
    <div className={`card ${p} animate-fade-in ${className}`}>
      {children}
    </div>
  );
}

export function Badge({ children, tone = "default", className = "" }: {
  children: React.ReactNode;
  tone?: "success" | "danger" | "warning" | "info" | "default" | "primary";
  className?: string;
}) {
  const styles = {
    success: "bg-emerald-500/15 text-emerald-400 border-emerald-500/20",
    danger: "bg-rose-500/15 text-rose-400 border-rose-500/20",
    warning: "bg-amber-500/15 text-amber-400 border-amber-500/20",
    info: "bg-sky-500/15 text-sky-400 border-sky-500/20",
    default: "bg-slate-500/15 text-slate-300 border-slate-500/20",
    primary: "bg-indigo-500/15 text-indigo-300 border-indigo-500/20",
  }[tone];
  return (
    <span className={`inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium ${styles} ${className}`}>
      {children}
    </span>
  );
}

export function Button({ children, onClick, variant = "primary", className = "", type = "button" }: {
  children: React.ReactNode;
  onClick?: () => void;
  variant?: "primary" | "secondary" | "ghost" | "danger";
  className?: string;
  type?: "button" | "submit";
}) {
  const styles = {
    primary: "bg-indigo-600 hover:bg-indigo-500 text-white shadow-lg shadow-indigo-500/20",
    secondary: "bg-slate-800 hover:bg-slate-700 text-slate-100 border border-slate-700",
    ghost: "hover:bg-slate-800 text-slate-300",
    danger: "bg-rose-600 hover:bg-rose-500 text-white",
  }[variant];
  return (
    <button
      type={type}
      onClick={onClick}
      className={`inline-flex items-center justify-center gap-2 rounded-lg px-4 py-2 text-sm font-medium transition-all active:scale-95 ${styles} ${className}`}
    >
      {children}
    </button>
  );
}

export function StatusDot({ tone }: { tone: "success" | "danger" | "warning" | "info" | "default" }) {
  const colors = {
    success: "bg-emerald-400",
    danger: "bg-rose-400",
    warning: "bg-amber-400",
    info: "bg-sky-400",
    default: "bg-slate-400",
  }[tone];
  return <span className={`inline-block h-2 w-2 rounded-full ${colors}`} />;
}

export function ProgressBar({ value, tone = "primary" }: { value: number; tone?: "primary" | "success" | "danger" | "warning" }) {
  const colors = {
    primary: "bg-indigo-500",
    success: "bg-emerald-500",
    danger: "bg-rose-500",
    warning: "bg-amber-500",
  }[tone];
  const pct = Math.max(0, Math.min(100, value));
  return (
    <div className="h-2 w-full overflow-hidden rounded-full bg-slate-800">
      <div
        className={`h-full rounded-full transition-all duration-700 ${colors}`}
        style={{ width: `${pct}%` }}
      />
    </div>
  );
}

export function Skeleton({ className = "" }: { className?: string }) {
  return <div className={`animate-pulse rounded-lg bg-slate-800 ${className}`} />;
}

export function SparkLine({ data, color = "#818cf8" }: { data: number[]; color?: string }) {
  if (!data.length) return <div className="h-10 w-full rounded bg-slate-800/50" />;
  const min = Math.min(...data);
  const max = Math.max(...data);
  const range = max - min || 1;
  const points = data.map((v, i) => {
    const x = (i / (data.length - 1 || 1)) * 100;
    const y = 100 - ((v - min) / range) * 100;
    return `${x},${y}`;
  }).join(" ");
  return (
    <svg viewBox="0 0 100 100" preserveAspectRatio="none" className="h-10 w-full overflow-visible">
      <polyline
        points={points}
        fill="none"
        stroke={color}
        strokeWidth="3"
        strokeLinecap="round"
        strokeLinejoin="round"
        vectorEffect="non-scaling-stroke"
      />
      <circle cx={points.split(" ").pop()?.split(",")[0]} cy={points.split(" ").pop()?.split(",")[1]} r="3" fill={color} />
    </svg>
  );
}

export function SectionTitle({ title, subtitle }: { title: string; subtitle?: string }) {
  return (
    <div className="mb-5">
      <h2 className="text-lg font-semibold tracking-tight text-slate-100">{title}</h2>
      {subtitle && <p className="text-sm text-slate-400">{subtitle}</p>}
    </div>
  );
}
