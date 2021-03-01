package draco

// #include "packaged/include/c_api.h"
import "C"
import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

type PointCloud struct {
	ref *C.draco_point_cloud
}

func (pc *PointCloud) free() {
	if pc.ref != nil {
		C.dracoPointCloudRelease(pc.ref)
	}
}

func NewPointCloud() *PointCloud {
	pc := &PointCloud{C.dracoNewPointCloud()}
	runtime.SetFinalizer(pc, (*PointCloud).free)
	return pc
}

func (pc *PointCloud) NumPoints() uint32 {
	return uint32(C.dracoPointCloudNumPoints(pc.ref))
}

func (pc *PointCloud) NumAttrs() int32 {
	return int32(C.dracoPointCloudNumAttrs(pc.ref))
}

func (pc *PointCloud) Attr(i int32) *PointAttr {
	attr := C.dracoPointCloudGetAttribute(pc.ref, C.int32_t(i))
	if attr == nil {
		return nil
	}
	return &PointAttr{ref: attr}
}

func (pc *PointCloud) AttrByUniqueID(id uint32) *PointAttr {
	attr := C.dracoPointCloudGetAttributeByUniqueId(pc.ref, C.uint32_t(id))
	if attr == nil {
		return nil
	}
	return &PointAttr{ref: attr}
}

func (pc *PointCloud) NamedAttributeID(gt GeometryAttrType) int32 {
	return int32(C.dracoPointCloudGetNamedAttributeId(pc.ref, C.draco_geometry_attr_type(gt)))
}

func (pc *PointCloud) AttrData(pa *PointAttr, buffer interface{}) (interface{}, bool) {
	var dt DataType
	n := pc.NumPoints() * uint32(pa.NumComponents())
	if buffer == nil {
		dt = pa.DataType()
		buffer = reflect.MakeSlice(reflect.SliceOf(dt.goType()), int(n), int(n)).Interface()
	} else {
		v := reflect.ValueOf(buffer)
		if v.IsNil() {
			buffer = reflect.MakeSlice(reflect.SliceOf(dt.goType()), int(n), int(n)).Interface()
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
			tmp := reflect.MakeSlice(reflect.SliceOf(dt.goType()), int(n)-l, int(n)-l).Interface()
			buffer = reflect.AppendSlice(reflect.ValueOf(buffer), reflect.ValueOf(tmp)).Interface()
		}
	}
	v := reflect.ValueOf(buffer).Index(0)
	size := n * dt.Size()
	ok := C.dracoPointCloudGetAttributeData(pc.ref, pa.ref, C.draco_data_type(dt), C.size_t(size), unsafe.Pointer(v.UnsafeAddr()))
	return buffer, bool(ok)
}
