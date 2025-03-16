"use client";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import L, { LatLngExpression } from "leaflet";
import "leaflet/dist/leaflet.css";

// Override default icon options for Leaflet
L.Icon.Default.mergeOptions({
  iconRetinaUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png",
  iconUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png",
  shadowUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png",
});

interface MapProps {
  meteorites: { name: string; lat: number; lon: number }[];
}

export default function Map({ meteorites }: MapProps) {
  const defaultCenter: LatLngExpression = [0, 0];

  return (
    <MapContainer 
      center={defaultCenter as [number, number]} 
      zoom={2} 
      className="h-[500px] w-full rounded-lg"
    >
      <TileLayer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      {meteorites.map((rock, index) => {
        const position: [number, number] = [rock.lat, rock.lon];
        return (
          <Marker key={index} position={position}>
            <Popup>
              <strong>{rock.name}</strong> <br />
              Found at: {rock.lat.toFixed(2)}, {rock.lon.toFixed(2)}
            </Popup>
          </Marker>
        );
      })}
    </MapContainer>
  );
}
