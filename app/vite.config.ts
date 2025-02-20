import { Agent } from 'node:http'
import { fileURLToPath, URL } from 'node:url'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import UnoCSS from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
import DefineOptions from 'unplugin-vue-define-options/vite'
import { defineConfig, loadEnv } from 'vite'
import vitePluginBuildId from 'vite-plugin-build-id'
import svgLoader from 'vite-svg-loader'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    base: './',
    server: {
      host: true,
      port: Number.parseInt(env.VITE_PORT) || 3002,
      proxy: {
        '/api': {
          target: env.VITE_PROXY_TARGET || 'http://localhost:9000',
          changeOrigin: true,
          secure: false,
          ws: true,
          timeout: 50000,  // Tambah timeout untuk request lama
          agent: new Agent({
            keepAlive: true,  // Aktifkan keepAlive
            keepAliveMsecs: 1000,
            maxSockets: 10,
            timeout: 60000  // Timeout untuk koneksi agent
          }),
          onProxyReq(proxyReq, req) {
            // Hanya set header untuk non-event-stream
            if (req.headers.accept !== 'text/event-stream') {
              proxyReq.setHeader('Connection', 'keep-alive')
              proxyReq.setHeader('Keep-Alive', 'timeout=5, max=1000')
            } else {
              proxyReq.setHeader('Cache-Control', 'no-cache')
              proxyReq.setHeader('Content-Type', 'text/event-stream')
            }
          },
          onProxyReqWs(proxyReq, req, socket) {
            socket.on('close', () => {
              proxyReq.destroy()
            })
          },
        },
      },
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
      extensions: [
        '.mjs',
        '.js',
        '.ts',
        '.jsx',
        '.tsx',
        '.json',
        '.vue',
        '.less',
      ],
    },
    plugins: [
      vue(),
      vueJsx(),
      vitePluginBuildId(),
      svgLoader(),
      UnoCSS(),
      Components({
        resolvers: [AntDesignVueResolver({ importStyle: false })],
        directoryAsNamespace: true,
      }),
      AutoImport({
        imports: [
          'vue',
          'vue-router',
          'pinia',
          {
            '@/gettext': [
              '$gettext',
              '$pgettext',
              '$ngettext',
              '$npgettext',
            ],
          },
        ],
        vueTemplate: true,
        eslintrc: {
          enabled: true,
          filepath: '.eslint-auto-import.mjs',
        },
      }),
      DefineOptions(),
    ],
    css: {
      preprocessorOptions: {
        less: {
          modifyVars: {
            'border-radius-base': '5px',
          },
          javascriptEnabled: true,
        },
      },
    },
    build: {
      chunkSizeWarningLimit: 1000,
    },
  }
})
