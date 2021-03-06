package magick

// #include <string.h>
// #include <magick/api.h>
import "C"

import (
	"fmt"
	"runtime"
	"strings"
	"unsafe"
)

// Info is used to specify the encoding parameters like
// format and quality when encoding and image.
type Info struct {
	info *C.ImageInfo
}

// Format returns the format used for encoding the image.
func (in *Info) Format() string {
	return C.GoString(&in.info.magick[0])
}

// SetFormat sets the image format for encoding this image.
// See http://www.graphicsmagick.org for a list of supported
// formats.
func (in *Info) SetFormat(format string) {
	if format == "" {
		in.info.magick[0] = 0
	} else {
		s := C.CString(format)
		defer C.free(unsafe.Pointer(s))
		C.strncpy(&in.info.magick[0], s, C.MaxTextExtent)
	}
}

// Quality returns the quality used when compressing the image.
// This parameter does not affect all formats.
func (in *Info) Quality() uint {
	return uint(in.info.quality)
}

// SetQuality sets the quality used when compressing the image.
// This parameter does not affect all formats.
func (in *Info) SetQuality(q uint) {
	in.info.quality = magickSize(q)
}

func (in *Info) SetCompression(compression string) error {
	switch strings.ToLower(compression) {
	case "undefined":
		in.info.compression = C.UndefinedCompression
	case "none", "":
		in.info.compression = C.NoCompression
	case "bzip":
		in.info.compression = C.BZipCompression
	case "fax":
		in.info.compression = C.FaxCompression
	case "group4":
		in.info.compression = C.Group4Compression
	case "jpeg":
		in.info.compression = C.JPEGCompression
	case "lzw":
		in.info.compression = C.LZWCompression
	case "runlengthencoded":
		in.info.compression = C.RunlengthEncodedCompression
	case "zip":
		in.info.compression = C.ZipCompression
	default:
		return fmt.Errorf("invalid compression %#v", compression)
	}
	return nil
}

func (in *Info) SetDensity(x_density uint, y_density uint) {
	in.info.density = C.CString(fmt.Sprintf("%dx%d", x_density, y_density))
}

// Colorspace returns the colorspace used when encoding the image.
func (in *Info) Colorspace() Colorspace {
	return Colorspace(in.info.colorspace)
}

// SetColorspace set the colorspace used when encoding the image.
// Note that not all colorspaces are supported for encoding. See
// the documentation on Colorspace.
func (in *Info) SetColorspace(cs Colorspace) {
	in.info.colorspace = C.ColorspaceType(cs)
}

func NewBaseInfo() *Info {
	cinfo := C.CloneImageInfo(nil)
	info := new(Info)
	info.info = cinfo
	return info
}

// NewInfo returns a newly allocated *Info structure. Do not
// create Info objects directly, since they need to allocate
// some internal structures while being created.
func NewInfo() *Info {
	info := NewBaseInfo()
	runtime.SetFinalizer(info, func(i *Info) {
		if i.info != nil {
			C.DestroyImageInfo(i.info)
			i.info = nil
		}
	})
	return info
}
