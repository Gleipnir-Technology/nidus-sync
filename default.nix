{ pkgs ? import <nixpkgs> { }, proj ? pkgs.proj }:

pkgs.buildGoModule rec {
  meta = {
    description = "Nidus Sync";
    homepage = "https://github.com/Gleipnir-Technology/nidus-sync";
  };
  pname = "nidus-sync";
  src = ./.;
  subPackages = [];
  version = "0.0.12";
  vendorHash = "sha256-o+AFNJX/7/jnaRPVEWSn0R3MuT7lVCiXC7HRHYJOCK8=";

  buildInputs = [ pkgs.proj ];

  nativeBuildInputs = [
    pkgs.pkg-config
    pkgs.nodejs
    pkgs.pnpm.configHook
  ];

  # Fix: Filter out pnpm.configHook instead of replacing the whole list
  overrideModAttrs = old: {
    nativeBuildInputs = builtins.filter 
      (pkg: pkg != pkgs.pnpm.configHook && pkg != pkgs.nodejs) 
      old.nativeBuildInputs;
    preBuild = "";
  };

  pnpmDeps = pkgs.pnpm.fetchDeps {
    inherit pname src version;
    fetcherVersion = 2;
    hash = "sha256-UvE49UmVw8zVFHywxRWyzL0EiZvuZjmm9hA1U98o2sA=";
  };

  preBuild = ''
    pnpm install --offline --frozen-lockfile --ignore-scripts
    mkdir -p "./ts/gen"
    pnpm generate-icons
    pnpm build-rmo
    pnpm build-sync
  '';
}
