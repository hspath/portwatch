# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected listeners.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a baseline of allowed ports:

```bash
portwatch --allow 22,80,443 --interval 30s
```

Run a one-time scan and print all open ports:

```bash
portwatch scan
```

Watch for unexpected listeners and send an alert to a webhook:

```bash
portwatch --allow 22,80,443 --webhook https://hooks.example.com/alert
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--allow` | `""` | Comma-separated list of permitted ports |
| `--interval` | `60s` | How often to poll open ports |
| `--webhook` | `""` | URL to POST alerts to |
| `--verbose` | `false` | Enable verbose logging |

---

## How It Works

`portwatch` periodically reads active listening sockets from the system, compares them against your defined allowlist, and fires an alert whenever an unexpected port is detected. It runs as a lightweight background daemon with minimal resource overhead.

---

## Requirements

- Go 1.21+
- Linux / macOS

---

## License

MIT © 2024 yourusername