// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package nf

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// fontNameMatchers is the list of regular expressions to match font names.
var fontNameMatchers = []*regexp.Regexp{
	regexp.MustCompile(`(?i)nerd[\s_-]*fonts?`),
	regexp.MustCompile(`(?i)\bnf[pm]?\b`),
}

// InstallStatus represents the detected status of Nerd Fonts being installed.
type InstallStatus int

func (s InstallStatus) String() string {
	switch s {
	case StatusDisabled:
		return "disabled"
	case StatusEnabled:
		return "enabled"
	case StatusNotInstalled:
		return "not installed"
	case StatusInstalled:
		return "installed"
	default:
		return fmt.Sprintf("unknown status: %d", s)
	}
}

const (
	// StatusDisabled indicates that Nerd Fonts was explicitly disabled by the user,
	// detected through one of the [InstallDetector] functions (e.g. explicitly
	// disabled through an environment variable).
	StatusDisabled InstallStatus = iota + 1
	// StatusEnabled indicates that Nerd Fonts was explicitly enabled by the user,
	// detected through one of the [InstallDetector] functions (e.g. explicitly
	// enabled through an environment variable).
	StatusEnabled
	// StatusNotInstalled indicates that none of the [InstallDetector] functions
	// detected that Nerd Fonts is installed on the system.
	StatusNotInstalled
	// StatusInstalled indicates that Nerd Fonts was detected as installed on the
	// system. This does not necessarily mean that the fonts are configured on the
	// terminal emulator, though. Use this to potentially prompt the user to enable
	// Nerd Font glyphs, e.g. "We've detected that Nerd Fonts may be installed on
	// your system, would you like to enable advanced icons?".
	StatusInstalled
)

// InstallDetector is a function that detects whether Nerd Fonts is installed on
// the system. It returns the [InstallStatus] of the detected installation, and
// any errors that occurred during detection. If [StatusEnabled], [StatusDisabled],
// or [StatusInstalled] are returned, the detection process will halt and return,
// not running any further detectors.
type InstallDetector func(ctx context.Context) (InstallStatus, error)

// DefaultDetectors returns a list of default detectors that will be used to
// detect whether Nerd Fonts is installed on the system. If you want to use your
// own detectors, it is recommended to use these as a base, prepending your own,
// or using it as a reference to pick and choose.
func DefaultDetectors() []InstallDetector {
	return []InstallDetector{
		DetectorEnvVar("NERD_FONTS"),
		DetectorEnvVar("NERDFONTS"),
		DetectorEnvVar("NF_FONTS"),
		DetectorWindowsGDI(),
		DetectorFontConfig(),
		DetectorFilesystem(),
	}
}

// DetectInstalled detects whether Nerd Fonts is installed on the system using
// the provided detectors. If no detectors are provided, the default list from
// [DefaultDetectors] will be used. The detection process will halt and return
// once a detector returns a non-StatusNotInstalled status.
//
// It is advised to order your detectors by ones which may return [StatusEnabled]
// or [StatusDisabled] first, and ones which may return [StatusInstalled] last,
// and by least invasive to most invasive (e.g. reading env vars, cli flags, vs
// reading filesystem, running commands, etc.)
//
// Default detectors are:
//   - env vars: NERD_FONTS (preferred), NERDFONTS, NF_FONTS (all platforms)
//   - Windows GDI: uses the Windows GDI API to enumerate installed fonts (windows only)
//   - FontConfig: uses the fontconfig CLI to enumerate installed fonts (unix only)
//   - Filesystem: checks the filesystem for font files in common locations (unix only)
//
// Why is it difficult to detect if Nerd Fonts are being used?: See
// https://github.com/ryanoasis/nerd-fonts/discussions/829
func DetectInstalled(ctx context.Context, detectors ...InstallDetector) (InstallStatus, error) {
	if len(detectors) == 0 {
		detectors = DefaultDetectors()
	}

	var errs []error

	for _, detector := range detectors {
		status, err := detector(ctx)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		switch status { //nolint:exhaustive
		case StatusEnabled, StatusDisabled, StatusInstalled:
			return status, nil
		}
	}

	if len(errs) > 0 {
		return StatusNotInstalled, errors.Join(errs...)
	}
	return StatusNotInstalled, nil
}

// DetectorEnvVar returns an [InstallDetector] that checks for the presence of
// the provided environment variable name, and parses it as a boolean, returning
// [StatusEnabled] if true, [StatusDisabled] if false, and [StatusNotInstalled]
// if the environment variable is not set.
func DetectorEnvVar(name string) InstallDetector {
	return func(_ context.Context) (InstallStatus, error) {
		v := os.Getenv(name)
		if v == "" {
			return StatusNotInstalled, nil
		}
		vv, err := strconv.ParseBool(v)
		if err != nil {
			return StatusNotInstalled, fmt.Errorf("failed to parse env var %q: %w", name, err)
		}
		if vv {
			return StatusEnabled, nil
		}
		return StatusDisabled, nil
	}
}
