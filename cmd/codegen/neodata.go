// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/Shopify/go-lua"
)

const nvimTreeIconBaseURL = "https://raw.githubusercontent.com/nvim-tree/nvim-web-devicons/master/lua/nvim-web-devicons"

type NeoGlyphEntry struct {
	Matcher        string `json:"matcher" validate:"required,min=1,max=50"`
	Name           string `json:"name" validate:"required,min=1,max=50"`
	Glyph          *Glyph `json:"glyph" validate:"required"`
	DarkColor      string `json:"dark_color" validate:"required,hexcolor"`
	DarkANSIColor  int    `json:"dark_ansicolor" validate:"required,min=0,max=255"`
	LightColor     string `json:"light_color" validate:"required,hexcolor"`
	LightANSIColor int    `json:"light_ansicolor" validate:"required,min=0,max=255"`
}

type NeoData struct {
	DesktopEnvironments []*NeoGlyphEntry `json:"desktop_environments" validate:"required,min=1,dive,required"`
	FileExtensions      []*NeoGlyphEntry `json:"file_extensions" validate:"required,min=1,dive,required"`
	Filenames           []*NeoGlyphEntry `json:"filenames" validate:"required,min=1,dive,required"`
	OperatingSystems    []*NeoGlyphEntry `json:"operating_systems" validate:"required,min=1,dive,required"`
	WindowManagers      []*NeoGlyphEntry `json:"window_managers" validate:"required,min=1,dive,required"`
	Classes             []string         `json:"classes" validate:"required,min=1,dive,required"`
}

type rawLuaIconEntry struct {
	icon       string
	color      string
	ctermColor int
	name       string
}

var luaIconFilenames = map[string]string{
	"desktop_environment": "icons_by_desktop_environment.lua",
	"file_extension":      "icons_by_file_extension.lua",
	"filename":            "icons_by_filename.lua",
	"operating_system":    "icons_by_operating_system.lua",
	"window_manager":      "icons_by_window_manager.lua",
}

// mergeToNeoGlyphEntries takes the default and light raw Lua maps for a category,
// resolves the *Glyph via charIndex, and returns sorted []*NeoGlyphEntry.
func mergeToNeoGlyphEntries(
	defaultMap map[string]rawLuaIconEntry,
	lightMap map[string]rawLuaIconEntry,
	glyphData *GlyphData,
) []*NeoGlyphEntry {
	var entries []*NeoGlyphEntry

	for matcher, dark := range defaultMap {
		g := glyphData.ByChar(dark.icon)
		if g == nil {
			logger.Warn("unmatched nvim-tree icon", "matcher", matcher, "icon", dark.icon) //nolint:all
			continue
		}

		entry := &NeoGlyphEntry{
			Matcher:       matcher,
			Name:          dark.name,
			Glyph:         g,
			DarkColor:     dark.color,
			DarkANSIColor: dark.ctermColor,
		}

		if light, ok := lightMap[matcher]; ok {
			entry.LightColor = light.color
			entry.LightANSIColor = light.ctermColor
		}

		entries = append(entries, entry)
	}

	slices.SortFunc(entries, func(a, b *NeoGlyphEntry) int {
		return strings.Compare(a.Matcher+a.Name, b.Matcher+b.Name)
	})

	return entries
}

// parseLuaIconTable iterates over the Lua table at tableIndex and extracts
// identifier -> rawLuaIconEntry mappings.
func parseLuaIconTable(l *lua.State, tableIndex int) (map[string]rawLuaIconEntry, error) {
	tableIndex = l.AbsIndex(tableIndex)
	result := make(map[string]rawLuaIconEntry)

	l.PushNil()
	for l.Next(tableIndex) {
		key, ok := l.ToString(-2)
		if !ok {
			l.Pop(2)
			continue
		}

		if !l.IsTable(-1) {
			l.Pop(2)
			continue
		}

		entry := rawLuaIconEntry{}
		valueIdx := l.AbsIndex(-1)

		var err error

		for _, field := range []string{"icon", "color", "cterm_color", "name"} {
			l.Field(valueIdx, field)
			if s, sok := l.ToString(-1); sok {
				switch field {
				case "icon":
					entry.icon = s
				case "color":
					entry.color = s
				case "cterm_color":
					entry.ctermColor, err = strconv.Atoi(s)
					if err != nil {
						return nil, fmt.Errorf("invalid cterm color: %w: %s", err, s)
					}
				case "name":
					entry.name = s
				}
			}
			l.Pop(1)
		}

		result[key] = entry
		l.Pop(1)
	}

	return result, nil
}

// executeLuaIconFile loads and executes the Lua content, then parses the
// returned table.
func executeLuaIconFile(l *lua.State, content string) (map[string]rawLuaIconEntry, error) {
	if err := lua.DoString(l, content); err != nil {
		return nil, fmt.Errorf("execute lua: %w", err)
	}

	if l.Top() < 1 {
		return nil, errors.New("lua returned no value")
	}

	if !l.IsTable(-1) {
		return nil, errors.New("lua return value is not a table")
	}

	return parseLuaIconTable(l, -1)
}

// fetchAndParseLuaIconFile fetches the Lua file from the given URL and parses it.
func fetchAndParseLuaIconFile(ctx context.Context, url string) (map[string]rawLuaIconEntry, error) {
	b := readCache(ctx, url)
	if b != nil {
		return executeLuaIconFile(lua.NewState(), string(b))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: status %d", url, resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", url, err)
	}

	writeCache(ctx, url, content)
	return executeLuaIconFile(lua.NewState(), string(content))
}

// FetchNeoGlyphData fetches all nvim-web-devicons Lua files, matches each
// icon character against the GlyphData char index, and returns resolved
// NeoGlyphData.
func FetchNeoGlyphData(ctx context.Context, glyphData *GlyphData) (*NeoData, error) {
	type categoryResult struct {
		defaultMap map[string]rawLuaIconEntry
		lightMap   map[string]rawLuaIconEntry
	}

	categories := make(map[string]*categoryResult, len(luaIconFilenames))

	for category, filename := range luaIconFilenames {
		cr := &categoryResult{}

		for _, variant := range []string{"default", "light"} {
			url := strings.Join([]string{nvimTreeIconBaseURL, variant, filename}, "/")
			m, err := fetchAndParseLuaIconFile(ctx, url)
			if err != nil {
				return nil, fmt.Errorf("parse %s/%s: %w", variant, filename, err)
			}

			if variant == "default" {
				cr.defaultMap = m
			} else {
				cr.lightMap = m
			}
		}

		categories[category] = cr
	}

	data := &NeoData{}
	for category, cr := range categories {
		entries := mergeToNeoGlyphEntries(cr.defaultMap, cr.lightMap, glyphData)

		for _, entry := range entries {
			if !slices.Contains(data.Classes, entry.Glyph.Class) {
				data.Classes = append(data.Classes, entry.Glyph.Class)
			}
		}

		switch category {
		case "desktop_environment":
			data.DesktopEnvironments = entries
		case "file_extension":
			data.FileExtensions = entries
		case "filename":
			data.Filenames = entries
		case "operating_system":
			data.OperatingSystems = entries
		case "window_manager":
			data.WindowManagers = entries
		}
	}

	if err := val.Struct(data); err != nil {
		return nil, fmt.Errorf("validate data: %w", err)
	}

	return data, nil
}
