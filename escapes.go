/*
 * Copyright (c) 2023-present Maxim Geraskin
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package escapes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"
)

//go:noinline
func ReturnValueParamAddress(escapesToHeap Point) *Point {
	if escapesToHeap.X < 0 {
		return nil
	}
	return &escapesToHeap
}

func Slices() int {
	doesNotEscape := make([]byte, 10000)
	escapes := make([]byte, 100000)
	return len(doesNotEscape) + len(escapes)
}

// `y`: leaking param
func YIfLongest(x, y *string) *string {
	if len(*y) > len(*x) {
		return y
	}
	// `s`: escapes to heap
	s := ""
	return &s
}

//go:noinline
func YIfLongest_noinline(x, y *string) *string {
	if len(*y) > len(*x) {
		return y
	}
	s := ""
	return &s
}

type Point struct {
	X, Y int
}

func ReturnPointerParam(leaking *Point) *Point {
	return leaking
}

// `p` neither leaks nor escapes
func ReturnPointerParamField(p *Point) int {
	return p.X
}

func ReturnSlice(leaking []byte) []byte {
	return leaking
}

// `p` does NOT leak
func SliceLen(p []byte) int {
	return len(p)
}

// ***** Leaking examples **********************************************************************************************

// `a` escapes to heap because of the ReturnSlice() leaking param
func CallReturnSlice() {
	a := make([]byte, 8)
	fmt.Println(ReturnSlice(a))
}

// `a` is kept on the stack
func CallSliceLen(f func([]byte) int) {
	a := make([]byte, 8)
	fmt.Println(SliceLen(a))
}

// ***** Closures ******************************************************************************************************

// `v` and `closure` escape
func ProvideClosure(closureCaller func(func() int) int) int {
	var v int
	closure := func() int {
		v++
		return 2
	}
	return closureCaller(closure)
}

type IClosure interface {
	Do() int
	ProvideInterfaceMethodAsClosure(closureCaller func(func() int) int) int
	ProvideFieldAsClosure(func(func() int) int) int
}

type Closure struct {
	v       int
	closure func() int
}

func (c *Closure) Do() int {
	c.v++
	return 2
}

// c.Do escapes
func (c *Closure) ProvideInterfaceMethodAsClosure(closureCaller func(func() int) int) int {
	return closureCaller(c.Do)
}

// c.closure does NOT escape
func (c *Closure) ProvideFieldAsClosure(closureCaller func(func() int) int) int {
	return closureCaller(c.closure)
}

// ***** binary.Read ******************************************************************************************************

func ReadInt64UsingBinaryRead(r io.Reader) (int64, error) {
	var v int64
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func ReadInt64(buf *bytes.Buffer) (res int64, err error) {
	res, err = int64(binary.BigEndian.Uint64(buf.Next(int(unsafe.Sizeof(res))))), nil
	return
}
