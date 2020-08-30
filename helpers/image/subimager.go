package image

import "image"

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}
