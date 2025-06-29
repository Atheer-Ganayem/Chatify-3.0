import type { NextConfig } from "next";

const nextConfig: NextConfig =
  process.env.NODE_ENV === "production"
    ? {
        env: {
          BACKEND_URL: "https://chatifiy-3.fly.dev",
          WS_URL: "wss://chatifiy-3.fly.dev",
          AWS: "https://atheer-web-projects.s3.eu-central-1.amazonaws.com/",
        },
      }
    : {
        env: {
          AWS: "https://atheer-web-projects.s3.eu-central-1.amazonaws.com/",
          BACKEND_URL: "http://localhost:8080",
          WS_URL: "ws://localhost:8080",
        },
      };

export default nextConfig;
