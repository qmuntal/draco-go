package main

import (
	"io/ioutil"
	"log"

	"github.com/qmuntal/draco-go"
)

func main() {
	data, err := ioutil.ReadFile("../../testdata/test_nm.obj.edgebreaker.cl4.2.2.drc")
	if err != nil {
		log.Fatalf("failed to read test file: %v", err)
	}
	m := draco.NewMesh()
	d := draco.NewDecoder()
	if err := d.DecodeMesh(data, m); err != nil {
		log.Fatalf("failed to decode mesh: %v", err)
	}
	log.Println(m.NumFaces())
	if n := m.NumFaces(); n != 170 {
		log.Fatalf("Mesh.NumFaces got 170, got %d", n)
	}
	if n := m.NumPoints(); n != 99 {
		log.Fatalf("Mesh.NumFaces got 99, got %d", n)
	}
	faces := m.Faces(nil)
	want := [3]uint32{0, 1, 2}
	if got := faces[0]; got != want {
		log.Fatalf("Mesh.Faces[0] got %v, got %v", want, got)
	}
}
