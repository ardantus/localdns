import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: true, // Needed for Docker
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://backend:8080', // Docker service name
        changeOrigin: true,
        // rewrite: (path) => path.replace(/^\/api/, ''), // If backend doesn't expect /api prefix, but it does in our routes
      }
    }
  }
})
