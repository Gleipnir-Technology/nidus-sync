import esbuild from "esbuild";
import vue from "esbuild-plugin-vue3";

const args = process.argv.slice(2);
const watch = args.includes("--watch");
const minify = args.includes("--minify");

const config = {
	entryPoints: ["ts/main.ts"],
	bundle: true,
	format: "esm",
	plugins: [vue()],
	define: {
		__VUE_OPTIONS_API__: "true",
		__VUE_PROD_DEVTOOLS__: "false",
		__VUE_PROD_HYDRATION_MISMATCH_DETAILS__: "false",
	},
	minify,
	loader: {
		".css": "css",
	},
	outdir: "static/gen",
	outbase: "ts",
	assetNames: "css/[name]",
};

if (watch) {
	const ctx = await esbuild.context(config);
	await ctx.watch();
	console.log("Watching for changes...");
} else {
	await esbuild.build(config);
}
