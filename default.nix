{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule rec {
        meta = {
                description = "Nidus Sync";
                homepage = "https://github.com/Gleipnir-Technology/nidus-sync";
        };
        pname = "nidus-sync";
        src = ./.;
        subPackages = [];
        version = "0.0.10";
        # Needs to be updated after every modification of go.mod/go.sum
        vendorHash = "sha256-5E5gQJh2cr/XwDg+XRQEdXW7mkObZMoyqQnfToVuZ10=";
}
