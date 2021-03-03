package draco

import (
	"reflect"
	"testing"

	"github.com/qmuntal/gltf"
)

func TestUnmarshal(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"base", args{[]byte(`{
			"bufferView" : 5,
			"attributes" : {
				"POSITION" : 0,
				"NORMAL" : 1,
				"TEXCOORD_0" : 2,
				"WEIGHTS_0" : 3,
				"JOINTS_0" : 4
			}
		}`)}, &PrimitiveExt{BufferView: 5, Attributes: gltf.Attribute{
			"JOINTS_0":   4,
			"NORMAL":     1,
			"POSITION":   0,
			"TEXCOORD_0": 2,
			"WEIGHTS_0":  3,
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalMesh(t *testing.T) {
	doc, err := gltf.Open("testdata/box/Box.gltf")
	if err != nil {
		t.Fatal(err)
	}
	pd, err := UnmarshalMesh(doc, doc.BufferViews[0])
	if err != nil {
		t.Fatal(err)
	}
	p := doc.Meshes[0].Primitives[0]
	indWant := []uint32{2, 5, 6, 3, 11, 8, 8, 11, 12, 14, 9, 17}
	if got := pd.ReadIndices(nil); !reflect.DeepEqual(indWant, got) {
		t.Errorf("ReadIndices want %v, got %v", indWant, got)
	}
	_, err = pd.ReadAttr(p, "POSITION", nil)
	if err != nil {
		t.Error(err)
	}
	_, err = pd.ReadAttr(p, "NORMAL", [][3]float32{{1, 2, 3}})
	if err != nil {
		t.Error(err)
	}
	_, err = pd.ReadAttr(p, "OTHER", nil)
	if err != nil {
		t.Error(err)
	}
}

func TestGetPrimitiveExt(t *testing.T) {
	type args struct {
		p *gltf.Primitive
	}
	tests := []struct {
		name string
		args args
		want *PrimitiveExt
	}{
		{"no extension", args{&gltf.Primitive{}}, nil},
		{"other extension", args{&gltf.Primitive{Extensions: gltf.Extensions{"other": nil}}}, nil},
		{"draco other extension", args{&gltf.Primitive{Extensions: gltf.Extensions{ExtensionName: nil}}}, nil},
		{"draco extension", args{&gltf.Primitive{Extensions: gltf.Extensions{ExtensionName: &PrimitiveExt{
			BufferView: 1,
		}}}}, &PrimitiveExt{BufferView: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPrimitiveExt(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPrimitiveExt() = %v, want %v", got, tt.want)
			}
		})
	}
}
