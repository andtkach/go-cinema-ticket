import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 17070,
    proxy: {
      '/cinemas': 'http://localhost:17080',
      '/movies': 'http://localhost:17080',
      '/sessions': 'http://localhost:17080',
    },
  },
})
