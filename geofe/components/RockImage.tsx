"use client";
import React from "react";

interface RockImageProps {
  rockId: number;
}

const rockImages = [
  "https://source.unsplash.com/200x200/?rock&sig=1",
  "https://source.unsplash.com/200x200/?rock&sig=2",
  "https://source.unsplash.com/200x200/?rock&sig=3"
];

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
