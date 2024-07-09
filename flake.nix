{
  inputs.nixpkgs.url = "nixpkgs/nixos-24.05";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.treefmt-nix.url = "github:numtide/treefmt-nix";
  inputs.treefmt-nix.inputs.nixpkgs.follows = "nixpkgs";

  outputs = { self, nixpkgs, flake-utils, treefmt-nix }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        formatter = treefmt-nix.lib.mkWrapper pkgs
          {
            projectRootFile = "flake.nix";
            programs.nixpkgs-fmt.enable = true;
            programs.gofmt.enable = true;
          };

        defaultPackage = self.packages.${system}.default;

        packages = {
          default = self.packages.${system}.bookarr;

          bookarr = pkgs.callPackage self {
            src = self;
          };
        };

        devShells = {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls go-tools gotools jq yq ];
          };
        };
      });
}
