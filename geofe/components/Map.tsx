"use client";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
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

export default function Map({ meteorites }: MapProps) {
  // Filter out any invalid coords
  const validRocks = meteorites.filter(
    (rock) => typeof rock.lat === "number" && typeof rock.lon === "number"
  );

  // Center on the first valid rock or [0,0] if none
  const center: LatLngExpression =
    validRocks.length > 0
      ? ([validRocks[0].lat, validRocks[0].lon] as [number, number])
      : ([0, 0] as [number, number]);

  return (
    <div className="w-full h-[500px] mb-8">
      <MapContainer
        center={center}
        zoom={2}
        className="w-full h-full rounded-lg"
      >
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


