import { defineConfig } from "@playwright/test";

export default defineConfig({
  workers: 1,
  testDir: "./tests",
  use: {
    baseURL: "http://localhost:3000",
    headless: true,
    viewport: { width: 1280, height: 720 },
  },
});
