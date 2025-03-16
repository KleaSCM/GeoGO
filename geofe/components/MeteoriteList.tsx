"use client";
import dynamic from "next/dynamic";
import Geocode from "./Geocode";
import RockImage from "./RockImage";
import { useMemo } from "react";
import styles from "./MeteoriteList.module.scss";

// Lazy load the Leaflet map (disable SSR)
const Map = dynamic(() => import("./Map"), { ssr: false });

interface Meteorite {
  id: number;
  name: string;
  recclass: string;
  mass: number;
  year: number;
  lat: number;
  lon: number;
}

interface MeteoriteListProps {
  results: Meteorite[];
}

export default function MeteoriteList({ results }: MeteoriteListProps) {
  if (!results.length) {
    return <p className="text-gray-400 text-center mt-6">No results found.</p>;
  }

  
  const mappedLocations = useMemo(() => {
    return results
      .filter((m) => m.lat !== 0 && m.lon !== 0)
      .map((m) => ({
        name: m.name,
        lat: m.lat,
        lon: m.lon,
      }));
  }, [results]);

  return (
    <div>
      <h2 className="text-2xl font-bold text-blue-400 mb-4">
        🗺️ Found {results.length} Meteorites
      </h2>

     
      <div className="mt-4">
        <Map meteorites={mappedLocations} />
      </div>

     
      <div className={styles.gridContainer}>
        {results.map((rock) => (
          <div key={rock.id} className={styles.meteoriteCard}>
            <RockImage rockId={rock.id} />
            <h3 className="text-xl font-bold text-blue-400 mt-2">
              {rock.name}
            </h3>
            <p className="text-gray-300">Class: {rock.recclass}</p>
            <p className="text-gray-300">Mass: {rock.mass}g</p>
            <p className="text-gray-300">Year: {rock.year}</p>

            
            {rock.lat && rock.lon ? (
              <>
                <p className="text-gray-300">🌍 Coordinates: {rock.lat.toFixed(4)}, {rock.lon.toFixed(4)}</p>
                <p className="text-gray-300">📍 Location: <Geocode lat={rock.lat} lon={rock.lon} /></p>
              </>
            ) : (
              <p className="text-gray-300">📍 Location: Unknown</p>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
