import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  define: {
    __INFER_ENABLED__: process.env.VITE_INFER_ENABLED === 'true',
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8999',
    },
  },
});
