package libc

import (
	"testing"
	"unsafe"
)

var cxaHit byte

func cxaTestDtor(p unsafe.Pointer) {
	cxaHit = 1
	if p != nil {
		*(*byte)(p) = 7
	}
}

func TestDynamicCastNullSrc(t *testing.T) {
	var dst byte
	if DynamicCast(nil, &dst, &dst, 0) != nil {
		t.Fatal("null src must be null")
	}
}

func TestDynamicCastSI(t *testing.T) {
	baseName := append([]byte("4Base"), 0)
	derName := append([]byte("7Derived"), 0)
	var baseTI struct {
		vptr unsafe.Pointer
		name *byte
	}
	baseTI.vptr = unsafe.Pointer(&ClassTypeInfoVT[2])
	baseTI.name = &baseName[0]
	var derTI struct {
		vptr unsafe.Pointer
		name *byte
		base unsafe.Pointer
	}
	derTI.vptr = unsafe.Pointer(&SIClassTypeInfoVT[2])
	derTI.name = &derName[0]
	derTI.base = unsafe.Pointer(&baseTI)

	var vt [3]unsafe.Pointer
	vt[1] = unsafe.Pointer(&derTI)
	var obj struct{ vptr unsafe.Pointer }
	obj.vptr = unsafe.Pointer(&vt[2])

	got := DynamicCast((*byte)(unsafe.Pointer(&obj)), (*byte)(unsafe.Pointer(&baseTI)), (*byte)(unsafe.Pointer(&derTI)), 0)
	if got != (*byte)(unsafe.Pointer(&obj)) {
		t.Fatalf("downcast %p want %p", got, &obj)
	}

	var baseVT [3]unsafe.Pointer
	baseVT[1] = unsafe.Pointer(&baseTI)
	var baseObj struct{ vptr unsafe.Pointer }
	baseObj.vptr = unsafe.Pointer(&baseVT[2])
	if DynamicCast((*byte)(unsafe.Pointer(&baseObj)), (*byte)(unsafe.Pointer(&baseTI)), (*byte)(unsafe.Pointer(&derTI)), 0) != nil {
		t.Fatal("base as derived must be null")
	}
}

func TestDynamicCastUnknownKindPanics(t *testing.T) {
	var srcTI, dstTI [2]unsafe.Pointer
	var vt [3]unsafe.Pointer
	vt[1] = unsafe.Pointer(&srcTI[0])
	var obj struct{ vptr unsafe.Pointer }
	obj.vptr = unsafe.Pointer(&vt[2])
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	DynamicCast((*byte)(unsafe.Pointer(&obj)), (*byte)(unsafe.Pointer(&srcTI[0])), (*byte)(unsafe.Pointer(&dstTI[0])), -1)
}

func TestCxaAtexit(t *testing.T) {
	cxaHit = 0
	tmp := cxaTestDtor
	fn := *(**byte)(unsafe.Pointer(&tmp))
	var obj byte
	if CxaAtexit(fn, &obj, nil) != 0 {
		t.Fatal("ret")
	}
	runCxaAtExit()
	if cxaHit != 1 || obj != 7 {
		t.Fatalf("hit=%d obj=%d", cxaHit, obj)
	}
}
