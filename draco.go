package draco

/*
#include "c_api.h"
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("draco: [%d] %s", e.Code, e.Message)
}

func newError(s *C.draco_status) error {
	if C.dracoStatusOk(s) {
		return nil
	}
	err := &Error{
		Code:    int(C.dracoStatusCode(s)),
		Message: C.GoString(C.dracoStatusErrorMsg(s)),
	}
	C.dracoStatusRelease(s)
	s = nil
	return err
}

type Mesh struct {
	ref *C.draco_mesh
}

func (m *Mesh) free() {
	if m.ref != nil {
		C.dracoMeshRelease(m.ref)
	}
}

func NewMesh() *Mesh {
	m := &Mesh{C.dracoNewMesh()}
	runtime.SetFinalizer(m, (*Mesh).free)
	return m
}

func (m *Mesh) NumFaces() uint32 {
	return uint32(C.dracoMeshNumFaces(m.ref))
}

func (m *Mesh) NumPoints() uint32 {
	return uint32(C.dracoMeshNumPoints(m.ref))
}

func (m *Mesh) Faces(buffer [][3]uint32) [][3]uint32 {
	n := m.NumFaces()
	if len(buffer) < int(n) {
		buffer = append(buffer, make([][3]uint32, int(n)-len(buffer))...)
	}
	C.dracoMeshGetTrianglesUint32(m.ref, C.size_t(n*3*4), (*C.uint32_t)(unsafe.Pointer(&buffer[0])))
	return buffer[:n]
}

type Decoder struct {
	ref *C.draco_decoder
}

func (d *Decoder) free() {
	if d.ref != nil {
		C.dracoDecoderRelease(d.ref)
	}
}

func NewDecoder() *Decoder {
	d := &Decoder{C.dracoNewDecoder()}
	runtime.SetFinalizer(d, (*Decoder).free)
	return d
}

func (d *Decoder) DecodeMesh(data []byte, m *Mesh) error {
	s := C.dracoDecoderArrayToMesh(d.ref, (*C.char)(unsafe.Pointer(&data[0])), C.size_t(len(data)), m.ref)
	return newError(s)
}
