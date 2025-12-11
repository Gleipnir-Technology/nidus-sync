{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule rec {
        meta = {
                description = "Nidus Sync";
                homepage = "https://github.com/Gleipnir-Technology/nidus-sync";
        };
        pname = "nidus-sync";
        src = ./.;
        subPackages = [];
        version = "0.0.3";
        # Needs to be updated after every modification of go.mod/go.sum
        vendorHash = "sha256-Bxv6OPxW9VERETEi2QRk/t3b0HpsoHZFPUsMP85G7Ug=";
}
