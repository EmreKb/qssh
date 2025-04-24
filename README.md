# QSSH – Quick SSH Connector

QSSH is a tiny but powerful **terminal user-interface (TUI)** written in Go that lets you browse the hosts defined in your SSH configuration and connect with a single keystroke.

## ✨ Features

- **Instant discovery** – Parses your `~/.ssh/config` (and any nested `Include` files) to build the list of hosts.
- **Fuzzy navigation** – Arrow keys or `j / k` to move, `Enter` to connect.
- **Rich details** – Shows user, hostname and non-default ports at-a-glance.
- **Keyboard-only** – Works wherever a terminal works (no mouse required).
- **Zero-friction** – Once you hit _Enter_, QSSH disappears and your normal SSH session starts (`qssh` never stays resident).

---

## 📦 Installation

### Go install (requires Go ≥1.20)
```bash
go install github.com/EmreKb/qssh@latest
```

> `go install` builds QSSH from source on your machine.

---

## 🚀 Usage

Just run:
```bash
qssh
```

> QSSH will parse your `~/.ssh/config` (and any files referenced via `Include`) and show all matching hosts.

### Example SSH config
```ssh
# ~/.ssh/config
Host my-server
    HostName xxx.xxx.xxx.xxx
    User root
    IdentityFile ~/.ssh/my-server-identity-file
```

After running `qssh`, you might see something like this:

![qssh screenshot](./assets/qssh.png)

Key bindings:

| Key            | Action                |
| -------------- | --------------------- |
| ↑ / k          | Move cursor up        |
| ↓ / j          | Move cursor down      |
| ↵ Enter / Space| Connect to selection  |
| q / Ctrl+C     | Quit                  |

---

## ⚙️  How it works

QSSH uses:
* [`ssh_config`](https://github.com/kevinburke/ssh_config) to parse your SSH config (including nested `Include` files).
* [`bubbletea`](https://github.com/charmbracelet/bubbletea) and [`lipgloss`](https://github.com/charmbracelet/lipgloss) for the TUI.
* Regular `ssh` under the hood – once a host is chosen, QSSH **replaces** itself with the `ssh` process.

---

## 🛠  Development
1. Clone the repo and `cd qssh`.
2. `go run .` – launches the dev build.

---

## 🤝 Contributing
Bug reports, feature ideas and PRs are welcome!

---

## 📄 License
MIT – see [LICENSE](LICENSE).

---

### Acknowledgements
Thanks to the Charm Bracelet crew for the amazing terminal UI ecosystem and to Kevin Burke for the `ssh_config` parser. 