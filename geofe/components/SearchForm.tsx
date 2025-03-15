"use client";
import { useState } from "react";
import styles from "./SearchForm.module.scss";

interface SearchProps {
  onSearch: (query: Record<string, string>) => void;
  loading: boolean;
}

export default function SearchForm({ onSearch, loading }: SearchProps) {
  const [query, setQuery] = useState({
    year_start: "1900",
    year_end: "2025",
    recclass: "",
    mass_min: "",
    mass_max: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setQuery({ ...query, [e.target.name]: e.target.value });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log("Submitting query:", query);  
    onSearch(query);
};


  return (
    <div className={styles["search-container"]}>
      <h2 className="text-3xl font-bold text-blue-400">🔎 Search Meteorites ☄️</h2>
      <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4 mt-4">
        {Object.entries(query).map(([key, value]) => (
          <input
            key={key}
            type="text"
            name={key}
            value={value}
            onChange={handleChange}
            placeholder={key.replace("_", " ").toUpperCase()}
            className={styles["search-input"]}
          />
        ))}
        <button type="submit" disabled={loading} className={styles["search-button"]}>
          {loading ? "Searching..." : "Search"}
        </button>
      </form>
    </div>
  );
}
