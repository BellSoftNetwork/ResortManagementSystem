import type { CapacitorConfig } from "@capacitor/cli";

// Get environment variable with fallback
const API_URL = process.env.VITE_API_URL || "";
const NODE_ENV = process.env.NODE_ENV || "development";

// Determine the server URL for mobile apps
// For Android/iOS apps, we need to specify the full API URL
// For local development with emulator, use the special Android emulator IP
// For production builds, use the API_URL environment variable
let serverUrl = "";

if (process.env.CAPACITOR_ANDROID_EMULATOR === "true") {
  // Special IP for Android emulator that points to host machine's localhost
  serverUrl = "http://10.0.2.2:8080";
} else if (API_URL) {
  // Use the provided API URL for mobile builds
  serverUrl = API_URL;
}

const config: CapacitorConfig = {
  appId: "net.bellsoft.rms",
  appName: "Resort Management System",
  webDir: "dist",
  // Server configuration for loading the app
  server: {
    url: serverUrl,
    cleartext: true,
  },
};

console.log(`Capacitor config using server URL: ${config.server.url}`);
console.log(`Environment: ${NODE_ENV}`);

export default config;
