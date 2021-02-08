// +build !customenv
package draco

/*
#cgo LDFLAGS: -L${SRCDIR}/lib
#cgo windows CFLAGS: -DDRACO_STATIC
#cgo windows,amd64 LDFLAGS: -lcdraco_windows_amd64 -lstdc++
#cgo linux,amd64 LDFLAGS: -lcdraco_linux_amd64 -lstdc++ -lm
#cgo darwin,amd64 LDFLAGS: -lcdraco_darwin_amd64
*/
import "C"
