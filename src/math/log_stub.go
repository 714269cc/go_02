// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build math_pure_go || (!amd64 && !s390x)

package math

const haveArchLog = false

func archLog(x float64) float64 {
	panic("not implemented")
}
