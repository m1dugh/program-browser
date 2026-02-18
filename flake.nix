# vim: ts=2 sw=2
{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    flake-utils = {
      url = "github:numtide/flake-utils";
      inputs.systems.follows = "systems";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
    };
  };

  outputs =
    { self
    , nixpkgs
    , flake-utils
    , treefmt-nix
    , ...
    }: flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs {
        inherit system;
      };
      treefmt = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
    in
    {
      packages = {
        default = self.outputs.packages.${system}.program-browser;

        program-browser = pkgs.callPackage ./nix/default.nix { };
      };

      formatter = treefmt.config.build.wrapper;

      devShells.default = pkgs.mkShell {
        nativeBuildInputs = with pkgs; [
          go
          goreleaser
        ];
      };
    });
}
