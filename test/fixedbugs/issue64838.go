// errorcheck

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p

type I interface{}
type T struct{ X *I }

func g() {
	var i I
	var t = T{i} // ERROR "cannot use i \(variable of type I\) as *I value in struct literal: type *I is pointer to interface, not interface"
	_ = i
}
