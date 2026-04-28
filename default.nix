{ pkgs ? import <nixpkgs> { }, proj ? pkgs.proj }:

pkgs.buildGoModule rec {
  # Try to get git info, fallback to version if .git doesn't exist
  # Note: This runs at eval time, so it captures the version when you build
  gitRevision = 
    if builtins.pathExists ./.git 
    then pkgs.lib.commitIdFromGitRepo ./.git 
    else "unknown";
  gitDescribe = builtins.readFile (pkgs.runCommand "git-describe" {} ''
    ${pkgs.git}/bin/git -C ${./.} describe --always --dirty --tags 2>/dev/null > $out || echo "${version}" > $out
  '');
  ldflags = [
    "-s"
    "-w"
    "-X main.Version=${version}"
    "-X main.Commit=${gitRevision}"
  ];
  meta = {
    description = "Nidus Sync";
    homepage = "https://github.com/Gleipnir-Technology/nidus-sync";
  };
  pname = "nidus-sync";
  src = ./.;
  subPackages = [];
  version = "0.0.12";
  vendorHash = "sha256-d2cTpzEtUmxAKPKmM0zeKhW4nbLnI08X8lm2yayl+NA=";

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
    hash = "sha256-4XPkwVKSuDlErKfD59iBPHLuf44iRkjEcS6tIityCjo=";
  };

  preBuild = ''
    pnpm install --offline --frozen-lockfile --ignore-scripts
    mkdir -p "./ts/gen"
    pnpm generate-icons
    pnpm build-rmo
    pnpm build-sync
  '';
}
