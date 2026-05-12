// Package watcher provides file system watching capabilities for .env files.
// It monitors one or more env files for changes and triggers a callback
// when modifications are detected, enabling live diff updates.
package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Event represents a file change notification.
type Event struct {
	// Path is the absolute path of the file that changed.
	Path string
	// ModTime is the modification time at the moment of detection.
	ModTime time.Time
}

// Handler is a function called when a watched file changes.
type Handler func(event Event)

// Watcher monitors a set of files for modifications.
type Watcher struct {
	mu       sync.Mutex
	files    map[string]time.Time // path -> last known mod time
	handler  Handler
	interval time.Duration
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// New creates a Watcher that polls the given files at the specified interval.
// interval controls how frequently the files are stat'd; 500ms is a reasonable default.
func New(paths []string, interval time.Duration, handler Handler) (*Watcher, error) {
	if handler == nil {
		return nil, fmt.Errorf("watcher: handler must not be nil")
	}
	if interval <= 0 {
		return nil, fmt.Errorf("watcher: interval must be positive")
	}

	files := make(map[string]time.Time, len(paths))
	for _, p := range paths {
		abs, err := filepath.Abs(p)
		if err != nil {
			return nil, fmt.Errorf("watcher: resolving path %q: %w", p, err)
		}
		info, err := os.Stat(abs)
		if err != nil {
			return nil, fmt.Errorf("watcher: stat %q: %w", abs, err)
		}
		files[abs] = info.ModTime()
	}

	return &Watcher{
		files:    files,
		handler:  handler,
		interval: interval,
		stopCh:   make(chan struct{}),
	}, nil
}

// Start begins polling in a background goroutine.
// It returns immediately; call Stop to shut down.
func (w *Watcher) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.poll()
			case <-w.stopCh:
				return
			}
		}
	}()
}

// Stop halts the background polling goroutine and waits for it to exit.
func (w *Watcher) Stop() {
	close(w.stopCh)
	w.wg.Wait()
}

// poll checks each watched file for modification and fires the handler if changed.
func (w *Watcher) poll() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for path, lastMod := range w.files {
		info, err := os.Stat(path)
		if err != nil {
			// File may have been temporarily unavailable; skip this cycle.
			continue
		}
		if info.ModTime().After(lastMod) {
			w.files[path] = info.ModTime()
			w.handler(Event{
				Path:    path,
				ModTime: info.ModTime(),
			})
		}
	}
}
