// Next.js requires this file. It wraps every page in <html> and <body>.
// All the app's logic lives in page.tsx.
import "./globals.css";

export const metadata = {
  title: "Servery Reviews",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
