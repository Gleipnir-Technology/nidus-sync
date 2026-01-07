{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule rec {
        meta = {
                description = "Nidus Sync";
                homepage = "https://github.com/Gleipnir-Technology/nidus-sync";
        };
        pname = "nidus-sync";
        src = ./.;
        subPackages = [];
        version = "0.0.5";
        # Needs to be updated after every modification of go.mod/go.sum
        vendorHash = "sha256-JuRB/2pNidonwQm6iCyUXiFxkAwizEXdwieRRjgSBwI=";
}
