# logdrain

A lightweight CLI for tailing and filtering structured JSON logs from multiple sources.

---

## Installation

```bash
go install github.com/yourusername/logdrain@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logdrain.git
cd logdrain
go build -o logdrain .
```

## Usage

Tail logs from a file and filter by log level:

```bash
logdrain tail --file app.log --filter level=error
```

Pipe logs from multiple sources:

```bash
logdrain tail --file app.log --file worker.log --filter service=api
```

Pretty-print raw JSON log output:

```bash
cat app.log | logdrain fmt
```

### Flags

| Flag | Description |
|------|-------------|
| `--file` | Log file to tail (repeatable) |
| `--filter` | Key=value filter applied to JSON fields |
| `--follow` | Continue watching for new log lines |
| `--level` | Minimum log level to display |

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

## License

[MIT](LICENSE)