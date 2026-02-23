// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.
//
// Example modified from the following:
//  * https://github.com/charmbracelet/bubbletea/blob/test/examples/package-manager/main.go

package main

import (
	"fmt"
	"math/rand/v2"
	"regexp"

	"github.com/lrstanley/go-nf/glyphs/dev"
	"github.com/lrstanley/go-nf/glyphs/fa"
	"github.com/lrstanley/go-nf/glyphs/md"
)

var versionSuffix = regexp.MustCompile(`-\d+\.\d+\.\d+$`)

var packages = map[string]string{
	"vegeutils":        md.Carrot.String(),
	"libgardening":     md.Leaf.String(),
	"currykit":         fa.PepperHot.String(),
	"spicerack":        md.SpoonSugar.String(),
	"fullenglish":      md.EggFried.String(),
	"eggy":             md.Egg.String(),
	"bad-kitty":        md.Cat.String(),
	"chai":             md.Tea.String(),
	"hojicha":          md.Tea.String(),
	"libtacos":         md.Taco.String(),
	"babys-monads":     md.CodeBraces.String(),
	"libpurring":       md.Cat.String(),
	"currywurst-devel": md.FoodSteak.String(),
	"xmodmeow":         md.Cat.String(),
	"licorice-utils":   md.Candy.String(),
	"cashew-apple":     md.FoodApple.String(),
	"rock-lobster":     md.Fish.String(),
	"standmixer":       md.Blender.String(),
	"coffee-CUPS":      md.Coffee.String(),
	"libesszet":        md.Application.String(),
	"zeichenorientierte-benutzerschnittstellen": md.Application.String(),
	"schnurrkit":      md.Cat.String(),
	"old-socks-devel": fa.Socks.String(),
	"jalapeño":        fa.PepperHot.String(),
	"molasses-utils":  md.SpoonSugar.String(),
	"xkohlrabi":       md.Carrot.String(),
	"party-gherkin":   dev.Cucumber.String(),
	"snow-peas":       md.Leaf.String(),
	"libyuzu":         fa.Lemon.String(),
}

// GlyphForPackage returns the glyph for a package, or a default package glyph if unknown.
func GlyphForPackage(pkg string) string {
	base := versionSuffix.ReplaceAllString(pkg, "")
	if g, ok := packages[base]; ok {
		return g
	}
	return md.Package.String()
}

func getPackages() []string {
	pkgs := make([]string, 0, len(packages))
	for name := range packages {
		pkgs = append(pkgs, name)
	}

	rand.Shuffle(len(pkgs), func(i, j int) {
		pkgs[i], pkgs[j] = pkgs[j], pkgs[i]
	})

	for k := range pkgs {
		pkgs[k] += fmt.Sprintf("-%d.%d.%d", rand.IntN(10), rand.IntN(10), rand.IntN(10)) //nolint:gosec
	}
	return pkgs
}
