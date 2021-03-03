package draco

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/qmuntal/draco-go/draco"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/binary"
	"github.com/qmuntal/gltf/modeler"
)

const (
	// ExtensionName defines the KHR_draco_mesh_compression unique key.
	ExtensionName = "KHR_draco_mesh_compression"
)

func init() {
	gltf.RegisterExtension(ExtensionName, Unmarshal)
}

// Unmarshal decodes the json data into the correct type.
func Unmarshal(data []byte) (interface{}, error) {
	drc := new(PrimitiveExt)
	err := json.Unmarshal(data, drc)
	return drc, err
}

// PrimitiveExt extends the gltf.Primtive info to handle draco compressed meshes.
type PrimitiveExt struct {
	Extensions gltf.Extensions `json:"extensions,omitempty"`
	Extras     interface{}     `json:"extras,omitempty"`
	BufferView uint32          `json:"bufferView"`
	Attributes gltf.Attribute  `json:"attributes"`
}

// GetPrimitiveExt retrieve a PrimitiveExt from p.
// If p does not contain the draco extensions it returns nil.
func GetPrimitiveExt(p *gltf.Primitive) *PrimitiveExt {
	var pe *PrimitiveExt
	if ext, ok := p.Extensions[ExtensionName]; ok {
		if pe, ok = ext.(*PrimitiveExt); !ok {
			return nil
		}
	} else {
		return nil
	}
	return pe
}

// Mesh contains the necessary information to process a draco-encoded
// in a gltf context.
type Mesh struct {
	doc *gltf.Document
	m   *draco.Mesh
}

// UnmarshalMesh unmarshal the draco-encoded mesh from a gltf.BufferView
func UnmarshalMesh(doc *gltf.Document, bv *gltf.BufferView) (*Mesh, error) {
	data, err := modeler.ReadBufferView(doc, bv)
	if err != nil {
		return nil, err
	}
	if tp := draco.GetEncodedGeometryType(data); tp != draco.EGT_TRIANGULAR_MESH {
		return nil, fmt.Errorf("draco-go: unsupported geometry type %v", tp)
	}
	m := draco.NewMesh()
	d := draco.NewDecoder()
	if err := d.DecodeMesh(m, data); err != nil {
		return nil, err
	}
	return &Mesh{
		doc: doc,
		m:   m,
	}, nil
}

// ReadIndices reads the faces of the Mesh.
// buffer can be nil.
func (m Mesh) ReadIndices(buffer []uint32) []uint32 {
	return m.m.Faces(buffer)
}

// ReadAttr reads the named attribute of a gltf.Primitive.
// If the attribute is defined in the primitive but not in the mesh
// it fallbacks to modeler.ReadAccessor.
// buffer can be nil.
func (m Mesh) ReadAttr(p *gltf.Primitive, name string, buffer interface{}) (interface{}, error) {
	var (
		gltfIndex, dracoID uint32
		ok                 bool
	)
	if gltfIndex, ok = p.Attributes[name]; !ok {
		return nil, nil
	}
	pe := GetPrimitiveExt(p)
	ok = false
	if pe != nil {
		dracoID, ok = pe.Attributes[name]
	}
	acr := m.doc.Accessors[gltfIndex]
	if !ok {
		return modeler.ReadAccessor(m.doc, acr, buffer)
	}
	attr := m.m.AttrByUniqueID(dracoID)
	if attr == nil {
		return nil, fmt.Errorf("draco: mesh does not contain attribute %v", dracoID)
	}
	buffer = binary.MakeSliceBuffer(acr.ComponentType, acr.Type, acr.Count, buffer)
	sh := new(reflect.SliceHeader)
	sh.Data = reflect.ValueOf(buffer).Pointer()
	sh.Len = int(acr.Type.Components() * acr.Count)
	sh.Cap = sh.Len
	var data interface{}
	switch acr.ComponentType {
	case gltf.ComponentByte:
		data = *(*[]int8)(unsafe.Pointer(sh))
	case gltf.ComponentUbyte:
		data = *(*[]uint8)(unsafe.Pointer(sh))
	case gltf.ComponentShort:
		data = *(*[]int16)(unsafe.Pointer(sh))
	case gltf.ComponentUshort:
		data = *(*[]uint16)(unsafe.Pointer(sh))
	case gltf.ComponentUint:
		data = *(*[]uint32)(unsafe.Pointer(sh))
	case gltf.ComponentFloat:
		data = *(*[]float32)(unsafe.Pointer(sh))
	default:
		panic("draco-go: unsupported data type")
	}

	_, _ = m.m.AttrData(attr, data)
	return buffer, nil
}
