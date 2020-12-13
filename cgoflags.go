package draco

/*
#cgo LDFLAGS: -L${SRCDIR}/lib
#cgo windows,amd64 LDFLAGS: -ldraco_windows_amd64 -lstdc++
#cgo linux,amd64 LDFLAGS: -ldraco_linux_amd64 -lstdc++ -lm
*/
import "C"
