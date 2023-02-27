package mem

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

//return pointer need free
func Malloc(size int, typeSize int) unsafe.Pointer {
	return C.malloc(C.ulonglong(size * typeSize))
}

//return pointer need free
func Realloc(p unsafe.Pointer, size int, typeSize int) unsafe.Pointer {
	return C.realloc(p, C.ulonglong(size*typeSize))
}

//free pointer
func Free(p unsafe.Pointer) {
	C.free(p)
}

//Sizeof get type size
func Sizeof(a any) int {
	return int(reflect.TypeOf(a).Size())
}

//Push push p address to pointer, return pointer need free
func Push(p unsafe.Pointer) unsafe.Pointer {
	r := Malloc(1, 8)
	*(*unsafe.Pointer)(r) = p
	return r
}

//PushTo push src address to dst
func PushTo(dst unsafe.Pointer, src unsafe.Pointer) {
	*(*unsafe.Pointer)(dst) = src
}

//PushAt push src address to dst next n ptr
func PushAt(dst unsafe.Pointer, n int, src unsafe.Pointer) {
	p := unsafe.Pointer(uintptr(dst) + uintptr(n)*8)
	*(*unsafe.Pointer)(p) = src
}

//Pop pop a pointer that address saved in p
func Pop(p unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(p)
}

//PopAt pop a pointer that address saved in p[n]
func PopAt(p unsafe.Pointer, n int) unsafe.Pointer {
	return Pop(unsafe.Pointer(uintptr(p) + uintptr(n)*8))
}

//Slice returns a pointer to a slice
func Slice(p unsafe.Pointer, n int) unsafe.Pointer {
	s := &struct {
		addr unsafe.Pointer
		len  int
		cap  int
	}{addr: p, len: n, cap: n}
	return unsafe.Pointer(s)
}

func RangePop(p unsafe.Pointer, n int, fn func(p unsafe.Pointer, n int)) {
	ps := *(*[]unsafe.Pointer)(Slice(p, n))
	for i := range ps {
		fn(ps[i], i)
	}
}
func RangePtr(p unsafe.Pointer, n int, fn func(p unsafe.Pointer, n int)) {
	for i := 0; i < n; i++ {
		_p := unsafe.Pointer(uintptr(p) + uintptr(i)*8)
		fn(_p, i)
	}
}

//IdxPtr return a pointer that is p[n]
func IdxPtr(p unsafe.Pointer, n int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + uintptr(n)*8)
}

//Memset copy nbyte char from src to dest
func Memcpy(dest, src unsafe.Pointer, nbyte int) { C.memcpy(dest, src, C.ulonglong(nbyte)) }

//Memset full mem with byteV,
func Memset(p unsafe.Pointer, byteV byte, byteN int) { C.memset(p, C.int(byteV), C.ulonglong(byteN)) }

type Str C.char

func NewStr(str string) *Str       { return (*Str)(C.CString(str)) }
func (s *Str) _toc() *C.char       { return (*C.char)(s) }
func (s *Str) ToC() unsafe.Pointer { return unsafe.Pointer(s) }
func (s *Str) Free()               { Free(unsafe.Pointer(s)) }
func (s *Str) Len() int            { return int(int32(C.strlen(s._toc()))) }
