{ src ? ./.
, buildGoModule
}:

buildGoModule  {
  name = "bookarr";
  inherit src;

  vendorHash = "";

  subPackages = [ "cmd/bookarr" ];

}
