{
	description = "Nidus sync";

	inputs = {
		nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
		flake-utils.url = "github:numtide/flake-utils";
		proj.url = "github:Gleipnir-Technology/proj";
	};

	outputs = { self, nixpkgs, flake-utils, proj }:
		flake-utils.lib.eachDefaultSystem (system:
			let
				pkgs = nixpkgs.legacyPackages.${system};
				projPkg = proj.packages.${system}.default;

				# Override pkgs.proj with your custom proj
				customPkgs = pkgs // {
					proj = proj.packages.${system}.default;
				};
				package = pkgs.callPackage ./default.nix {
				proj = projPkg;
			};
			in
			{
				packages.default = package;
				packages.nidus-sync = package;

				# Development shell configuration
				devShells.default = pkgs.mkShell {
					buildInputs = [
						pkgs.air
						pkgs.autoprefixer
						pkgs.dart-sass
						pkgs.go
						pkgs.go-jet
						pkgs.golangci-lint
						pkgs.goose
						pkgs.gotools
						pkgs.lefthook
						pkgs.nodejs
						pkgs.pkg-config
						pkgs.pnpm
						pkgs.prettier
						pkgs.prettier-plugin-go-template
						proj.packages.${system}.default
						pkgs.typescript
						pkgs.watchexec
					];
				};
			}
		);
}
