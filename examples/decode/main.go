package main

import (
	"io/ioutil"
	"log"

	"github.com/qmuntal/draco-go/draco"
)

func main() {
	data, err := ioutil.ReadFile("../../testdata/test_nm.obj.edgebreaker.cl4.2.2.drc")
	if err != nil {
		log.Fatalf("failed to read test file: %v", err)
	}
	m := draco.NewMesh()
	d := draco.NewDecoder()
	if err := d.DecodeMesh(m, data); err != nil {
		log.Fatalf("failed to decode mesh: %v", err)
	}
	log.Println("point count:", m.NumPoints())
	log.Println("face count:", m.NumFaces())
	log.Println("faces:", m.Faces(nil))
}
