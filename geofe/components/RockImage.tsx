"use client";
import React from "react";

interface RockImageProps {
  rockId: number;
}


const fallbackImages = [
  "/rock.jpg",
"/rock1.png",
"/rock2.jpg",
"/rock3.webp",
"/rock4.jpg",
];

export default function RockImage({ rockId }: RockImageProps) {
    const imageUrl = fallbackImages[rockId % fallbackImages.length];
    return (
      <img
        src={imageUrl}
        alt="Meteorite"
        className="w-full h-32 object-cover rounded-md"
      />
    );
  }