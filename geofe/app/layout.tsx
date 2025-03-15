import React from "react";
import "@/app/globals.css";
export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="bg-gray-900 text-gray-100 min-h-screen flex flex-col">
        <nav className="bg-gray-800 p-4 shadow-md">
          <div className="container mx-auto flex justify-between items-center">
            <h1 className="text-xl font-bold text-blue-400">GeoGO 🌍</h1>
            <span className="text-gray-300">ROCKS</span>
          </div>
        </nav>
        <main className="flex-1 container mx-auto p-6">{children}</main>
        <footer className="bg-gray-800 text-gray-400 text-center py-4 mt-6">
          © {new Date().getFullYear()} GeoGO - ROCKS VERYROCKS
        </footer>
      </body>
    </html>
  );
}
