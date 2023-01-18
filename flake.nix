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
        name = "program-browser";
        inherit (nixpkgs) lib;
        supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-linux" "aarch64-linux" ];
        forAllSystems = lib.genAttrs supportedSystems;
        nixpkgsFor = forAllSystems(system: import nixpkgs {
            config.allowUnfree = true;
            inherit system;
        });
    in {

        packages = forAllSystems(system: 
        let
            pkgs = nixpkgsFor.${system};
        in {
            program-browser = pkgs.buildGoModule {
                pname = "program-browser";
                src = ./.;
                version = "0.0.1";
                vendorHash = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
            };

            subdomain-finder = pkgs.buildGoModule {
                pname = "subdomain-finder";
                src = ./.;
                version = "0.0.1";
                vendorHash = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
            };
        });

        apps = forAllSystems(system:
        let
            pkgs = nixpkgsFor.${system};
            mypkgs = self.packages.${system};
            inherit (mypkgs) program-browser subdomain-finder;
        in {
            program-browser = {
                type = "app";
                program = "${program-browser}/bin/program-browser";
            };

            subdomain-finder = {
                type = "app";
                program = "${subdomain-finder}/bin/subdomain-finder";
            };
            default = self.apps.${system}.program-browser;
        });

        devShells = forAllSystems(system: 
        let
            pkgs = nixpkgsFor.${system};
        in {
            default = pkgs.mkShell {
                nativeBuildInputs = with pkgs; [
                    gnumake
                    go
                    jq
                ];
            };
        });
    };
}
