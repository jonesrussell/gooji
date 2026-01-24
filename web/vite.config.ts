import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  build: {
    manifest: true,
    outDir: 'dist',
    rollupOptions: {
      input: resolve(__dirname, 'src/main.ts'),
    },
  },
  server: {
    origin: 'http://localhost:5173',
  },
})
