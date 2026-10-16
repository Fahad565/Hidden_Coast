"use client";

import { useEffect, useState } from "react";

export default function Home() {
  const [health, setHealth] = useState<string>("Loading...");

  useEffect(() => {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    fetch(`${apiUrl}/api/health`)
      .then((res) => res.json())
      .then((data) => setHealth(data.status || "Unknown"))
      .catch(() => setHealth("Error connecting to backend"));
  }, []);

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 p-24 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
        <h1 className="text-3xl font-bold">Hidden Coast</h1>
        <p className="text-lg">Connecting tourists with Mombasa&apos;s local SMEs.</p>
        <div className="p-4 border rounded-lg bg-gray-50 dark:bg-gray-900">
          <span className="font-semibold">Backend Status: </span>
          <span className={health === "ok" ? "text-green-600 font-bold" : "text-amber-600"}>
            {health}
          </span>
        </div>
      </main>
    </div>
  );
}
