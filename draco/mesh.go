package draco

// #include "packaged/include/c_api.h"
import "C"
import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

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

func (m *Mesh) AttrData(pa *PointAttr, buffer interface{}) (interface{}, bool) {
	var dt DataType
	n := m.NumPoints() * uint32(pa.NumComponents())
	if buffer == nil {
		dt = pa.DataType()
		buffer = reflect.MakeSlice(reflect.SliceOf(dt.Type()), int(n), int(n)).Interface()
	} else {
		v := reflect.ValueOf(buffer)
		if v.IsNil() {
			buffer = reflect.MakeSlice(reflect.SliceOf(dt.Type()), int(n), int(n)).Interface()
		}
		if v.Kind() != reflect.Slice {
			panic(fmt.Sprintf("draco-go: expecting a slice but got %s", v.Kind()))
		}
		l := v.Len()
		switch buffer.(type) {
		case []int8:
			dt = DT_INT8
		case []uint8:
			dt = DT_UINT8
		case []int16:
			dt = DT_INT16
		case []uint16:
			dt = DT_UINT16
		case []int32:
			dt = DT_INT32
		case []uint32:
			dt = DT_UINT32
		case []int64:
			dt = DT_INT64
		case []uint64:
			dt = DT_UINT64
		case []float32:
			dt = DT_FLOAT32
		case []float64:
			dt = DT_FLOAT64
		default:
			panic("draco-go: unsupported data type")
		}
		if l < int(n) {
			tmp := reflect.MakeSlice(reflect.SliceOf(dt.Type()), int(n)-l, int(n)-l).Interface()
			buffer = reflect.AppendSlice(reflect.ValueOf(buffer), reflect.ValueOf(tmp)).Interface()
		}
	}
	v := reflect.ValueOf(buffer).Index(0)
	size := n * dt.Size()
	ok := C.dracoMeshGetAttributeData(m.ref, pa.ref, C.dracoDataType(dt), C.size_t(size), unsafe.Pointer(v.UnsafeAddr()))
	return buffer, bool(ok)
}
