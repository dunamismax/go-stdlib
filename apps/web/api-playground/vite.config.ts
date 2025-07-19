import { defineConfig } from 'vite'

export default defineConfig({
  root: 'frontend',
  build: {
    outDir: '../dist',
    emptyOutDir: true,
    rollupOptions: {
      input: 'frontend/main.ts'
    }
  },
  server: {
    port: 3000,
    proxy: {
      // Proxy all API calls to Echo backend
      '^/(?!static|frontend).*': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        ws: true
      }
    }
  }
})