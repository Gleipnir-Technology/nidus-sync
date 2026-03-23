import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import path from "path";

export default defineConfig({
    plugins: [vue()],

    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./ts"),
        },
    },

    css: {
        preprocessorOptions: {
            scss: {
                api: "modern-compiler",
                silenceDeprecations: ["import", "global-builtin", "if-function"],
            },
        },
    },

    build: {
        manifest: true,
        outDir: "static/gen",
        emptyOutDir: true,
        rollupOptions: {
            input: {
                main: path.resolve(__dirname, "ts/main.ts"),
            },
            output: {
                entryFileNames: "js/bundle.[hash].js",
                chunkFileNames: "js/[name].[hash].js",
                assetFileNames: (assetInfo) => {
                    if (/\.(woff2?|ttf|eot)$/.test(assetInfo.name || "")) {
                        return "fonts/[name].[hash][extname]";
                    }
                    if (/\.css$/.test(assetInfo.name || "")) {
                        return "css/style.[hash][extname]";
                    }
                    return "assets/[name].[hash][extname]";
                },
            },
        },
        sourcemap: true,
    },

    server: {
        allowedHosts: ["poweredge.local", "dev-sync.nidus.cloud"],
        port: 9000,
        proxy: {
            "/api": {
                target: "http://127.0.0.1:9002",
                changeOrigin: true,
            },
            "/signup": {
                target: "http://localhost:9002",
                changeOrigin: false,
            }
        },
        strictPort: true,
    },
});
