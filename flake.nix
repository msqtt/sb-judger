{
  description = "A Nix-flake-based Go 1.20 development environment";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }:
    let
		
      goVersion = 20; # Change this to update the whole stack
      overlays = [ (final: prev: { go = prev."go_1_${toString goVersion}"; }) ];
      supportedSystems = [ "x86_64-linux"];
      forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs { inherit overlays system; };
      });
    in
    {
      devShells = forEachSupportedSystem ({ pkgs }: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            # go 1.18 (specified by overlay)
            go
            # goimports, godoc, etc.
            gopls
            gotools

						protobuf
          ];
          shellHook = ''
            export GOPATH="$HOME/env/develop/go/gopath"
						export PATH="$PATH:$GOPATH/bin"
            go version
          '';
        };
      });
    };
}
