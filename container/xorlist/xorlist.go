package xorlist

import "unsafe"

type Element struct {
	P     uintptr
	Value interface{}
}

func uptr(p *Element) uintptr {
	return uintptr(unsafe.Pointer(p))
}
func ptr(u uintptr) *Element {
	return (*Element)(unsafe.Pointer(u))
}
func (e *Element) Prev(next *Element) *Element {
	if e == nil || e.P == 0 {
		return nil
	}

	prev := ptr(uptr(next) ^ e.P)
	if prev != nil && ptr(prev.P) == e {
		return nil
	}
	return prev
}
func (e *Element) Next(prev *Element) *Element {
	if e == nil || e.P == 0 {
		return nil
	}

	next := ptr(uptr(prev) ^ e.P)
	if next != nil && ptr(next.P) == e {
		return nil
	}
	return next
}

type XorList struct {
	head Element
	tail Element
	len  int
}
