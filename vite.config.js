import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [tailwindcss()],
  build: {
    outDir: 'static',
    rollupOptions: {
      input: 'static/css/input.css',
      output: {
        assetFileNames: 'css/output.css',
      },
    },
    emptyOutDir: false,
  },
})
