{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule rec {
        meta = {
                description = "Nidus Sync";
                homepage = "https://github.com/Gleipnir-Technology/nidus-sync";
        };
        pname = "nidus-sync";
        src = ./.;
        subPackages = [];
        version = "0.0.4";
        # Needs to be updated after every modification of go.mod/go.sum
        vendorHash = "sha256-7dEwIQMFGhNIMAlu3tiZ3PQoi5fq3sma85d0mEL98E0=";
}
