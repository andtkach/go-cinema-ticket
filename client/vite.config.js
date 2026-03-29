import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 17070,
    proxy: {
      '/movies': 'http://localhost:17080',
      '/sessions': 'http://localhost:17080',
    },
  },
})
