![Apache-2.0](https://img.shields.io/badge/License-Apache--2.0-green.svg) ![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8)

# P2P Sync

> *"Your files. Your LAN. No cloud, no accounts, no telemetry. Just sync."*

A peer-to-peer file synchronization tool designed for local networks. No central server. No cloud dependency. Two machines on the same WiFi see each other, shake hands, and keep their folders in sync.

## Comparison

| Feature | P2P Sync | Syncthing | Dropbox |
|---------|----------|-----------|---------|
| Central server | No | Relay only* | Yes |
| Cloud storage | No | No | Yes |
| LAN-first | Yes | Fallback | No |
| Setup complexity | Minimal | Moderate | Account-based |
| Delta sync | v2+ | Yes | Yes |
| TUI monitor | v3 | Web UI | Web UI |

*Syncthing uses relays for NAT traversal. P2P Sync is LAN-only by design (v4 adds optional relay).

## How It Works

1. **Discovery**: Each node broadcasts itself via mDNS on the local network
2. **Connection**: Peers connect over TCP+TLS using self-signed certificates
3. **Sync Engine**:
   - Compare file hashes (BLAKE3)
   - If files differ, compute delta via rsync-style rolling hash
   - Transfer only changed blocks
4. **Conflict Resolution**: Last-write-wins by default; conflicts saved as `.sync-conflict-*.ext`
5. **Monitoring**: Optional TUI shows sync status, bandwidth, history

```
[ Laptop A ]  <--mDNS-->  [ Laptop B ]
     |                         |
     +---- TCP+TLS tunnel -----+
     |                         |
  ~/Documents  <===========>  ~/Documents
```

## Install

```bash
go install github.com/yourname/p2p-sync@latest
# or
cd p2p-sync
go build -o p2psync ./cmd/p2psync
```

Requires Go 1.22+.

## Usage

```bash
# Start a sync node
p2psync --dir ~/Documents --port 9001

# On another machine
p2psync --dir ~/Documents --port 9001 --peer 192.168.1.42:9001

# With TUI monitor
p2psync --dir ~/Documents --tui
```

## Trust Model

First connection to a peer triggers a fingerprint verification:

```
Peer discovered: neon-laptop @ 192.168.1.42:9001
TLS Fingerprint: SHA256:abc123...
Accept this peer? [y/N] y
Trust saved to ~/.config/p2p-sync/known_peers.json
```

Subsequent connections are automatic.

---

*Sync happens in the shadows of your network. The cloud never knows.*

---

## ☕ Support the Developer

If this project saved you time, solved a problem, or just made your day a little more neon, you can fuel the next one:

[![Buy Me A Coffee](https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png)](https://buymeacoffee.com/synthalorian)
