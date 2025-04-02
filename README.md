# git-llama

![Go Report Card](https://goreportcard.com/badge/github.com/tlehman/git-llama)

Git Llama provides an offline-first, LLM command line interface to [git source control](https://git-scm.com/) that translates your prompt like:

```sh
git-llama "show me the last 5 commits, just the commit messages"
```

Into the real command:
```sh
git log -5 --pretty=format:%s
```

# Build with Nix

If you have [nix flakes](https://nhttps://wiki.nixos.org/wiki/Flakes) enabled, you can install `git-llama` with:

```sh
nix develop
```

Then run `make`.

# Install
Git Llama depends on [ollama](https://github.com/ollama/ollama).
The first iteration of `git-llama` has a decoupled architecture
that doesn't try to install `ollama` for you. If this is a feature
you really want, please
[file a feature request issue](https://github.com/tlehman/git-llama/issues)
and I will write it for you.

# Architecture
![git llama architecture](./architecture.png)

# TODO
- [x] Check if `ollama` is started in background
- [x] Add SQLite vector database [sqlite-vec](https://github.com/asg017/sqlite-vec)
- [x] Write the vector db module with an interface that lets you insert, update and get the embeddings from the db
- [ ] Write the input -> Embedding Model code
- [ ] Write Semantic Search code
- [ ] Pair every command with it's inverse (if possible) and lots of unit tests

# Appendix

## Check the `sqlite_version` and the `vec_version`
```go
// check the sqlite_version and the vec_version
stmt, _, err := vecdb.DB.Prepare(`SELECT sqlite_version(), vec_version()`)
if err != nil {
	fmt.Printf("failed getting vec_version(): %s\n", err)
	return nil, err
}
stmt.Step()
fmt.Printf("sqlite_version() = %s, vec_version() = %s\n", stmt.ColumnText(0), stmt.ColumnText(1))
```
