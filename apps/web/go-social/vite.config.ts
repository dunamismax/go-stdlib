import { defineConfig } from 'vite'

export default defineConfig({
  root: 'frontend',
  build: {
    outDir: '../dist',
    emptyOutDir: true,
    rollupOptions: {
      input: {
        main: 'frontend/main.ts'
      },
      output: {
        entryFileNames: 'assets/[name]-[hash].js',
        chunkFileNames: 'assets/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash].[ext]'
      }
    }
  },
  server: {
    port: 3001,
    proxy: {
      '/api': 'http://localhost:8081',
      '/static': 'http://localhost:8081',
      '/login': 'http://localhost:8081',
      '/register': 'http://localhost:8081',
      '/logout': 'http://localhost:8081',
      '/post': 'http://localhost:8081',
      '/like': 'http://localhost:8081'
    }
  }
})