// +build !customenv
package draco

/*
#cgo CFLAGS: -DDRACO_C_STATIC
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/packaged/lib/windows_amd64
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/packaged/lib/linux_amd64
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/packaged/lib/darwin_amd64
#cgo LDFLAGS: -ldraco_c -lstdc++
#cgo !windows LDFLAGS: -lm
*/
import "C"
