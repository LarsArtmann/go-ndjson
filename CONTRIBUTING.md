# Contributing

Thanks for your interest in contributing!

## How to Contribute

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## Development Setup

The project uses [Nix flakes](https://nixos.wiki/wiki/Flakes) for build automation.

```bash
nix develop                    # enter devShell (sets GOEXPERIMENT=jsonv2)
nix run .#test                 # run tests
nix run .#test-race            # run tests with race detector
nix run .#lint                 # run golangci-lint
nix run .#vet                  # run go vet
nix flake check                # validate flake + formatting
```

If running raw `go` commands outside the flake, you must set the experiment flag:

```bash
export GOEXPERIMENT=jsonv2     # required for encoding/json/v2
go test ./... -race
golangci-lint run ./...
```

## Reporting Issues

Please use GitHub Issues to report bugs or request features.
