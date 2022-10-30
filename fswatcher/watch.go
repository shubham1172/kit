package fswatcher

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watch for changes to a directory on the filesystem and sends a notification to notifyCh every time a file in the folder is changed.
// Although it's possible to watch for individual files, that's not recommended; watch for the file's parent folder instead.
// Note that changes are batched for 0.5 seconds before notifications are sent
func Watch(ctx context.Context, dir string, notifyCh chan<- struct{}) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	err = watcher.Add(dir)
	if err != nil {
		return fmt.Errorf("watcher error: %w", err)
	}

	eventCh := make(chan struct{}, 1)
	defer close(eventCh)

	go startPublishEvents(ctx, eventCh, notifyCh)

	for {
		select {
		// Watch for events
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write {
				if strings.Contains(event.Name, dir) {
					eventCh <- struct{}{}
				}
			}

		// Abort in case of errors
		case err = <-watcher.Errors:
			return fmt.Errorf("watcher listen error: %w", err)

		// Stop on context canceled
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// startPublishEvents sends a notification to notifyCh every time an event is received on eventCh.
// This should be run in a goroutine.
func startPublishEvents(ctx context.Context, eventCh <-chan struct{}, notifyCh chan<- struct{}) {
	shouldPublish := false
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-eventCh:
			shouldPublish = true
		case <-ticker.C:
			if shouldPublish {
				notifyCh <- struct{}{}
				shouldPublish = false
			}
		case <-ctx.Done():
			return
		}
	}
}
