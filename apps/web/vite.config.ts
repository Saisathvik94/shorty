import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');
  const apiUrl = env.VITE_API_URL || 'http://localhost:3000';

  return {
    plugins: [react()],
    server: {
      port: 5173,
      proxy: {
        '/api': {
          target: apiUrl,
          changeOrigin: true,
        },
        // Match short code routes like /nQYAYb (alphanumeric, -, _)
        '^/[a-zA-Z0-9_-]{4,32}$': {
          target: apiUrl,
          changeOrigin: true,
          bypass: (req) => {
            const url = req.url || '';
            // Bypass proxy for static assets, scripts, styles, or Vite internals
            if (
              url.includes('.') ||
              url.startsWith('/@') ||
              url.startsWith('/node_modules') ||
              url.startsWith('/src')
            ) {
              return req.url;
            }
          },
        },
      },
    },
  };
});
