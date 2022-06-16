import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  root: "frontend",
  plugins: [react()],
  server: {
    proxy: {
      "/api": {
        // assumes port 8080 in .env for development
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
