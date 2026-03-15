package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type fileChangedMsg struct{}

// StartWatcher watches for .go file changes and sends fileChangedMsg to the program.
// Returns a cleanup function.
func StartWatcher(paths []string, notify chan<- fileChangedMsg) (func(), error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// recursively add directories
	for _, p := range paths {
		filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				// skip hidden dirs and vendor
				name := info.Name()
				if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" {
					return filepath.SkipDir
				}
				watcher.Add(path)
			}
			return nil
		})
	}

	go func() {
		var debounce *time.Timer
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// only react to .go file writes
				if !strings.HasSuffix(event.Name, ".go") {
					continue
				}
				if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
					continue
				}
				// debounce: wait 500ms after last change
				if debounce != nil {
					debounce.Stop()
				}
				debounce = time.AfterFunc(500*time.Millisecond, func() {
					notify <- fileChangedMsg{}
				})
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()

	cleanup := func() {
		watcher.Close()
	}
	return cleanup, nil
}
