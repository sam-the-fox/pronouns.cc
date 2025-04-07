import adapter from "@sveltejs/adapter-node";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";
import * as child_process from "node:child_process";
import * as path from "node:path";

import { config as dotenv } from "dotenv";
dotenv({
  path: [path.resolve(process.cwd(), ".env"), path.resolve(process.cwd(), ".env.local")],
});

console.log(process.env.NODE_ENV);
const isProd = process.env.NODE_ENV === "production";

/** @type {import('@sveltejs/kit').Config} */
const config = {
  // Consult https://kit.svelte.dev/docs/integrations#preprocessors
  // for more information about preprocessors
  preprocess: vitePreprocess(),

  kit: {
    adapter: adapter(),
    version: {
      name: child_process.execSync("git describe --tags --long --always").toString().trim(),
    },
    // If PRIVATE_ASSETS_PREFIX is set in .env, serve assets from there.
    paths: {
      assets: isProd ? process.env.PRIVATE_ASSETS_PREFIX || "" : "",
    },
  },
};

export default config;
