import type { NextConfig } from "next";

const nextConfig: NextConfig =
  process.env.NODE_ENV === "production"
    ? {
        env: {
          BACKEND_URL: "https://chatifiy-3.fly.dev",
          WS_URL: "wss://chatifiy-3.fly.dev",
          AWS: "https://atheer-web-projects.s3.eu-central-1.amazonaws.com/",
        },
        images: {
          remotePatterns: [
            {
              protocol: "https",
              hostname: "atheer-web-projects.s3.eu-central-1.amazonaws.com",
              port: "",
              pathname: "/**",
            },
          ],
        },
      }
    : {
        env: {
          AWS: "https://atheer-web-projects.s3.eu-central-1.amazonaws.com/",
          BACKEND_URL: "http://localhost:8080",
          WS_URL: "ws://localhost:8080",
          NEXTAUTH_SECRET: "N9rW6PeFVMGUgRIkxMXH7B5YChatify-Deployed",
        },
        images: {
          remotePatterns: [
            {
              protocol: "https",
              hostname: "atheer-web-projects.s3.eu-central-1.amazonaws.com",
              port: "",
              pathname: "/**",
            },
          ],
        },
      };

export default nextConfig;
