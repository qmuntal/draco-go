package draco

import (
	"io/ioutil"
	"testing"
)

func TestDecode(t *testing.T) {
	data, err := ioutil.ReadFile("./testdata/test_nm.obj.edgebreaker.cl4.2.2.drc")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}
	m := NewMesh()
	d := NewDecoder()
	if err := d.DecodeMesh(data, m); err != nil {
		t.Fatalf("failed to decode mesh: %v", err)
	}
	if n := m.NumFaces(); n != 170 {
		t.Errorf("Mesh.NumFaces got 170, got %d", n)
	}
	if n := m.NumPoints(); n != 99 {
		t.Errorf("Mesh.NumFaces got 99, got %d", n)
	}
	faces := m.Faces(nil)
	want := [3]uint32{0, 1, 2}
	if got := faces[0]; got != want {
		t.Errorf("Mesh.Faces[0] got %v, got %v", want, got)
	}
}
