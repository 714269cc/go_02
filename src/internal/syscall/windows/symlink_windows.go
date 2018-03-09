// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package windows

import "syscall"

const (
	ERROR_INVALID_PARAMETER syscall.Errno = 87

	// symlink support for CreateSymbolicLink() starting with Windows 10 (1607, v10.0.14393)
	SYMBOLIC_LINK_FLAG_ALLOW_UNPRIVILEGED_CREATE = 0x2
)
