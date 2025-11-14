import { defineConfig } from '@playwright/test';

export default defineConfig({
	workers: 1,
	webServer: {
		command: 'npm run build && npm run preview',
		port: 4173
	},
});
