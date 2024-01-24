// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fiat_test

import (
	"crypto/internal/nistec/fiat"
	"testing"
)

func BenchmarkMul(b *testing.B) {
	b.Run("P224", func(b *testing.B) {
		v := new(fiat.P224Element).One()
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			v.Mul(v, v)
		}
	})
	b.Run("P384", func(b *testing.B) {
		v := new(fiat.P384Element).One()
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			v.Mul(v, v)
		}
	})
	b.Run("P521", func(b *testing.B) {
		v := new(fiat.P521Element).One()
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			v.Mul(v, v)
		}
	})
}

func BenchmarkSquare(b *testing.B) {
	b.Run("P224", func(b *testing.B) {
		v := new(fiat.P224Element).One()
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			v.Square(v)
		}
	})
	b.Run("P384", func(b *testing.B) {
		v := new(fiat.P384Element).One()
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			v.Square(v)
		}
	})
	b.Run("P521", func(b *testing.B) {
		v := new(fiat.P521Element).One()
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			v.Square(v)
		}
	})
}
