package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/synth/p2p-sync/internal/transport"
)

type Sync struct {
	Dir      string
	OnChange func(files []transport.FileInfo)

	watcher *fsnotify.Watcher
	stop    chan struct{}
}

func New(dir string) (*Sync, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := watcher.Add(absDir); err != nil {
		watcher.Close()
		return nil, err
	}
	return &Sync{
		Dir:     absDir,
		watcher: watcher,
		stop:    make(chan struct{}),
	}, nil
}

func (s *Sync) Start() {
	go s.watch()
}

func (s *Sync) Stop() {
	close(s.stop)
	s.watcher.Close()
}

func (s *Sync) Scan() ([]transport.FileInfo, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		return nil, err
	}
	var files []transport.FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		path := filepath.Join(s.Dir, name)
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		hash, err := hashFile(path)
		if err != nil {
			continue
		}
		files = append(files, transport.FileInfo{
			Name: name,
			Hash: hash,
			Size: info.Size(),
		})
	}
	return files, nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (s *Sync) watch() {
	for {
		select {
		case <-s.stop:
			return
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}
			if event.Op.Has(fsnotify.Write) || event.Op.Has(fsnotify.Create) || event.Op.Has(fsnotify.Remove) || event.Op.Has(fsnotify.Rename) {
				files, err := s.Scan()
				if err == nil && s.OnChange != nil {
					s.OnChange(files)
				}
			}
		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("fsnotify error: %v", err)
		}
	}
}
