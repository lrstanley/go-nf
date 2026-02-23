// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

//go:build !windows && !unix

package nf

import "context"

// DetectorFilesystem is a detector that checks the filesystem for font files,
// in common locations. Returns [StatusInstalled] if any font files are found,
// [StatusNotInstalled] otherwise. Permissions errors are ignored.
//
// This is a no-op on non-Unix systems.
func DetectorFilesystem() InstallDetector {
	return func(_ context.Context) (InstallStatus, error) {
		return StatusNotInstalled, nil // No-op.
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
// [StatusNotInstalled] otherwise.
//
// This is a no-op on non-Unix systems.
func DetectorFontConfig() InstallDetector {
	return func(_ context.Context) (InstallStatus, error) {
		return StatusNotInstalled, nil // No-op.
	}
}
