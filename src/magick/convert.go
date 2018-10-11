package magick

/*
#cgo pkg-config: MagickCore
#include <magick/MagickCore.h>
#include "magick.h"
*/
import "C"
import (
	"fmt"
	"image"
	"reflect"
	"unsafe"

	"github.com/uoregon-libraries/gopkg/logger"
)

var Logger = logger.DefaultLogger

// Image returns a native Go image interface.  For now, this is always RGBA for
// simplicity, but it would be a good idea to use a gray image when it makes
// sense to improve performance and RAM usage.
func (i *Image) Image() (image.Image, error) {
	// Create and prep-for-freeing the exception
	exception := C.AcquireExceptionInfo()
	defer C.DestroyExceptionInfo(exception)

	img := image.NewRGBA(image.Rect(0, 0, i.decodeWidth, i.decodeHeight))

	area := i.decodeWidth * i.decodeHeight
	pixLen := area << 2
	pixels := make([]byte, pixLen)
	pi := reflect.ValueOf(pixels).Interface()
	ptr := unsafe.Pointer(&pixels[0])

	// Dimensions as C types
	w := C.size_t(i.decodeWidth)
	h := C.size_t(i.decodeHeight)

	i.attemptExportRGBA(w, h, ptr, exception, 0)

	img.Pix = pi.([]uint8)

	return img, nil
}

func (i *Image) attemptExportRGBA(w, h C.size_t, ptr unsafe.Pointer, ex *C.ExceptionInfo, tries int) (err error) {
	defer func() {
		if x := recover(); x != nil {
			if tries < 3 {
				Logger.Warnf("Error trying to decode from ImageMagick (trying again): %s", x)
				i.attemptExportRGBA(w, h, ptr, ex, tries+1)
			} else {
				Logger.Errorf("Error trying to decode from ImageMagick: %s", x)
				err = fmt.Errorf("imagemagick failure: %s", x)
			}
		}
	}()

	C.ExportRGBA(i.image, w, h, ptr, ex)
	return
}
