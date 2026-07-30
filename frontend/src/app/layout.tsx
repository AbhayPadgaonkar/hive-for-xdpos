import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import Link from "next/link";

const geistSans = Geist({ variable: "--font-geist-sans", subsets: ["latin"] });
const geistMono = Geist_Mono({ variable: "--font-geist-mono", subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Hive Dashboard",
  description: "Hive test framework results for xdc-geth-audit",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className={`${geistSans.variable} ${geistMono.variable}`}>
      <body className="min-h-screen bg-zinc-50 text-zinc-900 font-sans antialiased">
        <header className="border-b border-zinc-200 bg-white">
          <div className="mx-auto max-w-6xl flex items-center justify-between px-6 h-14">
            <Link href="/" className="text-lg font-bold tracking-tight hover:text-blue-600">
              Hive Dashboard
            </Link>
            <nav className="flex gap-6 text-sm font-medium">
              <Link href="/" className="hover:text-blue-600">Runs</Link>
              <Link href="/comparisons" className="hover:text-blue-600">Comparisons</Link>
              <Link href="/gap-matrix" className="hover:text-blue-600">Gap Matrix</Link>
            </nav>
          </div>
        </header>
        <main className="mx-auto max-w-6xl px-6 py-8">
          {children}
        </main>
      </body>
    </html>
  );
}
