// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package nf

import "testing"

func TestDetectorEnvVar(t *testing.T) {
	tests := []struct {
		value   string
		want    InstallStatus
		wantErr bool
	}{
		{value: "", want: StatusNotInstalled},
		{value: "auto", want: StatusNotInstalled},
		{value: " AUTO ", want: StatusNotInstalled},
		{value: "1", want: StatusEnabled},
		{value: "t", want: StatusEnabled},
		{value: "TRUE", want: StatusEnabled},
		{value: "yes", want: StatusEnabled},
		{value: "y", want: StatusEnabled},
		{value: "On", want: StatusEnabled},
		{value: "0", want: StatusDisabled},
		{value: "f", want: StatusDisabled},
		{value: "false", want: StatusDisabled},
		{value: "No", want: StatusDisabled},
		{value: "n", want: StatusDisabled},
		{value: "off", want: StatusDisabled},
		{value: "maybe", want: StatusNotInstalled, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			t.Setenv("NF_TEST_ENV", tt.value)

			got, err := DetectorEnvVar("NF_TEST_ENV")(t.Context())
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error state: got %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectorTerminal(t *testing.T) {
	envs := []string{
		"KITTY_WINDOW_ID",
		"GHOSTTY_RESOURCES_DIR",
		"WEZTERM_PANE",
		"TERM",
		"TERM_PROGRAM",
	}

	tests := []struct {
		name string
		env  map[string]string
		want InstallStatus
	}{
		{name: "none", env: map[string]string{"TERM": "xterm-256color"}, want: StatusNotInstalled},
		{name: "unrelated-term-program", env: map[string]string{"TERM_PROGRAM": "vscode"}, want: StatusNotInstalled},
		{name: "kitty-window-id", env: map[string]string{"KITTY_WINDOW_ID": "1"}, want: StatusEnabled},
		{name: "kitty-term", env: map[string]string{"TERM": "xterm-kitty"}, want: StatusEnabled},
		{name: "ghostty-resources", env: map[string]string{"GHOSTTY_RESOURCES_DIR": "/usr/share/ghostty"}, want: StatusEnabled},
		{name: "ghostty-term", env: map[string]string{"TERM": "xterm-ghostty"}, want: StatusEnabled},
		{name: "ghostty-term-program", env: map[string]string{"TERM_PROGRAM": "Ghostty"}, want: StatusEnabled},
		{name: "wezterm-pane", env: map[string]string{"WEZTERM_PANE": "0"}, want: StatusEnabled},
		{name: "wezterm-term-program", env: map[string]string{"TERM_PROGRAM": "WezTerm"}, want: StatusEnabled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, k := range envs {
				t.Setenv(k, tt.env[k])
			}

			got, err := DetectorTerminal()(t.Context())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectInstalledPrecedence(t *testing.T) {
	t.Setenv("NF_TEST_ENV", "false")
	t.Setenv("KITTY_WINDOW_ID", "1")

	got, err := DetectInstalled(t.Context(), DetectorEnvVar("NF_TEST_ENV"), DetectorTerminal())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != StatusDisabled {
		t.Fatalf("got %v, want %v", got, StatusDisabled)
	}
}
