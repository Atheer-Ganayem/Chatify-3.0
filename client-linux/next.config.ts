import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  env: {
    NEXTAUTH_URL: "http://192.168.1.16:3000",
    BACKEND_URL: "http://192.168.1.16:8080",
    WS_URL: "ws://192.168.1.16:8080",
    AWS: "https://atheer-web-projects.s3.eu-central-1.amazonaws.com/",
  },
  /* config options here */
};

export default nextConfig;
