"use client";
import { useState } from "react";
import MeteoriteList from "../components/MeteoriteList";
import SearchForm from "../components/SearchForm";

export default function HomePage() {
  const [results, setResults] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  const searchMeteorites = async (query: Record<string, string>) => {
    setLoading(true);
    try {
      const BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
      const url = new URL(`${BASE_URL}/meteorites`);
      
      Object.entries(query).forEach(([key, value]) => {
        if (value.trim() !== "") url.searchParams.append(key, value);
      });

      console.log("Fetching from:", url.toString());
      const response = await fetch(url.toString());

      if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);

      const data = await response.json();
      setResults(data);
    } catch (error) {
      console.error("❌ Error fetching meteorites:", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-5xl mx-auto p-6">
      <SearchForm onSearch={searchMeteorites} loading={loading} />
      <MeteoriteList results={results} />
    </div>
  );
}
