import { defineConfig } from 'vite';
import { resolve } from 'path';

export default defineConfig({
  root: 'web',
  build: {
    outDir: '../web/static/js',
    emptyOutDir: false,
    rollupOptions: {
      input: {
        app: resolve(__dirname, 'web/src/app.ts'),
        record: resolve(__dirname, 'web/src/record.ts'),
        gallery: resolve(__dirname, 'web/src/gallery.ts'),
        upload: resolve(__dirname, 'web/src/upload.ts'),
        editor: resolve(__dirname, 'web/src/editor.ts'),
      },
      output: {
        entryFileNames: '[name].js',
        format: 'iife',
        sourcemap: true,
      },
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  test: {
    globals: true,
    environment: 'happy-dom',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'web/src/__tests__/',
        '*.config.ts',
        '*.config.js',
      ],
    },
  },
});
