import { defineConfig } from 'vitest/config';
import { sveltekit } from '@sveltejs/kit/vite';
import { svelteTesting } from '@testing-library/svelte/vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
	plugins: [sveltekit(), svelteTesting(), tailwindcss()],
	test: {
		globals: true,
		environment: 'jsdom',
		setupFiles: ['./vitest-setup-client.ts'],
		include: ['src/**/*.svelte.{test,spec}.{js,ts}', 'src/**/*.test.{js,ts}'],
		exclude: ['tests/e2e/**', 'src/**/*.e2e.{js,ts}']
	}
});
