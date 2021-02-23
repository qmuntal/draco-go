package draco

// #include "packaged/include/c_api.h"
import "C"
import (
	"runtime"
	"unsafe"
)

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
