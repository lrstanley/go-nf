// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"encoding/json"
	"iter"
	"maps"
	"slices"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

var val = validator.New(validator.WithRequiredStructEnabled())

// Data is the main data structure for the glyph data.
type Data struct {
	Metadata *Metadata           `json:"metadata" validate:"required"`
	Glyphs   map[string][]*Glyph `json:"glyphs" validate:"required,min=5,dive,required"`
}

// Classes returns all the classes in the data (sorted).
func (d *Data) Classes() []string {
	classes := slices.Collect(maps.Keys(d.Glyphs))
	slices.Sort(classes)
	return classes
}

// AllIter returns an iterator over all the glyphs in the data (sorted).
func (d *Data) AllIter() iter.Seq[*Glyph] {
	return func(yield func(g *Glyph) bool) {
		for _, class := range d.Classes() {
			for _, g := range d.Glyphs[class] {
				if !yield(g) {
					return
				}
			}
		}
	}
}

// UnmarshalJSON converts the JSON data into the Data structure. The JSON structure is:
//
//	{
//	  "METADATA": {
//	    "website": "[url]",
//	    "development-website": "[url]",
//	    "version": "[semver-version]",
//	    "date": "[date]"
//	  },
//	  "[class]-[id]": {
//	    "char": "[raw-character]",
//	    "code": "[hex codepoint, utf-8 or utf-16]"
//	  },
//	  [...]
//	}
func (d *Data) UnmarshalJSON(b []byte) error {
	d.Metadata = &Metadata{}
	d.Glyphs = make(map[string][]*Glyph, 5000) // pre-allocate for ~5k glyphs

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	for k, v := range raw {
		if v == nil {
			continue
		}

		if k == "METADATA" {
			if err := json.Unmarshal(v, d.Metadata); err != nil {
				return err
			}
			continue
		}

		var g *Glyph
		if err := json.Unmarshal(v, &g); err != nil {
			return err
		}

		class, id, ok := strings.Cut(k, "-")
		if !ok {
			logger.Error("invalid glyph key", "key", k) //nolint:all
			continue
		}

		g.ID = id
		g.Class = class
		g.FullID = k
		g.PascalID = strcase.ToCamel(id)

		// If the PascalID starts with a number, prepend it with "Glyph"
		if unicode.IsDigit(rune(g.PascalID[0])) {
			g.PascalID = "Glyph" + g.PascalID
		}

		d.Glyphs[class] = append(d.Glyphs[class], g)
	}

	// Sort all the glyphs by ID.
	for _, glyphs := range d.Glyphs {
		slices.SortFunc(glyphs, func(a, b *Glyph) int {
			return strings.Compare(a.ID, b.ID)
		})
	}

	// Ensure all data is valid.
	if err := val.Struct(d); err != nil {
		return err
	}

	return nil
}

type Metadata struct {
	Website            string `json:"website" validate:"required,url"`
	DevelopmentWebsite string `json:"development-website" validate:"required,url"`
	Version            string `json:"version" validate:"required,semver"`
	Date               string `json:"date" validate:"required,min=10"`
}

type Glyph struct {
	ID       string `json:"id" validate:"required,min=1,max=50"`
	FullID   string `json:"full_id" validate:"required,min=1,max=55"`
	PascalID string `json:"pascal_id" validate:"required,min=1,max=50"`
	Class    string `json:"class" validate:"required,min=1,max=10"`
	Char     string `json:"char" validate:"required,min=1,max=10"`
	HexCode  string `json:"code" validate:"required,min=1,max=10"`
}
