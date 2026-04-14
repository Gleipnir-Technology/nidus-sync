import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import checker from "vite-plugin-checker";
import path from "path";

export default defineConfig({
	plugins: [
		vue(),
		checker({
			vueTsc: true,
		}),
	],

	resolve: {
		alias: {
			"@": path.resolve(__dirname, "../../ts"),
		},
	},

	css: {
		preprocessorOptions: {
			scss: {
				additionalData: `@use "sass:map";\n@import "@/style/variables.scss";`,
				api: "modern-compiler",
				silenceDeprecations: [
					"import",
					"global-builtin",
					"if-function",
					"color-functions",
				],
			},
		},
	},

	build: {
		manifest: false,
		outDir: "static/gen/sync",
		emptyOutDir: true,
		rollupOptions: {
			input: {
				main: path.resolve(__dirname, "./index.html"),
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
		allowedHosts: [
			"poweredge.local",
			"dev-report.mosquitoes.online",
			"dev-sync.nidus.cloud",
		],
		port: 9000,
		proxy: {
			"/api": {
				target: "http://127.0.0.1:9003",
				changeOrigin: false,
			},
			"/configuration/upload/pool/flyover": {
				target: "http://127.0.0.1:9003",
				changeOrigin: false,
			},
			"/mailer": {
				target: "http://127.0.0.1:9003",
				changeOrigin: false,
			},
			"/oauth": {
				target: "http://127.0.0.1:9003",
				changeOrigin: false,
			},
			"/qr-code": {
				target: "http://127.0.0.1:9003",
				changeOrigin: false,
			},
			"/signin": {
				target: "http://localhost:9003",
				changeOrigin: false,
			},
			"/signup": {
				target: "http://localhost:9003",
				changeOrigin: false,
			},
		},
		strictPort: true,
	},
});
