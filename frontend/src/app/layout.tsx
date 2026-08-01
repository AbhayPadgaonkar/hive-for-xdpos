import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { Sidebar, TopBar } from "@/components/sidebar";

const geistSans = Geist({ variable: "--font-geist-sans", subsets: ["latin"] });
const geistMono = Geist_Mono({ variable: "--font-geist-mono", subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Hive Dashboard",
  description: "Hive test framework results for xdc-geth-audit",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className={`${geistSans.variable} ${geistMono.variable}`}>
      <body className="min-h-screen antialiased">
        <Sidebar />
        <div className="flex min-h-screen flex-col lg:pl-64">
          <TopBar />
          <main className="flex-1 p-4 lg:p-8">
            {children}
          </main>
        </div>
      </body>
    </html>
  );
}
