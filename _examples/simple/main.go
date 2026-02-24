// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

//nolint:forbidigo
package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/lrstanley/go-nf"
	"github.com/lrstanley/go-nf/glyphs/all"
	"github.com/lrstanley/go-nf/glyphs/md"
)

var files = []string{
	"README.markdown",
	"packages.csv",
	"main.xml",
	"global.css",
	"index.html",
	"package.json",
	"report.wtpy",
}

func main() {
	// If you don't want to explicitly enable/disable usage of Nerd Fonts, you can
	// use the `nf.DetectInstalled` function to check if Nerd Fonts are installed
	// on the system.
	switch status, _ := nf.DetectInstalled(context.Background()); status {
	case nf.StatusEnabled:
		fmt.Println("status: nerd fonts were confirmed in use or explicitly enabled by the user")
	case nf.StatusInstalled:
		fmt.Println("status: nerd fonts are installed on the system")
	case nf.StatusDisabled, nf.StatusNotInstalled:
		fmt.Println("status: nerd fonts are not installed on the system or explicitly disabled by the user")
	}

	// Simple static references.
	fmt.Printf("%s beginning tests...", md.TestTube)
	time.Sleep(500 * time.Millisecond)
	fmt.Printf(" -- %s success!\n", md.Check)

	// Crude example that likely won't match exactly what you want, but shows how
	// to dynamically query glyphs.
	for _, file := range files {
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file), "."))

		glyph := md.FileQuestion
		for id := range all.GlyphFullIDs() {
			if strings.Contains(strings.ToLower(id), ext) {
				glyph = all.ByID(id)
				break
			}
		}
		fmt.Printf("--> processing %s (%s)...\n", file, glyph)
	}
	fmt.Printf("%s success!\n", md.Check)
}
