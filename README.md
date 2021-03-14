# tents

This is a solver for the [Tents and Trees](https://www.brainbashers.com/tents.asp) problem.

## Run

### Development

The setup uses [CompileDaemon](https://github.com/githubnemo/CompileDaemon) for
automatically rebuilding and running on file change.

```bash
$ make build-dev
$ make dev command="./tents -v <input.json>"
```

### Plain

```bash
$ make build
$ make -B <input.json>
```
