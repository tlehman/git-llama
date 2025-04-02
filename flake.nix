{
  description = "go lang, SQLite, ollama, deepseek-r1 model";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }: {
    devShells.aarch64-darwin.default = nixpkgs.legacyPackages.aarch64-darwin.mkShell {
      buildInputs = with nixpkgs.legacyPackages.aarch64-darwin; [
        go
        sqlite
        ollama
      ];
      shellHook = ''
        ollama serve 2> /dev/null &
      '';
    };
  };
}
