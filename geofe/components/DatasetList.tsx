"use client";
import dynamic from "next/dynamic";
import Geocode from "./Geocode";
import { useMemo } from "react";
import styles from "./MeteoriteList.module.scss";

// Lazy load Leaflet map (disable SSR)
const Map = dynamic(() => import("./Map"), { ssr: false });

interface DatasetItem {
  id: number;
  name: string;
  recclass?: string | null;
  mass?: number | null;
  year?: number | null;
  lat?: number | null; 
  lon?: number | null;
  value?: number | null;
  unit?: string | null;
  metadata?: string | null;
  // Climate specific fields
  djf?: number | null;
  jja?: number | null;
  mam?: number | null;
  son?: number | null;
  // Wind specific fields
  wind_speed?: number | null;
  wind_direction?: number | null;
  // Vegetation specific fields
  area_km2?: number | null;
  zone_type?: string | null;
}

interface DatasetListProps {
  results: DatasetItem[];
  datasetType?: string;
}

export default function DatasetList({ results, datasetType = "meteorite" }: DatasetListProps) {
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
    lon: typeof item.lon === 'number' ? item.lon : null,
    djf: typeof item.djf === 'number' ? item.djf : null,
    jja: typeof item.jja === 'number' ? item.jja : null,
    mam: typeof item.mam === 'number' ? item.mam : null,
    son: typeof item.son === 'number' ? item.son : null,
    wind_speed: typeof item.wind_speed === 'number' ? item.wind_speed : null,
    wind_direction: typeof item.wind_direction === 'number' ? item.wind_direction : null,
    area_km2: typeof item.area_km2 === 'number' ? item.area_km2 : null,
    zone_type: item.zone_type ? String(item.zone_type) : null
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

  const getValueDisplay = (item: DatasetItem) => {
    if (datasetType === "meteorite") {
      return item.mass ? `${item.mass}g` : "Unknown";
    }
    if (datasetType === "climate") {
      if (item.metadata) {
        try {
          const metadata = JSON.parse(item.metadata);
          const unit = metadata.unit || item.unit || "°C";
          return item.value ? `${item.value}${unit}` : "Unknown";
        } catch (e) {
          return item.value ? `${item.value}${item.unit || '°C'}` : "Unknown";
        }
      }
      return item.value ? `${item.value}${item.unit || '°C'}` : "Unknown";
    }
    if (datasetType === "wind") {
      return item.value ? `${item.value} m/s` : "Unknown";
    }
    if (datasetType === "vegetation") {
      return item.value ? `${(item.value / 1000000).toFixed(1)} km²` : "Unknown";
    }
    return item.value ? `${item.value}${item.unit || ''}` : "Unknown";
  };

  const getValueLabel = (item?: DatasetItem) => {
    if (datasetType === "meteorite") {
      return "Mass";
    }
    if (datasetType === "climate" && item?.metadata) {
      try {
        const metadata = JSON.parse(item.metadata);
        const climateType = metadata.climate_type;
        switch (climateType) {
          case "max_temperature":
            return "Max Temperature";
          case "min_temperature":
            return "Min Temperature";
          case "avg_temperature":
            return "Avg Temperature";
          case "humidity":
            return "Humidity";
          case "evaporation":
            return "Evaporation";
          default:
            return "Climate Value";
        }
      } catch (e) {
        return "Climate Value";
      }
    }
    if (datasetType === "climate") {
      return "Climate Value";
    }
    if (datasetType === "wind") {
      return "Wind Speed";
    }
    if (datasetType === "vegetation") {
      return "Area";
    }
    return "Value";
  };

  const getMetadataDisplay = (item: DatasetItem) => {
    if (datasetType === "meteorite") {
      return (
        <>
          <p className="text-gray-300">Class: {item.recclass || "N/A"}</p>
          <p className="text-gray-300">Year: {item.year || "Unknown"}</p>
        </>
      );
    }
    
    if (datasetType === "climate") {
      if (item.metadata) {
        try {
          const metadata = JSON.parse(item.metadata);
          return (
            <div className="text-gray-300 text-sm">
              {metadata.djf && <p>DJF: {metadata.djf.toFixed(2)}</p>}
              {metadata.jja && <p>JJA: {metadata.jja.toFixed(2)}</p>}
              {metadata.mam && <p>MAM: {metadata.mam.toFixed(2)}</p>}
              {metadata.son && <p>SON: {metadata.son.toFixed(2)}</p>}
              {metadata.station_id && <p>Station ID: {metadata.station_id}</p>}
            </div>
          );
        } catch (e) {
          console.error("Failed to parse climate metadata:", e);
        }
      }
      return null;
    }
    
    if (datasetType === "wind") {
      if (item.metadata) {
        try {
          const metadata = JSON.parse(item.metadata);
          return (
            <div className="text-gray-300 text-sm">
              {metadata.wind_direction && <p>Direction: {metadata.wind_direction}°</p>}
              {metadata.gust_speed && <p>Gust: {metadata.gust_speed} m/s</p>}
            </div>
          );
        } catch (e) {
          console.error("Failed to parse wind metadata:", e);
        }
      }
      return null;
    }
    
    if (datasetType === "vegetation") {
      if (item.metadata) {
        try {
          const metadata = JSON.parse(item.metadata);
          return (
            <div className="text-gray-300 text-sm">
              {metadata.zone && <p>Zone: {metadata.zone}</p>}
              {metadata.type && <p>Type: {metadata.type}</p>}
            </div>
          );
        } catch (e) {
          console.error("Failed to parse vegetation metadata:", e);
        }
      }
      return null;
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

  const getDatasetIcon = (item: DatasetItem) => {
    switch (datasetType) {
      case "meteorite":
        return "☄️";
      case "climate":
        return "🌡️";
      case "wind":
        return "💨";
      case "vegetation":
        return "🌿";
      case "infrastructure":
        return "🏗️";
      case "fire":
        return "🔥";
      default:
        return "📍";
    }
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
            <div className="text-6xl mb-4">{getDatasetIcon(item)}</div>
            <h3 className="text-xl font-bold text-blue-400 mt-2">
              {item.name || "Unknown Name"}
            </h3>
            <p className="text-gray-300">{getValueLabel(item)}: {getValueDisplay(item)}</p>
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