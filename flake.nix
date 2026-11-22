{
	inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

	outputs = { self, nixpkgs }:
	let
		system = "x86_64-linux";
		pkgs = nixpkgs.legacyPackages.${system};
		edil = self.packages.${system}.default;
		# Runs a script on a writable copy of the repo with the vendored modules.
		goCheck = name: script: pkgs.runCommand name { nativeBuildInputs = [ pkgs.go ]; } ''
			cp -r ${./.} src
			chmod -R u+w src
			cd src
			cp -r ${edil.goModules} vendor
			export HOME=$TMPDIR GOFLAGS=-mod=vendor GOTOOLCHAIN=local
			${script}
			touch $out
		'';
	in {
		packages.${system}.default = pkgs.buildGoModule {
			pname = "edil";
			version = "0.4.0";
			src = ./.;
			subPackages = [ "cmd/edil" "cmd/gendocs" ];
			vendorHash = "sha256-nFI8aITdM5ojcLiFXuB8/aQ8M9NgxsYa4UiSHEUcpu8=";
		};

		devShells.${system}.default = pkgs.mkShell {
			packages = [ pkgs.go ];
		};

		checks.${system} = {
			docs = pkgs.runCommand "docs" {} ''
				cd ${./.}
				${edil}/bin/gendocs -check
				touch $out
			'';
			test = goCheck "test" "go test ./...";
			lint = goCheck "lint" "go vet -vettool=${edil}/bin/edil ./...";
		};
	};
}
