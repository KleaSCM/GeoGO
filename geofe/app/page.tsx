"use client";
import { useState } from "react";
import MeteoriteList from "../components/MeteoriteList";
import SearchForm from "../components/SearchForm";

export default function HomePage() {
  const [results, setResults] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const searchMeteorites = async (query: Record<string, string | number>) => {
    setLoading(true);
    setError(null);
    setResults([]); 

    try {
      const BASE_URL =
        process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "") ||
        "http://localhost:8080";
      const url = new URL(`${BASE_URL}/meteorites`);

      Object.entries(query).forEach(([key, value]) => {
        if (value !== "" && value !== undefined) {
          url.searchParams.append(key, String(value).trim());
        }
      });

      console.log("🔍 Query Parameters:", query);
      console.log("🌍 Fetching from:", url.toString());

      const response = await fetch(url.toString());

      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      const data = await response.json();
      console.log("✅ API Response:", data);

      // Ensure the API returns an array
      if (!Array.isArray(data)) {
        throw new Error(
          `Unexpected API response format: Expected an array, received ${typeof data}`
        );
      }

      // Ensure objects contain valid lat/lon
      data.forEach((rock, index) => {
        if (
          typeof rock.lat !== "number" ||
          typeof rock.lon !== "number" ||
          isNaN(rock.lat) ||
          isNaN(rock.lon)
        ) {
          console.warn(`⚠️ Warning: Invalid lat/lon for item #${index}`, rock);
        }
      });

      setResults(data);
    } catch (error) {
      console.error("❌ Error fetching meteorites:", error);
      setError(
        error instanceof Error ? error.message : "An unexpected error occurred."
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-5xl mx-auto p-6">
      <SearchForm onSearch={searchMeteorites} loading={loading} />
      {error && <p className="text-red-500 text-center mt-4">{error}</p>}
      <MeteoriteList results={results} />
    </div>
  );
}
