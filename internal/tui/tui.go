package tui

import (
	"fmt"
	"sync"
)

type TUI struct {
	mu     sync.Mutex
	peers  []string
	events []string
}

func New() *TUI {
	return &TUI{}
}

func (t *TUI) AddPeer(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.peers = append(t.peers, name)
	fmt.Printf("[peer] %s connected\n", name)
}

func (t *TUI) RemovePeer(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, p := range t.peers {
		if p == name {
			t.peers = append(t.peers[:i], t.peers[i+1:]...)
			break
		}
	}
	fmt.Printf("[peer] %s disconnected\n", name)
}

func (t *TUI) LogEvent(event string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, event)
	if len(t.events) > 100 {
		t.events = t.events[1:]
	}
	fmt.Printf("[sync] %s\n", event)
}

func (t *TUI) Status() {
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Printf("Peers: %v\n", t.peers)
}
