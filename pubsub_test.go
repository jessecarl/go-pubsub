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
	// TODO [jesse@jessecarl.com][2013-10-08]: create test publishers
	// TODO [jesse@jessecarl.com][2013-10-08]: create test subscribers
	// TODO [jesse@jessecarl.com][2013-10-08]: create test filters
	// TODO [jesse@jessecarl.com][2013-10-08]: create test cases as follows:
	// * basic set of publisher, subscriber, filter combinations
	// * heavy load (large number of publishers, subscribers, etc.)
	// * concurrency (heavy concurrent load)
}
