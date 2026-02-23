// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

//go:build windows

package nf

import (
	"context"
	"strings"
	"syscall"
	"unsafe"
)

// logFontW matches the Windows LOGFONTW structure (wingdi.h).
// See: https://learn.microsoft.com/en-us/windows/win32/api/wingdi/ns-wingdi-logfontw
type logFontW struct {
	LfHeight         int32
	LfWidth          int32
	LfEscapement     int32
	LfOrientation    int32
	LfWeight         int32
	LfItalic         uint8
	LfUnderline      uint8
	LfStrikeOut      uint8
	LfCharSet        uint8
	LfOutPrecision   uint8
	LfClipPrecision  uint8
	LfQuality        uint8
	LfPitchAndFamily uint8
	LfFaceName       [32]uint16
}

// enumLogFontExW matches the Windows ENUMLOGFONTEXW structure (wingdi.h).
// See: https://learn.microsoft.com/en-us/windows/win32/api/wingdi/ns-wingdi-enumlogfontexw
type enumLogFontExW struct {
	ElfLogFont  logFontW
	ElfFullName [64]uint16 // LF_FULLFACESIZE
	ElfStyle    [32]uint16
	ElfScript   [32]uint16
}

// DetectorWindowsGDI is a detector that uses the Windows GDI EnumFontFamiliesEx
// API to enumerate installed fonts and match them against Nerd Font patterns.
// Returns [StatusInstalled] if any font name matches, [StatusNotInstalled] otherwise.
//
// This is a no-op on non-Windows systems.
func DetectorWindowsGDI() InstallDetector {
	return func(_ context.Context) (InstallStatus, error) {
		var usr, gdi *syscall.DLL
		var err error

		usr, err = syscall.LoadDLL("C:\\Windows\\System32\\user32.dll")
		if err != nil {
			return StatusNotInstalled, nil
		}

		var getDC, releaseDC, enumFonts *syscall.Proc

		// Required to get a device context for the screen.
		getDC, err = usr.FindProc("GetDC")
		if err != nil {
			return StatusNotInstalled, nil
		}

		// Required to release the device context.
		releaseDC, err = usr.FindProc("ReleaseDC")
		if err != nil {
			return StatusNotInstalled, nil
		}

		gdi, err = syscall.LoadDLL("C:\\Windows\\System32\\gdi32.dll")
		if err != nil {
			return StatusNotInstalled, nil
		}

		// The "W" version handles Unicode character sets, supporting international characters,
		// while "A" uses ANSI.
		enumFonts, err = gdi.FindProc("EnumFontFamiliesExW")
		if err != nil {
			return StatusNotInstalled, nil
		}

		// Get the device context handler for the screen.
		hdc, _, _ := getDC.Call(0)
		if hdc == 0 {
			return StatusNotInstalled, nil
		}
		defer releaseDC.Call(0, hdc)

		var found bool
		callback := syscall.NewCallback(func(lpElfe *enumLogFontExW, _ uintptr, _ uintptr, _ uintptr) uintptr {
			fontName := strings.TrimSpace(syscall.UTF16ToString(lpElfe.ElfLogFont.LfFaceName[:]))
			if fontName == "" || strings.HasPrefix(fontName, "@") {
				return 1
			}
			for _, matcher := range fontNameMatchers {
				if matcher.MatchString(fontName) {
					found = true
					return 0
				}
			}
			return 1
		})

		lf := logFontW{
			// DEFAULT_CHARSET
			// LfCharSet: 0xFF,
		}

		// Trigger the enumeration of font families, and have it invoke our
		// callback function for each font family found.
		_, _, _ = enumFonts.Call(hdc, uintptr(unsafe.Pointer(&lf)), callback, 0, 0)

		if found {
			return StatusInstalled, nil
		}
		return StatusNotInstalled, nil
	}
}

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
