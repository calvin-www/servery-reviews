import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Forward every request that starts with /api/ to the Go backend on port 8080.
  // For example, fetch("/api/reviews") in page.tsx is sent to
  // http://127.0.0.1:8080/api/reviews
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://127.0.0.1:8080/api/:path*",
      },
    ];
  },
};

export default nextConfig;
