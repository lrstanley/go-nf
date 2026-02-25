// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"maps"
	"net/http"
	"slices"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

var val = validator.New(validator.WithRequiredStructEnabled())

// GlyphData is the main data structure for the glyph data.
type GlyphData struct {
	Metadata *Metadata           `json:"metadata" validate:"required"`
	Glyphs   map[string][]*Glyph `json:"glyphs" validate:"required,min=5,dive,required"`
}

// Classes returns all the classes in the data (sorted).
func (d *GlyphData) Classes() []string {
	classes := slices.Collect(maps.Keys(d.Glyphs))
	slices.Sort(classes)
	return classes
}

// AllIter returns an iterator over all the glyphs in the data (sorted).
func (d *GlyphData) AllIter() iter.Seq[*Glyph] {
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

func (d *GlyphData) ByChar(char string) *Glyph {
	for _, class := range d.Classes() {
		for _, g := range d.Glyphs[class] {
			if g.Char == char {
				return g
			}
		}
	}
	return nil
}

var reservedClasses = []string{
	"all",
	"neo",
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
func (d *GlyphData) UnmarshalJSON(b []byte) error {
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

		if slices.Contains(reservedClasses, class) {
			panic(fmt.Sprintf("reserved class: %s", class))
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

func fetchGlyphData(ctx context.Context) (*GlyphData, error) {
	data := &GlyphData{}

	b := readCache(ctx, glyphDataURL)
	if b != nil {
		err := json.Unmarshal(b, data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}
		return data, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, glyphDataURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp *http.Response
	resp, err = httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	writeCache(ctx, glyphDataURL, b)

	err = json.Unmarshal(b, data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return data, nil
}
