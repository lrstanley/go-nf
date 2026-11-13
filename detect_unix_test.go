// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

//go:build unix

package nf

import (
	"context"
	"errors"
	"testing"
)

func TestDetectorFilesystemContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	got, err := DetectorFilesystem()(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got error %v, want %v", err, context.Canceled)
	}
	if got != StatusNotInstalled {
		t.Fatalf("got %v, want %v", got, StatusNotInstalled)
	}
}
