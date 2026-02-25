// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var cacheDir = sync.OnceValues(func() (string, error) {
	cdir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("get cache directory: %w", err)
	}

	dir := filepath.Join(cdir, "go-nf")
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return "", fmt.Errorf("create cache directory: %w", err)
	}
	return dir, nil
})

var fileCacheMutex = sync.Mutex{}

func cachePath(base, url string) string {
	hash := md5.Sum([]byte(url))
	return filepath.Join(base, time.Now().Format("20060102")+"-"+hex.EncodeToString(hash[:]))
}

func readCache(_ context.Context, url string) []byte {
	fileCacheMutex.Lock()
	defer fileCacheMutex.Unlock()

	dir, err := cacheDir()
	if err != nil {
		logger.Error("failed to get cache directory", "error", err) //nolint:all
		return nil
	}

	path := cachePath(dir, url)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		logger.Error("failed to read file", "file", path, "error", err) //nolint:all
		return nil
	}
	return b
}

func writeCache(_ context.Context, url string, data []byte) {
	fileCacheMutex.Lock()
	defer fileCacheMutex.Unlock()

	dir, err := cacheDir()
	if err != nil {
		logger.Error("failed to get cache directory", "error", err) //nolint:all
		return
	}

	path := cachePath(dir, url)
	if err := os.WriteFile(path, data, 0o640); err != nil {
		logger.Error("failed to write file", "file", path, "error", err) //nolint:all
	}
	logger.Info("cached file", "file", path) //nolint:all
}
