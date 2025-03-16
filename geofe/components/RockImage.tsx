"use client";
import React from "react";

const rockImages = [
  "https://source.unsplash.com/200x200/?rock&sig=1",
  "https://source.unsplash.com/200x200/?rock&sig=2",
  "https://source.unsplash.com/200x200/?rock&sig=3",
];

interface RockImageProps {
  rockId: number;
}

export default function RockImage({ rockId }: RockImageProps) {
  const imageUrl = rockImages[rockId % rockImages.length];
  return (
    <img
      src={imageUrl}
      alt="Rock"
      className="w-full h-32 object-cover rounded-md"
    />
  );
}
