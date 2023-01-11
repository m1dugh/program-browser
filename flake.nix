{
    description = "A golang package to fetch all programs from divers bug bounty platforms";

    inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    };

    outputs = {
        self,
        nixpkgs
    }:
    let
        system = "x86_64-linux";
        pkgs = import nixpkgs {
            inherit system;
            config.allowUnfree = true;
        };
        name = "program-browser";
        inherit (nixpkgs) lib;
    in {
        
        devShells.${system}.default = pkgs.mkShell {
            nativeBuildInputs = with pkgs; [
                gnumake
                go
            ];
        };
    };
}
