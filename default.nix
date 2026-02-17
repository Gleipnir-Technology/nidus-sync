{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule rec {
        meta = {
                description = "Nidus Sync";
                homepage = "https://github.com/Gleipnir-Technology/nidus-sync";
        };
        pname = "nidus-sync";
        src = ./.;
        subPackages = [];
        version = "0.0.11";
        # Needs to be updated after every modification of go.mod/go.sum
        vendorHash = "sha256-ax4suxqagDNBIRk/j2M3bMw3D/GXMmIKsALWrsf9zAA=";

	nativeBuildInputs = [ pkgs.dart-sass ];

	preBuild = ''

		SASS_SRC_DIR="./scss"
		CSS_OUTPUT_DIR="./html/static/css/"

		mkdir -p "$CSS_OUTPUT_DIR"

		echo "Compiling $SASS_SRC_DIR/custom.scss to $CSS_OUTPUT_DIR/bootstrap.css..."
		sass --style=compressed --trace "$SASS_SRC_DIR/custom.scss":"$CSS_OUTPUT_DIR/bootstrap.css"
	'';
}
