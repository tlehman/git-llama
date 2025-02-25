# git-llama

Git Llama provides an offline-first, command line interface to [git source control](https://git-scm.com/) that translates your prompt like:

```sh
git-llama "show me the last 5 commits, just the commit messages"
```

Into the real command:
```sh
git log -5 --pretty=format:%s
```

# Architecture
![git llama architecture](./architecture.png)

# TODO
- [ ] Start `ollama` in background if not already started
- [ ]
