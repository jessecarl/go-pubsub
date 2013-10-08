// Copyright 2012 Jesse Allen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

import (
	"testing"
)

func TestNew(t *testing.T) {
	// almost too dumb to test, but we will anyhow
	ps := New()
	if ps == nil {
		t.Errorf("New() returned nil pointer")
	}
}

func TestRegister(t *testing.T) {
	t.Errorf("Tests not created")
}
