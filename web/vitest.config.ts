import path from 'path'

import react from '@vitejs/plugin-react'
import { defineConfig, configDefaults } from 'vitest/config'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: [
      {
        find: /.*\.svg$/,
        replacement: path.resolve(__dirname, './src/__mocks__/svg.tsx'),
      },
      {
        find: '@',
        replacement: path.resolve(__dirname, './src'),
      },
    ],
  },
  test: {
    environment: 'jsdom',
    globals: true,
    exclude: [...configDefaults.exclude, 'e2e/**', 'playwright-report/**'],
  },
})
