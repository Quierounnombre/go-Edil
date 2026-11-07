{
	inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

	outputs = { self, nixpkgs }:
	let
		system = "x86_64-linux";
		pkgs = nixpkgs.legacyPackages.${system};
	in {
		packages.${system}.default = pkgs.buildGoModule {
			pname = "edil";
			version = "0.1.0";
			src = ./.;
			subPackages = [ "cmd/edil" "cmd/gendocs" ];
			vendorHash = "sha256-hid0eSPm/yfUFBnPnQn2VpHGt/vIEB+zFZEaHuCBJZU="
		};

		devShells.${system}.default = pkgs.mkShell {
			packages = [ pkgs.go ];
		};

		checks.${system}.docs = pkgs.runCommand "docs" {} ''
			cd ${./.}
			${self.packages.${system}.default}/bin/gendocs -check
			touch $out
		'';
	};
}
