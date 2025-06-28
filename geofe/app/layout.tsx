import React from "react";
import "@/app/globals.scss"; 

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="body">
        <main>{children}</main>
      </body>
    </html>
  );
}
