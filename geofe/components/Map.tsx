"use client";
import { MapContainer, TileLayer, Marker, Popup, useMap } from "react-leaflet";
import { useEffect } from "react";
import L, { LatLngExpression } from "leaflet";
import "leaflet/dist/leaflet.css";

// Override default icon options
L.Icon.Default.mergeOptions({
  iconRetinaUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png",
  iconUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png",
  shadowUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png",
});

interface MapData {
  name: string;
  lat: number;
  lon: number;
}

interface MapProps {
  meteorites: MapData[];
}

// Component to recenter the map dynamically
function ChangeMapView({ center }: { center: LatLngExpression }) {
  const map = useMap();
  useEffect(() => {
    map.setView(center, map.getZoom());
  }, [center, map]);
  return null;
}

export default function Map({ meteorites }: MapProps) {
  // Filter out any invalid coords
  const validRocks = meteorites.filter(
    (rock) => typeof rock.lat === "number" && typeof rock.lon === "number"
  );

  // Default center
  const defaultCenter: LatLngExpression = [20, 0];

  // Center on the first valid rock or use default
  const center: LatLngExpression =
    validRocks.length > 0
      ? ([validRocks[0].lat, validRocks[0].lon] as [number, number])
      : defaultCenter;

  //  Dbug debug bug bug bug bug pins
  useEffect(() => {
    console.log("📍 Markers being displayed:", validRocks);
  }, [validRocks]);

  return (
    <div className="w-full h-[500px] mb-8">
      <MapContainer
        center={center}
        zoom={2}
        className="w-full h-full rounded-lg"
      >
      
        <ChangeMapView center={center} />

        <TileLayer
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          attribution="&copy; OpenStreetMap contributors"
        />

        {validRocks.map((rock, index) => {
          const position: [number, number] = [rock.lat, rock.lon];
          return (
            <Marker key={index} position={position}>
              <Popup>
                <strong>{rock.name}</strong>
                <br />
                {rock.lat.toFixed(2)}, {rock.lon.toFixed(2)}
              </Popup>
            </Marker>
          );
        })}
      </MapContainer>
    </div>
  );
}
