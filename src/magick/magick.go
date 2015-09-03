// The magick package is a hacked up port of the minimal functionality we need
// to satisfy the IIIFImage interface.  Code is based in part on
// github.com/quirkey/magick
package magick

/*
#cgo pkg-config: MagickCore
#include <magick/MagickCore.h>
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

func init() {
	path, _ := os.Getwd()
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.MagickCoreGenesis(cPath, C.MagickFalse)
}

func makeError(exception *C.ExceptionInfo) error {
	return fmt.Errorf("%s: %s - %s", exception.severity, exception.reason, exception.description)
}
