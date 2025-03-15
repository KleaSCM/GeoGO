"use client";

interface Meteorite {
  id: number;
  name: string;
  recclass: string;
  mass: number;
  year: number;
  location: string;
}

export default function MeteoriteList({ results }: { results: Meteorite[] }) {
  return (
    <div className="mt-6">
      <h3 className="text-lg font-semibold text-gray-300">Results:</h3>
      {results.length === 0 ? (
        <p className="text-gray-400">No meteorites found.</p>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {results.map((m) => (
            <div
              key={m.id}
              className="bg-gray-800 p-4 rounded shadow-md flex flex-col"
            >
              <h4 className="text-blue-300 font-bold">{m.name}</h4>
              <p className="text-gray-400">Class: {m.recclass}</p>
              <p className="text-gray-400">Mass: {m.mass}g</p>
              <p className="text-gray-400">Year: {m.year}</p>
              <p className="text-gray-400">Location: {m.location}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
