import { defineConfig, loadEnv } from 'vite';
import { resolve } from 'path';
import react from '@vitejs/plugin-react-swc';
import { visualizer } from 'rollup-plugin-visualizer';

// report plugin config
const configVisualizerConfig = (mode) => {
  return mode === 'report'
    ? [
        visualizer({
          emitFile: false,
          sourcemap: true, // When is true,Always add plugin as last option
          filename: 'pack_analyze.html',
          open: true, // auto open
        }),
      ]
    : [];
};

// https://vitejs.dev/config/
/** @ts-ignore */
export default defineConfig(({ mode }) => {
  const envConfig = loadEnv(mode, './env/');
  return {
    envDir: './env',
    build: {
      chunkSizeWarningLimit: 1000,
      reportCompressedSize: false,
      minify: false, // 「terserOptions」need to set terser
      terserOptions: {
        compress: {
          // remove console in production
          drop_console: false,
          drop_debugger: false,
        },
      },
    },
    plugins: [react(), ...configVisualizerConfig(mode)],
    server: {
      host: envConfig.VITE_HOST,
      port: envConfig.VITE_PORT,
      open: true,
    },
    resolve: {
      alias: {
        '@': resolve(__dirname, './src'),
      },
    },
  };
});
