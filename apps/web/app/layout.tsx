import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Genkit Recipe Studio",
  description: "Type-safe recipe generator powered by Genkit and Gemini"
};

export default function RootLayout({
  children
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
