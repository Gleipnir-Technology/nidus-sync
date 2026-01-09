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
        vendorHash = "sha256-fU0FPvuDC3rwQ4ygQYA3sH48o8PaK5VqqowwdUsVNaA=";
}
