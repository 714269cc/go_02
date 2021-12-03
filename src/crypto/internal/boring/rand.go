// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && amd64 && !android && !cmd_go_bootstrap && !msan
// +build linux,amd64,!android,!cmd_go_bootstrap,!msan

package boring

import "crypto/internal/boring/boringcrypto"

const RandReader = boringcrypto.RandReader(0)
