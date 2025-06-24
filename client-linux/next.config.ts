import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  env: {
    // NEXTAUTH_URL: "http://192.168.1.16:3000",
    BACKEND_URL: "https://chatifiy-3.fly.dev",
    WS_URL: "wss://chatifiy-3.fly.dev",
    AWS: "https://atheer-web-projects.s3.eu-central-1.amazonaws.com/",
  },
  /* config options here */
};

export default nextConfig;
