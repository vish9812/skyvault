import { defineConfig } from "vite";
import solidPlugin from "vite-plugin-solid";
import devtools from "solid-devtools/vite";
import tailwindcss from "@tailwindcss/vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
  plugins: [
    devtools({ autoname: true }),
    tailwindcss(),
    solidPlugin(),
    tsconfigPaths(),
  ],
  server: {
    port: 3000,
    proxy: {
      // Forward all API requests to the Go server
      "/api": {
        target: "http://localhost:8090",
        changeOrigin: true,
      },
    },
  },
  build: {
    target: "esnext",
    outDir: "dist",
  },
});
