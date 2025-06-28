"use client";
import { useState } from "react";
import styles from "./SearchForm.module.scss";

interface SearchProps {
  onSearch: (query: Record<string, string>) => void;
  loading: boolean;
  datasetType?: string;
}

export default function SearchForm({ onSearch, loading, datasetType = "meteorite" }: SearchProps) {
  const [query, setQuery] = useState({
    year_start: "1900",
    year_end: "2025",
    recclass: "",
    mass_min: "",
    mass_max: "",
    location: "",
    value_min: "",
    value_max: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setQuery({ ...query, [e.target.name]: e.target.value });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log("Submitting query:", query);
    onSearch(query);
  };

  const getSearchFields = () => {
    switch (datasetType) {
      case "meteorite":
        return [
          { name: "year_start", placeholder: "YEAR START", type: "text" },
          { name: "year_end", placeholder: "YEAR END", type: "text" },
          { name: "recclass", placeholder: "METEORITE CLASS", type: "text" },
          { name: "mass_min", placeholder: "MASS MIN (g)", type: "text" },
          { name: "mass_max", placeholder: "MASS MAX (g)", type: "text" },
          { name: "location", placeholder: "LOCATION", type: "text" },
        ];
      case "climate":
        return [
          { name: "value_min", placeholder: "TEMP MIN (°C)", type: "text" },
          { name: "value_max", placeholder: "TEMP MAX (°C)", type: "text" },
          { name: "location", placeholder: "LOCATION", type: "text" },
        ];
      case "wind":
        return [
          { name: "value_min", placeholder: "WIND SPEED MIN (m/s)", type: "text" },
          { name: "value_max", placeholder: "WIND SPEED MAX (m/s)", type: "text" },
          { name: "location", placeholder: "LOCATION", type: "text" },
        ];
      case "vegetation":
        return [
          { name: "value_min", placeholder: "AREA MIN (m²)", type: "text" },
          { name: "value_max", placeholder: "AREA MAX (m²)", type: "text" },
          { name: "location", placeholder: "LOCATION", type: "text" },
        ];
      default:
        return [
          { name: "value_min", placeholder: "VALUE MIN", type: "text" },
          { name: "value_max", placeholder: "VALUE MAX", type: "text" },
          { name: "location", placeholder: "LOCATION", type: "text" },
        ];
    }
  };

  const getTitle = () => {
    switch (datasetType) {
      case "meteorite":
        return "☄️ Search Meteorites";
      case "climate":
        return "🌡️ Search Climate Data";
      case "wind":
        return "💨 Search Wind Data";
      case "vegetation":
        return "🌿 Search Vegetation Data";
      case "infrastructure":
        return "🏗️ Search Infrastructure Data";
      case "fire":
        return "🔥 Search Fire Data";
      default:
        return "🔎 Search Data";
    }
  };

  const searchFields = getSearchFields();

  return (
    <div className={styles["search-container"]}>
      <h2 className="text-3xl font-bold text-blue-400">{getTitle()}</h2>
      <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4 mt-4">
        {searchFields.map((field) => (
          <input
            key={field.name}
            type={field.type}
            name={field.name}
            value={query[field.name as keyof typeof query] || ""}
            onChange={handleChange}
            placeholder={field.placeholder}
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
