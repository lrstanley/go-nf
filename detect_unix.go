// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

//go:build unix

package nf

import (
	"bufio"
	"context"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// fsMaxDepth is the maximum depth to traverse when checking for font files in the filesystem.
const fsMaxDepth = 4

// fontExtensions is the list of font file extensions to check for.
var fontExtensions = []string{
	".ttf",
	".otf",
	".woff",
	".woff2",
}

// filesystemPaths is the list of paths to check for font files in the filesystem.
// Root directory trees only, they will be traversed recursively (ignoring errors).
var filesystemPaths = []string{
	"/usr/share/fonts/",
	"/usr/local/share/fonts/",
	"/var/lib/snapd/desktop/fonts/",
	"~/.fonts/",
	"~/.local/share/fonts/",
}

// DetectorFilesystem is a detector that checks the filesystem for font files,
// in common locations. Returns [StatusInstalled] if any font files are found,
// [StatusNotInstalled] otherwise. Permissions errors are ignored.
//
// This is a no-op on non-Unix systems.
func DetectorFilesystem() InstallDetector { //nolint:gocognit
	home, _ := os.UserHomeDir()
	return func(_ context.Context) (InstallStatus, error) {
		var errs []error
		for _, path := range filesystemPaths {
			if strings.HasPrefix(path, "~/") {
				path = filepath.Join(home, path[2:])
			}

			if _, err := os.Stat(path); err != nil {
				continue
			}

			sepCount := strings.Count(path, string(filepath.Separator))

			var found bool
			err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
				if err != nil {
					return filepath.SkipDir
				}

				if d.IsDir() {
					// Skip directories that are too deep.
					if c := strings.Count(p, string(filepath.Separator)) - sepCount; c > fsMaxDepth {
						return filepath.SkipDir
					}
					return nil
				}

				// Skip non-font files.
				if !slices.Contains(fontExtensions, filepath.Ext(d.Name())) {
					return nil
				}

				// Check if the font name matches any of the matchers.
				for _, matcher := range fontNameMatchers {
					if matcher.MatchString(d.Name()) {
						found = true
						return filepath.SkipAll
					}
				}
				return nil
			})
			if found {
				return StatusInstalled, nil
			}
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}
		if len(errs) > 0 {
			return StatusNotInstalled, errors.Join(errs...)
		}
		return StatusNotInstalled, nil
	}
}

// DetectorWindowsGDI is a detector that uses the Windows GDI EnumFontFamiliesEx
// API to enumerate installed fonts. This is a no-op on non-Windows systems.
func DetectorWindowsGDI() InstallDetector {
	return func(_ context.Context) (InstallStatus, error) {
		return StatusNotInstalled, nil
	}
}

// DetectorFontConfig is a detector that runs fc-list (fontconfig) and checks its
// output for Nerd Fonts. Returns [StatusInstalled] if any font name matches,
// [StatusNotInstalled] otherwise. Skips gracefully (no error) if fc-list is not
// in PATH or cannot be executed.
//
// This is a no-op on non-Unix systems.
func DetectorFontConfig() InstallDetector {
	return func(ctx context.Context) (InstallStatus, error) {
		if _, err := exec.LookPath("fc-list"); err != nil {
			return StatusNotInstalled, err
		}

		cmd := exec.CommandContext(ctx, "fc-list")
		cmd.WaitDelay = 10 * time.Millisecond

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return StatusNotInstalled, err
		}
		err = cmd.Start()
		if err != nil {
			return StatusNotInstalled, err
		}

		scanner := bufio.NewScanner(stdout)
		var line string
		for scanner.Scan() {
			line = scanner.Text()
			for _, matcher := range fontNameMatchers {
				if matcher.MatchString(line) {
					_ = cmd.Wait()
					return StatusInstalled, nil
				}
			}
		}

		_ = cmd.Wait()
		return StatusNotInstalled, nil
	}
}
