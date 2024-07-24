{ pkgs
, version ? "development"
, src ? ./.
, buildGoModule
}:

buildGoModule {
  name = "bookarr";
  inherit src;
  inherit version;

  vendorHash = pkgs.lib.fileContents ./go.mod.sri;

  subPackages = [ "cmd/bookarr" ];

}
