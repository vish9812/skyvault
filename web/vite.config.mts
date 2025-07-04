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
  },
  build: {
    target: "esnext",
    outDir: "dist",
  },
});
