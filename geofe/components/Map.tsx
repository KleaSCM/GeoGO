"use client";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import L, { LatLngExpression } from "leaflet";
import "leaflet/dist/leaflet.css";
L.Icon.Default.mergeOptions({
  iconRetinaUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png",
  iconUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png",
  shadowUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png",
});

interface MapProps {
  meteorites: { name: string; lat: number; lon: number }[];
}

export default function Map({ meteorites }: MapProps) {
  const defaultCenter: LatLngExpression = [0, 0];

  return (
    <div className="w-full h-[500px] mb-8">
      <MapContainer
        center={defaultCenter as [number, number]}
        zoom={2}
        className="w-full h-full rounded-lg"
      >
        <TileLayer
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          attribution="&copy; OpenStreetMap contributors"
        />
        {meteorites.map((rock, index) => {
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
