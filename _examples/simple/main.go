// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

//nolint:forbidigo
package main

import (
	"context"
	"fmt"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/lrstanley/go-nf"
	"github.com/lrstanley/go-nf/glyphs/md"
	"github.com/lrstanley/go-nf/glyphs/neo"
)

var files = []string{
	"README.md",
	"packages.csv",
	"main.xml",
	"global.css",
	"index.html",
	"package.json",
	"report.wtpy",
}

func printf(format string, a ...any) {
	_, _ = fmt.Printf(format, a...)
}

func main() {
	statusBlock := lipgloss.NewStyle().
		Background(lipgloss.Color("#080808")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Render("status:")

	// If you don't want to explicitly enable/disable usage of Nerd Fonts, you can
	// use the `nf.DetectInstalled` function to check if Nerd Fonts are installed
	// on the system.
	switch status, _ := nf.DetectInstalled(context.Background()); status {
	case nf.StatusEnabled:
		printf("%s nerd fonts were confirmed in use or explicitly enabled by the user\n", statusBlock)
	case nf.StatusInstalled:
		printf("%s nerd fonts are installed on the system\n", statusBlock)
	case nf.StatusDisabled, nf.StatusNotInstalled:
		printf("%s nerd fonts are not installed on the system or explicitly disabled by the user\n", statusBlock)
	}

	// Simple static references.
	cos := neo.CurrentOS()
	printf("%s beginning tests (os: %s)...", md.TestTube, lipgloss.NewStyle().Foreground(cos.Color(true)).Render(cos.String()))
	time.Sleep(500 * time.Millisecond)
	printf(" -- %s %s!\n", md.Check, lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("success"))

	blend := lipgloss.Blend1D(len(files), lipgloss.Color("#9A86FD"), lipgloss.Color("#FF69B4"))

	// Dynamically resolve glyphs based on the file name.
	for i, file := range files {
		g := neo.ByPath(file)
		var icon string
		if g != nil {
			icon = lipgloss.NewStyle().Foreground(g.Color(true)).Render(g.String())
		} else {
			icon = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render(md.FileQuestion.String())
		}

		printf(
			"--> %s (%s)...\n",
			lipgloss.NewStyle().Foreground(blend[i]).Render("processing "+file),
			icon,
		)
	}

	fmt.Printf("%s success!\n", md.Check)
}
