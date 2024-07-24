{
  inputs.nixpkgs.url = "nixpkgs/nixos-24.05";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.treefmt-nix.url = "github:numtide/treefmt-nix";
  inputs.treefmt-nix.inputs.nixpkgs.follows = "nixpkgs";

  outputs = { self, nixpkgs, flake-utils, treefmt-nix }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        appVersion =
          if (self ? shortRev)
          then self.shortRev
          else "development";
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
            src = pkgs.lib.cleanSource self;
            version = appVersion;
          };
        };

        devShells =
          let
            updateVendor = pkgs.writeScriptBin "nix-vendor-sri"
              ''
                #!${pkgs.runtimeShell}
                set -eu

                OUT=$(mktemp -d -t nar-hash-XXXXXX)
                rm -rf "$OUT"

                go mod vendor -o "$OUT"
                go run tailscale.com/cmd/nardump --sri "$OUT" > go.mod.sri
                rm -rf "$OUT"
              '';
          in
          {
            default = pkgs.mkShell {
              buildInputs = with pkgs; [
                go
                gopls
                go-tools
                gotools
                gotests
                updateVendor
              ];
            };
          };
      });
}
