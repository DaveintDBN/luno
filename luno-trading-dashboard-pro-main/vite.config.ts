import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import path from "path";
import { componentTagger } from "lovable-tagger";

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => ({
  server: {
    host: "::",
    port: 3000,
    proxy: {
      '/pairs': 'http://localhost:8080',
      '/config': 'http://localhost:8080',
      '/scan': 'http://localhost:8080',
      '/autoscan': 'http://localhost:8080',
      '/backtest': 'http://localhost:8080',
      '/logs': 'http://localhost:8080',
      '/orderbook': 'http://localhost:8080',
      '/balances': 'http://localhost:8080',
      '/percent-change': 'http://localhost:8080',
      '/healthz': 'http://localhost:8080',
      '/opportunities': 'http://localhost:8080',
      '/status': 'http://localhost:8080',
      '/simulate': 'http://localhost:8080',
      '/execute': 'http://localhost:8080',
    },
  },
  plugins: [
    react(),
    mode === 'development' &&
    componentTagger(),
  ].filter(Boolean),
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
}));
