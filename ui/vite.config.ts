import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [
    svelte({
      compilerOptions: {
        generate: 'client',
      },
    }),
  ],
  resolve: {
    conditions: ['browser'],
  },
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:5252',
    },
  },
});
