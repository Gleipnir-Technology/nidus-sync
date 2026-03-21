// generate-icons.js
import fs from "fs";
import path from "path";

const svgDir = "./svg";
const outputFile = "./ts/gen/custom-icons.scss";
const outputDir = path.dirname(outputFile);

function svgToDataUri(svgContent) {
	// Remove unnecessary attributes first
	svgContent = svgContent
		.replace(/\s+id="[^"]*"/g, "")
		.replace(/\s+id='[^']*'/g, "")
		.replace(/\s+data-name="[^"]*"/g, "")
		.replace(/\s+data-name='[^']*'/g, "");

	// Encode for data URI
	return svgContent
		.replace(/%/g, "%25")
		.replace(/</g, "%3C")
		.replace(/>/g, "%3E")
		.replace(/#/g, "%23")
		.replace(/"/g, "'") // Use single quotes in SVG
		.replace(/'/g, "%27") // Then encode them
		.replace(/\s+/g, " ") // Collapse whitespace
		.trim();
}

function generateIconStyles() {
	const svgFiles = fs.readdirSync(svgDir).filter((f) => f.endsWith(".svg"));

	let scss = `// Auto-generated custom icons\n\n`;
	scss += `.bi-custom {
  display: inline-block;
  width: 1em;
  height: 1em;
  vertical-align: -0.125em;
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
}\n\n`;

	svgFiles.forEach((file) => {
		const iconName = path.basename(file, ".svg");
		const svgContent = fs.readFileSync(path.join(svgDir, file), "utf-8");
		const dataUri = svgToDataUri(svgContent);

		scss += `.bi-${iconName} {
  @extend .bi-custom;
  background-image: url("data:image/svg+xml,${dataUri}");
}\n\n`;
	});

	fs.writeFileSync(outputFile, scss);
	console.log(`Generated ${svgFiles.length} icon styles`);
}

generateIconStyles();
