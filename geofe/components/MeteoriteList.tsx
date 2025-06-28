"use client";
import dynamic from "next/dynamic";
import Geocode from "./Geocode";
import RockImage from "./RockImage";
import { useMemo } from "react";
import styles from "./MeteoriteList.module.scss";

// Lazy load Leaflet map (disable SSR)
const Map = dynamic(() => import("./Map"), { ssr: false });

interface Meteorite {
  id: number;
  name: string;
  recclass?: string | null;
  mass?: number | null;
  year?: number | null;
  lat?: number; 
  lon?: number;
  value?: number | null;
  unit?: string | null;
  metadata?: string | null;
}

interface MeteoriteListProps {
  results: Meteorite[];
  datasetType?: string;
}

export default function MeteoriteList({ results, datasetType = "meteorite" }: MeteoriteListProps) {
  if (!Array.isArray(results) || results.length === 0) {
    return <p className="text-gray-400 text-center mt-6">No data found.</p>;
  }

  // Clean the results to ensure all values are renderable
  const cleanResults = results.map(item => ({
    ...item,
    name: String(item.name || "Unknown"),
    recclass: item.recclass ? String(item.recclass) : null,
    mass: typeof item.mass === 'number' ? item.mass : null,
    year: typeof item.year === 'number' ? item.year : null,
    value: typeof item.value === 'number' ? item.value : null,
    unit: item.unit ? String(item.unit) : null,
    metadata: item.metadata ? String(item.metadata) : null,
    lat: typeof item.lat === 'number' ? item.lat : null,
    lon: typeof item.lon === 'number' ? item.lon : null
  }));

  console.log("🔹 Cleaned data from API:", JSON.stringify(cleanResults, null, 2));

  const mappedLocations = useMemo(() => {
    return cleanResults
      .filter((m) => {
        const valid = typeof m.lat === "number" && typeof m.lon === "number" && !isNaN(m.lat) && !isNaN(m.lon);
        if (!valid) console.warn(`⚠️ Skipping invalid entry:`, m);
        return valid;
      })
      .map((m) => ({
        name: m.name,
        lat: m.lat!,
        lon: m.lon!,
      }));
  }, [cleanResults]);

  console.log("🗺️ Map Pins Data:", JSON.stringify(mappedLocations, null, 2));

  const getDatasetTitle = () => {
    switch (datasetType) {
      case "meteorite":
        return "☄️ Meteorites";
      case "climate":
        return "🌡️ Climate Stations";
      case "wind":
        return "💨 Wind Observations";
      case "vegetation":
        return "🌿 Vegetation Zones";
      case "infrastructure":
        return "🏗️ Infrastructure";
      case "fire":
        return "🔥 Fire Data";
      default:
        return "📍 Data Points";
    }
  };

  const getValueDisplay = (item: any) => {
    if (datasetType === "meteorite") {
      return item.mass ? `${item.mass}g` : "Unknown";
    }
    if (datasetType === "climate") {
      return item.value ? `${item.value}°C` : "Unknown";
    }
    if (datasetType === "wind") {
      return item.value ? `${item.value} m/s` : "Unknown";
    }
    if (datasetType === "vegetation") {
      return item.value ? `${(item.value / 1000000).toFixed(1)} km²` : "Unknown";
    }
    return item.value ? `${item.value}${item.unit || ''}` : "Unknown";
  };

  const getValueLabel = () => {
    switch (datasetType) {
      case "meteorite":
        return "Mass";
      case "climate":
        return "Temperature";
      case "wind":
        return "Wind Speed";
      case "vegetation":
        return "Area";
      default:
        return "Value";
    }
  };

  const getMetadataDisplay = (item: any) => {
    if (datasetType === "meteorite") {
      return (
        <>
          <p className="text-gray-300">Class: {item.recclass || "N/A"}</p>
          <p className="text-gray-300">Year: {item.year || "Unknown"}</p>
        </>
      );
    }
    
    if (item.metadata) {
      try {
        const metadata = JSON.parse(item.metadata);
        return (
          <div className="text-gray-300 text-sm">
            {Object.entries(metadata).slice(0, 3).map(([key, value]) => (
              <p key={key}>{key}: {typeof value === 'object' ? JSON.stringify(value) : String(value)}</p>
            ))}
          </div>
        );
      } catch (e) {
        return <p className="text-gray-300">Metadata: {item.metadata}</p>;
      }
    }
    
    return null;
  };

  return (
    <div>
      <h2 className="text-2xl font-bold text-blue-400 mb-4">
        {getDatasetTitle()} - Found {results.length} Records
      </h2>

    
      {mappedLocations.length > 0 && (
        <div className="mt-4">
          <Map meteorites={mappedLocations} />
        </div>
      )}

  
      <div className={styles.gridContainer}>
        {cleanResults.map((item) => (
          <div key={item.id || Math.random()} className={styles.meteoriteCard}>
            <RockImage rockId={item.id} />
            <h3 className="text-xl font-bold text-blue-400 mt-2">
              {item.name || "Unknown Name"}
            </h3>
            <p className="text-gray-300">{getValueLabel()}: {getValueDisplay(item)}</p>
            {getMetadataDisplay(item)}

         
            <p className="text-gray-300">
              📍 Location:{" "}
              {item.lat !== null && item.lon !== null ? (
                <>
                  {item.lat.toFixed(2)}, {item.lon.toFixed(2)}{" "}
                  <Geocode lat={item.lat} lon={item.lon} />
                </>
              ) : (
                "Unknown"
              )}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
