"use client";
import { themes } from "@qttdev/ui/themes";
import { lazy } from "react";
import "./globals.css";
const loader = () =>
  import("@qttdev/ui/themes").then((mod) => ({ default: mod.ThemeProvider }));
const ThemeProvider = lazy(loader);
export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <ThemeProvider theme={themes["light"]}>{children}</ThemeProvider>
      </body>
    </html>
  );
}
