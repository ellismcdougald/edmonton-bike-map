import { defineConfig } from '@playwright/test';

export default defineConfig({
	workers: 1,
	webServer: {
		command: 'npm run build && npm run preview',
		port: 4173,
		url: 'http://localhost:4173', // Explicitly define the URL for the web server
		timeout: 120 * 1000, // Increase webServer startup timeout to 120 seconds (2 minutes)
	},
	use: {
		baseURL: 'http://localhost:4173', // Base URL for all page.goto() calls
		trace: 'on-first-retry', // Record trace only when retrying a test for the first time.
	},
});
