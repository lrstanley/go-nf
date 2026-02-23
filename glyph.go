// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package nf

// Glyph represents a glyph in the Nerd Fonts project.
type Glyph struct {
	ID      string // ID of the glyph.
	Class   string // Class name of the glyph (e.g. "fa", "md", "pl").
	Char    string // Raw character of the glyph.
	Unicode string // Unicode codepoint(s) of the glyph.
}

// String returns the raw character of the glyph.
func (g *Glyph) String() string {
	return g.Char
}

// FullID returns the full ID of the glyph ("[class]-[id]").
func (g *Glyph) FullID() string {
	return g.Class + "-" + g.ID
}
