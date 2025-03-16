"use client";
import dynamic from "next/dynamic";
import Geocode from "./Geocode";
import RockImage from "./RockImage";
import { useMemo, useEffect } from "react";
import styles from "./MeteoriteList.module.scss";

// Lazy load Leaflet map (disable SSR)
const Map = dynamic(() => import("./Map"), { ssr: false });

interface Meteorite {
  id: number;
  name: string;
  recclass: string;
  mass: number;
  year: number;
  lat?: number; 
  lon?: number;
}

interface MeteoriteListProps {
  results: Meteorite[];
}

export default function MeteoriteList({ results }: MeteoriteListProps) {
  if (!Array.isArray(results) || results.length === 0) {
    return <p className="text-gray-400 text-center mt-6">No meteorites found.</p>;
  }

  // Get valid coordinates for the map
  const mappedLocations = useMemo(() => {
    return results
      .filter((m) => typeof m.lat === "number" && typeof m.lon === "number")
      .map((m) => ({
        name: m.name,
        lat: m.lat!,
        lon: m.lon!,
      }));
  }, [results]);

  //debugs bugs bugs bugs bugs 
  useEffect(() => {
    console.log("🗺️ Mapped Meteorites for Map:", mappedLocations);
  }, [mappedLocations]);

  return (
    <div>
      <h2 className="text-2xl font-bold text-blue-400 mb-4">
        🗺️ Found {results.length} Meteorites
      </h2>

      
      {mappedLocations.length > 0 ? (
        <div className="mt-4">
          <Map meteorites={mappedLocations} />
        </div>
      ) : (
        <p className="text-gray-400 text-center mt-6">No valid meteorite locations found.</p>
      )}

      
      <div className={styles.gridContainer}>
        {results.map((rock) => (
          <div key={rock.id || Math.random()} className={styles.meteoriteCard}>
            <RockImage rockId={rock.id} />
            <h3 className="text-xl font-bold text-blue-400 mt-2">
              {rock.name || "Unknown Name"}
            </h3>
            <p className="text-gray-300">Class: {rock.recclass || "N/A"}</p>
            <p className="text-gray-300">Mass: {rock.mass ? `${rock.mass}g` : "Unknown"}</p>
            <p className="text-gray-300">Year: {rock.year || "Unknown"}</p>

            
            <p className="text-gray-300">
              📍 Coordinates: {rock.lat && rock.lon ? `${rock.lat}, ${rock.lon}` : "Unknown"}
            </p>
            <p className="text-gray-300">
              🌍 Location: {rock.lat && rock.lon ? <Geocode lat={rock.lat} lon={rock.lon} /> : "Unknown"}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
