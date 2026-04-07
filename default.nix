{ pkgs ? import <nixpkgs> { }, proj ? pkgs.proj }:

let
	# Fetch pnpm dependencies
	pnpmDeps = pkgs.pnpm.fetchDeps {
		pname = "nidus-sync-frontend";
		version = "0.0.11";
		src = ./.;
		hash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="; # nix will tell you the correct hash
	};
in
pkgs.buildGoModule rec {
	meta = {
		description = "Nidus Sync";
		homepage = "https://github.com/Gleipnir-Technology/nidus-sync";
	};
	pname = "nidus-sync";
	src = ./.;
	subPackages = [];
	version = "0.0.11";
	vendorHash = "sha256-6g2gk7AyRmoqYfqwsTbzc5u5IKzlNFHjD3TeYXnf1zA=";

	buildInputs = [ pkgs.proj ];

	# Only inclue pkg-config here - it's needed for both phases
	nativeBuildInputs = [
		pkgs.pkg-config
		pkgs.nodejs
		pkgs.pnpm.configHook
	];

	# Override the go modules derivation to expclude pnpm stuff
	overrideModAttrs = (_: {
		preBuild = ""; # Don't run pnpm stuff during go modules fetch
	});

	pnpmDeps = pkgs.pnpm.fetchDeps {
		inherit pname src version;
		fetcherVersion = 2;
		hash = "sha256-UvE49UmVw8zVFHywxRWyzL0EiZvuZjmm9hA1U98o2sA=";
	};

	preBuild = ''
# Add pnpm and nodejs to PATH for this phase only
pnpm install --offline --frozen-lockfile --ignore-scripts

# Icon generation
mkdir -p "./ts/gen"
pnpm generate-icons

# Remove static links
# Build frontend
pnpm build-rmo
pnpm build-sync
'';
}
