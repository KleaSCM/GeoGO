"use client";
import { useEffect, useState } from "react";

interface GeocodeProps {
  lat: number;
  lon: number;
}

export default function Geocode({ lat, lon }: GeocodeProps) {
  const [location, setLocation] = useState<string | null>(null);

  useEffect(() => {
    async function fetchLocation() {
      if (!lat || !lon) {
        setLocation("Invalid Coordinates");
        return;
      }

      const cacheKey = `${lat},${lon}`;
      const cachedLocation = sessionStorage.getItem(cacheKey);

      if (cachedLocation) {
        setLocation(cachedLocation);
        return;
      }

      try {
        const response = await fetch(
          `http://localhost:8080/meteorites/location?lat=${lat}&lon=${lon}`
        );
        if (!response.ok) throw new Error("Failed to fetch location");

        const data = await response.json();
        if (!data.location) throw new Error("Empty response from server");

        setLocation(data.location);
        sessionStorage.setItem(cacheKey, data.location);
      } catch (error) {
        console.error("❌ Geocoding failed:", error);
        setLocation("Unknown Location");
      }
    }

    fetchLocation();
  }, [lat, lon]);

  return <span>{location || "Fetching..."}</span>;
}
