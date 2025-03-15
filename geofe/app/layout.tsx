import React from "react";
import "@/app/globals.scss"; 


export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
     <body className="body">
        
        <nav className="flex justify-between items-center p-4 shadow-md rounded-lg">
          <h1 className="text-3xl font-bold text-blue-400 transition-transform duration-200 hover:scale-105">
            🌍 GeoGO
          </h1>
          <span className="text-gray-300 text-sm tracking-wide italic">
            Discovering Space Rocks
          </span>
        </nav>

        
        <main className="flex-1 container mx-auto p-6">{children}</main>

      
        <footer>
          <p>Sylavanas is still Warchief</p>
          <p>© {new Date().getFullYear()} FLYING ROCKS</p>
        </footer>
      </body>
    </html>
  );
}
