// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package nf

// Class represents a class in the Nerd Fonts project.
type Class string

// String returns the name of the class.
func (c Class) String() string {
	return string(c)
}

// Glyph represents a glyph in the Nerd Fonts project.
type Glyph string

// String returns the raw character of the glyph.
func (g Glyph) String() string {
	return string(g)
}

// IsZero returns true if the glyph is the zero value (empty string).
func (g Glyph) IsZero() bool {
	return g == ""
}
