/*
 * Copyright (c) 2023-present Maxim Geraskin
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */
package escapes

import (
	"bytes"
	"io"
	"testing"
)

func Benchmark_YIfLongest1(b *testing.B) {
	x := "x"
	y := "y"
	for i := 0; i < b.N; i++ {
		l := YIfLongest(&x, &y)
		if l == nil {
			b.Fatal("l is nil")
		}
	}
}

func Benchmark_YIfLongest1_noinline(b *testing.B) {
	x := "x"
	y := "y"
	for i := 0; i < b.N; i++ {
		l := YIfLongest_noinline(&x, &y)
		if l == nil {
			b.Fatal("l is nil")
		}
	}
}

var closureCaller = func(f func() int) int {
	return f()
}

func Benchmark_ProvideClosure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ProvideClosure(closureCaller)
	}
}

func Benchmark_ProvideInterfaceMethodAsClosure_localCaller(b *testing.B) {
	closureCaller := func(f func() int) int {
		return f()
	}

	c := &Closure{}
	for i := 0; i < b.N; i++ {
		_ = c.ProvideInterfaceMethodAsClosure(closureCaller)
	}
}

func Benchmark_ProvideInterfaceMethodAsClosure_localCaller_intf(b *testing.B) {
	closureCaller := func(f func() int) int {
		return f()
	}

	var c IClosure = &Closure{}
	for i := 0; i < b.N; i++ {
		_ = c.ProvideInterfaceMethodAsClosure(closureCaller)
	}
}

func Benchmark_ProvideInterfaceMethodAsClosure_globalCaller(b *testing.B) {
	c := &Closure{}
	for i := 0; i < b.N; i++ {
		_ = c.ProvideInterfaceMethodAsClosure(closureCaller)
	}
}

func Benchmark_ProvideFieldAsClosure(b *testing.B) {
	c := &Closure{}
	c.closure = c.Do
	for i := 0; i < b.N; i++ {
		_ = c.ProvideFieldAsClosure(closureCaller)
	}
}

func Benchmark_ReadInt64UsingBinaryRead(b *testing.B) {

	reader := bytes.NewReader([]byte{0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04})

	var err error
	for i := 0; i < b.N; i++ {
		_, err = ReadInt64UsingBinaryRead(reader)
		if err != nil {
			b.Fatal(err)
		}
		_, err = reader.Seek(0, io.SeekStart)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_ReadInt64(b *testing.B) {

	data := []byte{0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04}
	bytebuffer := bytes.NewBuffer(data)

	var err error
	for i := 0; i < b.N; i++ {
		_, err = ReadInt64(bytebuffer)
		if err != nil {
			b.Fatal(err)
		}
		bytebuffer.Reset()
		bytebuffer.Write(data)
	}
}
