"use client";
import Geocode from "./Geocode";

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

  return (
    <div className="mt-8 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {results.map((meteorite) => {
        
        let lat = 0;
        let lon = 0;

        try {
          const coords = meteorite.location.match(/POINT\(([-\d.]+) ([-\d.]+)\)/);
          if (coords) {
            lon = parseFloat(coords[1]);
            lat = parseFloat(coords[2]); 
          } else {
            console.error("❌ Location format incorrect:", meteorite.location);
          }
        } catch (error) {
          console.error("❌ Failed to extract lat/lon:", error);
        }

        return (
          <div
            key={meteorite.id}
            className="bg-gray-800 p-4 rounded-lg shadow-lg hover:scale-105 transition-transform duration-200"
          >
            <h3 className="text-xl font-bold text-blue-400">{meteorite.name}</h3>
            <p className="text-gray-300">Class: {meteorite.recclass}</p>
            <p className="text-gray-300">Mass: {meteorite.mass}g</p>
            <p className="text-gray-300">Year: {meteorite.year}</p>
            <p className="text-gray-300">
              📍 Location: {lat !== 0 && lon !== 0 ? <Geocode lat={lat} lon={lon} /> : "Unknown"}
            </p>
          </div>
        );
      })}
    </div>
  );
}

