import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";
import { plugin as markdown, Mode as MarkdownMode } from "vite-plugin-markdown";

export default defineConfig({
  plugins: [sveltekit(), markdown({ mode: [MarkdownMode.HTML] })],
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
      "/media": {
        target: "http://localhost:9000",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/media/, "/pronouns.cc"),
      },
    },
  },
  ssr: {
    noExternal: ["@popperjs/core"],
  },
});
