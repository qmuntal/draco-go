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

type GeometryType int

const (
	GT_INVALID GeometryType = iota - 1
	GT_POSITION
	GT_NORMAL
	GT_COLOR
	GT_TEX_COORD
	GT_GENERIC
)

type DataType int

const (
	DT_INVALID DataType = iota
	DT_INT8
	DT_UINT8
	DT_INT16
	DT_UINT16
	DT_INT32
	DT_UINT32
	DT_INT64
	DT_UINT64
	DT_FLOAT32
	DT_FLOAT64
	DT_BOOL
)

type Face = [3]uint32

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

type PointAttr struct {
	ref *C.draco_point_attr
}

func (pa *PointAttr) Type() GeometryType {
	return GeometryType(C.dracoPointAttrType(pa.ref))
}

func (pa *PointAttr) DataType() DataType {
	return DataType(C.dracoPointAttrDataType(pa.ref))
}

func (pa *PointAttr) NumComponents() int8 {
	return int8(C.dracoPointAttrNumComponents(pa.ref))
}

func (pa *PointAttr) Normalized() bool {
	return bool(C.dracoPointAttrNormalized(pa.ref))
}

func (pa *PointAttr) ByteStride() int64 {
	return int64(C.dracoPointAttrByteStride(pa.ref))
}

func (pa *PointAttr) ByteOffset() int64 {
	return int64(C.dracoPointAttrByteOffset(pa.ref))
}

func (pa *PointAttr) UniqueID() uint32 {
	return uint32(C.dracoPointAttrUniqueId(pa.ref))
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

func (m *Mesh) NumAttrs() int32 {
	return int32(C.dracoMeshNumAttrs(m.ref))
}

func (m *Mesh) Faces(buffer []Face) []Face {
	n := m.NumFaces()
	if len(buffer) < int(n) {
		buffer = append(buffer, make([]Face, int(n)-len(buffer))...)
	}
	C.dracoMeshGetTrianglesUint32(m.ref, C.size_t(n*3*4), (*C.uint32_t)(unsafe.Pointer(&buffer[0])))
	return buffer[:n]
}

func (m *Mesh) Attr(i int32) *PointAttr {
	attr := C.dracoMeshGetAttribute(m.ref, C.int32_t(i))
	if attr == nil {
		return nil
	}
	return &PointAttr{ref: attr}
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
