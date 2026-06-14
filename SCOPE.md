# Scope of Work — P2P Sync

## v1: Discovery & One-Way Sync

**Goal**: Two machines can find each other and copy files.

- [ ] mDNS service discovery (via `github.com/grandcat/zeroconf` or raw UDP)
  - [ ] Broadcast: node name, sync port, advertised directories
  - [ ] Browse: list peers on LAN
- [ ] TCP transport with TLS (self-signed certs, fingerprint pinning)
  - [ ] Handshake protocol (version, capabilities, auth)
- [ ] Full file sync (one-way, push or pull)
  - [ ] File listing with metadata (size, mtime, mode)
  - [ ] Full file transfer over TLS socket
- [ ] Directory watching with `fsnotify`
  - [ ] Detect create / modify / delete / rename
  - [ ] Debounce rapid changes
- [ ] Basic CLI: `--dir`, `--port`, `--peer`

**v1 Acceptance**: Edit a file in `~/Sync` on Machine A. Machine B receives it within 5 seconds.

## v2: Bidirectional Sync & Delta

**Goal**: True two-way sync with minimal bandwidth.

- [ ] Bidirectional sync protocol
  - [ ] Sync both directions on connection
  - [ ] Handle simultaneous edits
- [ ] Conflict detection
  - [ ] Compare mtime + hash
  - [ ] Last-write-wins default strategy
  - [ ] Conflict files: `document.txt.sync-conflict-20240613-205612.txt`
- [ ] Delta sync (rsync algorithm)
  - [ ] Rolling hash (Rabin-Karp style) for block boundaries
  - [ ] Weak hash + strong hash (BLAKE3 per block)
  - [ ] Transfer only differing blocks
  - [ ] Reconstruct file on receiver
- [ ] Directory tree sync (recursive, preserve structure)

## v3: TUI Monitor & Polish

**Goal**: Professional-grade observable sync.

- [ ] Terminal UI using `ratatui` or `bubbletea`
  - [ ] Live peer list with connection status
  - [ ] File transfer progress bars
  - [ ] Sync history log (last N events)
  - [ ] Bandwidth stats (current / average / total)
  - [ ] Per-directory sync status
- [ ] Configuration file (`~/.config/p2p-sync/config.toml`)
  - [ ] Sync directories, ignore patterns, peer aliases
- [ ] ACL / permissions
  - [ ] Read-only vs read-write peers
  - [ ] Per-directory access control
- [ ] Resume interrupted transfers
- [ ] Compression (zstd) for full-file transfers

## v4: NAT Traversal Relay (Optional)

**Goal**: Sync beyond the LAN.

- [ ] Optional relay server (separate binary: `p2p-sync-relay`)
  - [ ] Hole punching coordination (STUN-lite)
  - [ ] Fallback TCP relay for symmetric NAT
- [ ] Peer routing via relay
- [ ] Relay authentication (token-based)

## Architecture

```
+----------+     +-----------+     +-------------+     +-----------+
|  mDNS    | <-> |  Peer     | <-> |  Sync       | <-> |  Filesystem|
| Discovery|     |  Manager  |     |  Engine     |     |  Watcher  |
+----------+     +-----------+     +-------------+     +-----------+
                      |                  |
                      v                  v
               +-----------+      +-------------+
               |  TLS      |      |  Delta /    |
               |  Transport|      |  Hash Engine|
               +-----------+      +-------------+
                      |                  |
                      v                  v
               +-----------+      +-------------+
               |  TUI /    |      |  Config /   |
               |  CLI      |      |  ACL Store  |
               +-----------+      +-------------+
```

## Milestones

| Day | Deliverable | Acceptance Criteria |
|-----|-------------|-------------------|
| 1   | mDNS Hello | Two machines discover each other on LAN |
| 3   | File Sync | One-way file transfer works, watcher detects changes |
| 7   | Bidirectional | Both directions sync; conflicts create `.sync-conflict` files |
| 14  | TUI + Polish | TUI shows peers, progress, history; config file works |

## Protocol Sketch

```
PEER_HELLO      -> {name, version, dirs[], fingerprint}
PEER_ACCEPT     <- {status: ok}
SYNC_REQUEST    -> {dir_id, file_list[]}
SYNC_RESPONSE   <- {needed_files[], deltas[]}
FILE_BLOCK      -> {path, offset, data, hash}
FILE_DONE       <- {path, status: ok}
```

All messages are length-prefixed protobuf or msgpack over TLS.

## Non-Goals

- Cloud storage backend (S3, GCS, etc.)
- Mobile apps (iOS/Android)
- Version history / time machine (only latest sync + conflicts)
- End-to-end encryption beyond TLS (LAN threat model)
- Web UI (TUI is the interface)
