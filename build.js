import esbuild from "esbuild";
import vue from "esbuild-plugin-vue3";
import { sassPlugin } from "esbuild-sass-plugin";

const args = process.argv.slice(2);
const watch = args.includes("--watch");
const minify = args.includes("--minify");

// Plugin to show build status
const buildStatusPlugin = {
	name: "build-status",
	setup(build) {
		let buildStart;

		build.onStart(() => {
			buildStart = Date.now();
			// Clear console and move cursor to top
			console.clear();
			console.log(
				"\x1b[36m%s\x1b[0m",
				`🔨 Building... [${new Date().toLocaleTimeString()}]`,
			);
		});

		build.onEnd((result) => {
			const buildTime = Date.now() - buildStart;
			if (result.errors.length > 0) {
				console.log(
					"\x1b[31m%s\x1b[0m",
					`❌ Build failed (${buildTime}ms) [${new Date().toLocaleTimeString()}]`,
				);
			} else {
				console.log(
					"\x1b[32m%s\x1b[0m",
					`✅ Build complete (${buildTime}ms) [${new Date().toLocaleTimeString()}]`,
				);
			}
			console.log("\x1b[33m%s\x1b[0m", "\n👀 Watching for changes...");
		});
	},
};

const config = {
	entryPoints: ["ts/main.ts"],
	bundle: true,
	format: "esm",
	plugins: [
		buildStatusPlugin, // Add this first
		sassPlugin({
			quietDeps: true,
			silenceDeprecations: ["import"],
			type: "css",
		}),
		vue(),
	],
	sourcemap: true,
	define: {
		__VUE_OPTIONS_API__: "true",
		__VUE_PROD_DEVTOOLS__: "false",
		__VUE_PROD_HYDRATION_MISMATCH_DETAILS__: "false",
	},
	minify,
	loader: {
		".css": "css",
		".woff": "file",
		".woff2": "file",
		".ttf": "file",
		".eot": "file",
	},
	outdir: "static/gen",
	outbase: "ts",
	assetNames: "fonts/[name]",
};

if (watch) {
	const ctx = await esbuild.context(config);
	await ctx.watch();
	console.log("\x1b[33m%s\x1b[0m", "👀 Watching for changes...\n");
} else {
	await esbuild.build(config);
}
