import "./globals.css";
import React from "react";
import NavHeader from "./_components/navHeader";
import Footer from "./_components/footer";


export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
    <body>
        <NavHeader/>
        {children}
        <Footer/>
    </body>
    </html>
  );
}
