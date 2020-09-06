package extendedimage

import "image"

type SubImager interface {
	SubImage(image.Rectangle) image.Image
}
