"use client"
import { useState } from "react";
import MeteoriteList from "../components/MeteoriteList";
export default function Home() {
  const [query, setQuery] = useState({
    year_start: "1900",
    year_end: "2025",
    recclass: "",
    mass_min: "",
    mass_max: "",
  });
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setQuery({ ...query, [e.target.name]: e.target.value });
  };
  const searchMeteorites = async () => {
    setLoading(true);
    try {
      const BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
      const url = new URL(`${BASE_URL}/meteorites`);
      Object.entries(query).forEach(([key, value]) => {
        if (value) url.searchParams.append(key, value);
      });
      console.log("Fetching from:", url.toString()); 
      const response = await fetch(url.toString());
      if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);

      const data = await response.json();
      setResults(data);
    } catch (error) {
      console.error("Error fetching meteorites:", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h2 className="text-2xl font-bold text-blue-400">Search Meteorites ☄️</h2>
      <div className="grid grid-cols-2 gap-4 mt-4">
        <input
          type="text"
          name="year_start"
          value={query.year_start}
          onChange={handleChange}
          placeholder="Start Year"
          className="p-2 bg-gray-700 text-white border border-gray-600 rounded"
        />
        <input
          type="text"
          name="year_end"
          value={query.year_end}
          onChange={handleChange}
          placeholder="End Year"
          className="p-2 bg-gray-700 text-white border border-gray-600 rounded"
        />
        <input
          type="text"
          name="recclass"
          value={query.recclass}
          onChange={handleChange}
          placeholder="Classification"
          className="p-2 bg-gray-700 text-white border border-gray-600 rounded"
        />
        <input
          type="text"
          name="mass_min"
          value={query.mass_min}
          onChange={handleChange}
          placeholder="Min Mass (g)"
          className="p-2 bg-gray-700 text-white border border-gray-600 rounded"
        />
        <input
          type="text"
          name="mass_max"
          value={query.mass_max}
          onChange={handleChange}
          placeholder="Max Mass (g)"
          className="p-2 bg-gray-700 text-white border border-gray-600 rounded"
        />
      </div>

      <button
        onClick={searchMeteorites}
        disabled={loading}
        className="mt-4 bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded"
      >
        {loading ? "Searching..." : "Search"}
      </button>

      <MeteoriteList results={results} />
    </div>
  );
}
