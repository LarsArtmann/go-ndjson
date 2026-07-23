{
  description = "go-ndjson — NDJSON reader with format detection for Go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    systems.url = "github:nix-systems/default";
  };

  outputs =
    inputs@{
      self,
      flake-parts,
      treefmt-nix,
      systems,
      ...
    }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import systems;

      imports = [
        treefmt-nix.flakeModule
      ];

      perSystem =
        {
          config,
          pkgs,
          ...
        }:
        let
          inherit (pkgs) lib;
          goPkg = pkgs.go_1_26;

          goFiles = lib.fileset.fileFilter (file: file.hasExt "go") ./.;
          src = lib.fileset.toSource {
            root = ./.;
            fileset = lib.fileset.unions [
              ./go.mod
              goFiles
            ];
          };

          mkApp =
            name: runtimeInputs: text:
            let
              script = pkgs.writeShellApplication {
                inherit name runtimeInputs text;
              };
            in
            {
              type = "app";
              program = lib.getExe script;
            };
        in
        {
          treefmt = {
            projectRootFile = "go.mod";
            programs = {
              gofumpt.enable = true;
              goimports.enable = true;
              golines.enable = true;
              nixfmt.enable = true;
            };
          };

          checks.format = config.treefmt.build.check self;

          devShells.default = pkgs.mkShell {
            packages = [
              goPkg
              pkgs.golangci-lint
              pkgs.gofumpt
              pkgs.golines
              pkgs.gopls
              pkgs.gotools
              pkgs.govulncheck
              pkgs.trash-cli
            ];

            GOWORK = "off";
            GOEXPERIMENT = "jsonv2";

            shellHook = ''
              echo "go-ndjson dev shell — $(go version)"
              echo "GOEXPERIMENT=$GOEXPERIMENT"
            '';
          };

          apps = {
            test = mkApp "test" [ goPkg ] ''
              export GOEXPERIMENT=jsonv2
              go test ./... -count=1 "$@"
            '';

            test-race = mkApp "test-race" [ goPkg ] ''
              export GOEXPERIMENT=jsonv2
              go test ./... -race -count=1 "$@"
            '';

            build = mkApp "build" [ goPkg ] ''
              export GOEXPERIMENT=jsonv2
              go build ./...
            '';

            vet = mkApp "vet" [ goPkg ] ''
              export GOEXPERIMENT=jsonv2
              go vet ./...
            '';

            lint = mkApp "lint" [ pkgs.golangci-lint ] ''
              export GOEXPERIMENT=jsonv2
              golangci-lint run ./...
            '';

            coverage = mkApp "coverage" [ goPkg ] ''
              export GOEXPERIMENT=jsonv2
              go test ./... -coverprofile=coverage.out -covermode=atomic "$@"
              go tool cover -func=coverage.out
            '';

            vulncheck = mkApp "vulncheck" [ pkgs.govulncheck ] ''
              export GOEXPERIMENT=jsonv2
              govulncheck ./...
            '';

            clean =
              mkApp "clean"
                [
                  goPkg
                  pkgs.trash-cli
                ]
                ''
                  trash-put coverage.out 2>/dev/null || true
                  go clean -testcache
                '';
          };
        };
    };
}
