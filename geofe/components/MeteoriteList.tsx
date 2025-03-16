"use client";
import Geocode from "./Geocode";
import RockImage from "./RockImage";
import dynamic from "next/dynamic";
import { useMemo } from "react";

// Lazy load the Leaflet Map to avoid SSR issues
const Map = dynamic(() => import("./Map"), { ssr: false });

interface Meteorite {
  id: number;
  name: string;
  recclass: string;
  mass: number;
  year: number;
  location: string;
}

interface MeteoriteListProps {
  results: Meteorite[];
}

export default function MeteoriteList({ results }: MeteoriteListProps) {
  if (!results.length) return <p className="text-gray-400 text-center mt-6">No results found.</p>;

  const meteorites = useMemo(() => {
    return results
      .map((meteorite) => {
        let lat = 0;
        let lon = 0;
        try {
          const coords = meteorite.location.match(/POINT\(([-\d.]+) ([-\d.]+)\)/);
          if (coords) {
            lon = parseFloat(coords[1]);
            lat = parseFloat(coords[2]);
          }
        } catch (error) {
          console.error("❌ Failed to extract lat/lon:", error);
        }
        return lat !== 0 && lon !== 0 ? { name: meteorite.name, lat, lon } : null;
      })
      .filter((m) => m !== null);
  }, [results]);

  return (
    <div>
      <h2 className="text-2xl font-bold text-blue-400">🗺️ Found {results.length} Meteorites</h2>

      {/* MAP COMPONENT */}
      <div className="mt-4">
        <Map meteorites={meteorites as { name: string; lat: number; lon: number }[]} />
      </div>

      {/* ROCK LIST */}
      <div className="mt-8 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        {results.map((rock) => {
          let lat = 0, lon = 0;
          try {
            const coords = rock.location.match(/POINT\(([-\d.]+) ([-\d.]+)\)/);
            if (coords) {
              lon = parseFloat(coords[1]);
              lat = parseFloat(coords[2]);
            }
          } catch (error) {
            console.error("❌ Failed to extract lat/lon:", error);
          }

          return (
            <div
              key={rock.id}
              className="bg-gray-800 p-4 rounded-lg shadow-lg hover:scale-105 transition-transform duration-200"
            >
              <RockImage rockId={rock.id} />
              <h3 className="text-xl font-bold text-blue-400">{rock.name}</h3>
              <p className="text-gray-300">Class: {rock.recclass}</p>
              <p className="text-gray-300">Mass: {rock.mass}g</p>
              <p className="text-gray-300">Year: {rock.year}</p>
              <p className="text-gray-300">
                📍 Location: {lat !== 0 && lon !== 0 ? <Geocode lat={lat} lon={lon} /> : "Unknown"}
              </p>
            </div>
          );
        })}
      </div>
    </div>
  );
}
