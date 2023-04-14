package imgresizecrop

import (
	"errors"
	"fmt"
	"image"
	//_ "image/jpeg"
	"bytes"
	"image/color"
	//	"image/png"
	"io"
	"math"
	"os"
	"runtime"
	"sync"
)

func handleDeleteOriginal(path, filename string, deleteOriginal bool) {
	if deleteOriginal {
		deleteFile(fmt.Sprintf("%s/%s", path, filename))
	}
}

func deleteFile(path string) {

	// make sure the file to delete exists
	// befor the remove func is exectued
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	} else {
		err := os.Remove(path)
		if err != nil {
			//			fmt.Println(err)
			panic(err)
		}
	}

}

func handleCrop(path *Path, boundsToCrop *BoundsToCrop, rotateAngle float64) (*image.NRGBA, error) {
	//var imageSrc, _, err = readImageFromFilePath(path.PathSourceImage, path.SourceImage)
	imageSrc, err := Open(fmt.Sprintf("%s/%s", path.PathSourceImage, path.SourceImage))
	if err != nil {
		return nil, err
	}

	if rotateAngle != 0 {
		imageSrc = Rotate(imageSrc, rotateAngle, color.Black)
	}

	fixCropperJSConstrainsBug(imageSrc.Bounds(), boundsToCrop)

	var errBounds = validateCropBounds(imageSrc.Bounds(), boundsToCrop)
	if errBounds != nil {
		return nil, errBounds
	}

	var croppedImage = _crop(imageSrc, image.Rect(boundsToCrop.X, boundsToCrop.Y, boundsToCrop.X1, boundsToCrop.Y1))
	//	fmt.Println(boundsToCrop.X, boundsToCrop.Y, boundsToCrop.X1, boundsToCrop.Y1)
	//	fmt.Println("Cropped Image Bounds:", croppedImage.Bounds())
	saveImage(croppedImage, (boundsToCrop.X1 - boundsToCrop.X), (boundsToCrop.Y1 - boundsToCrop.Y), path.DestinationCroppedImage, path.RenameTo)
	imageSrc = nil

	return croppedImage, nil
}

func saveImage(srcImage *image.NRGBA, width, height int, path, filename string) {

	// Create a new image and paste the four produced images into it.
	dst := New(width, height, color.NRGBA{0, 0, 0, 0})
	dst = Paste(dst, srcImage, image.Pt(0, 0))
	//	fmt.Println(dst)

	// Save the resulting image as JPEG.
	var err = Save(dst, fmt.Sprintf("%s/%s", path, filename))
	if err != nil {
		fmt.Printf("failed to save image: %v, path: %s, filename: %s\n", err, path, filename)
	}
	//	dst = nil

}

/*
func readImageFromFilePath(filePath, fileName string) (image.Image, string, error) {
	f, err := os.Open(fmt.Sprintf("%s/%s", filePath, fileName))
	if err != nil {
		return nil, "", err
	}
	image, imageType, err := image.Decode(f)
	defer f.Close()
	return image, imageType, err
}
*/
// Crop cuts out a rectangular region with the specified bounds
// from the image and returns the cropped image.
func _crop(img image.Image, rect image.Rectangle) *image.NRGBA {

	r := rect.Intersect(img.Bounds()).Sub(img.Bounds().Min)
	if r.Empty() {
		return &image.NRGBA{}
	}
	src := newScanner(img)
	dst := image.NewNRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
	rowSize := r.Dx() * 4
	parallel(r.Min.Y, r.Max.Y, func(ys <-chan int) {
		for y := range ys {
			i := (y - r.Min.Y) * dst.Stride
			src.scan(r.Min.X, y, r.Max.X, y+1, dst.Pix[i:i+rowSize])
		}
	})
	return dst
}

func _calculateThumbnailDimensions(path *Path, thumbnailSize int) *ResizeTo {
	imageSrc, err := os.Open(fmt.Sprintf("%s/%s", path.PathSourceImage, path.SourceImage))
	if err != nil {
		fmt.Println(err)
	}
	defer imageSrc.Close()
	var imgReader io.Reader
	imgReader = imageSrc
	imgDecoder, _, err := image.DecodeConfig(imgReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding image: %v\n", err)
	}

	var aspectRatio = calculateAspectRatio(imgDecoder.Width, imgDecoder.Height, thumbnailSize)

	//	fmt.Printf("Dimensiion: %d %d\n Aspect Ratio: %d\n", imgDecoder.Width, imgDecoder.Height, aspectRatio)

	//	var height = int(float64(thumbnailSize) * aspectRatio)
	//	fmt.Println(height)

	var resizeTo = new(ResizeTo)
	if imgDecoder.Width == imgDecoder.Height {
		// it's a square image
		resizeTo.Width = thumbnailSize
		resizeTo.Height = thumbnailSize
	} else if imgDecoder.Width > imgDecoder.Height {
		resizeTo.Width = thumbnailSize
		resizeTo.Height = int(float64(thumbnailSize) * aspectRatio)
	} else {
		resizeTo.Width = int(float64(thumbnailSize) * aspectRatio)
		resizeTo.Height = thumbnailSize
	}

	//	fmt.Printf("Width: %d, Height: %d\n", resizeTo.Width, resizeTo.Height)

	return resizeTo

}

func calculateAspectRatio(srcWidth, srcHeight, thumbnailSize int) float64 {

	var wRatio = (float64(srcWidth) / float64(srcHeight))
	var hRatio = (float64(srcHeight) / float64(srcWidth))
	//	fmt.Println("W Ratio:", wRatio)
	//	fmt.Println("H Ratio:", hRatio)

	// return the minimum most of ratios
	if wRatio > hRatio {
		//		fmt.Println("ratio:", hRatio)
		return hRatio
	}
	//	fmt.Println("ratio:", wRatio)
	return wRatio
}

// resizeNearest is a fast nearest-neighbor resize, no filtering.
func resizeNearest(img image.Image, width, height int) *image.NRGBA {

	if height == 0 {
		tmpH := float64(width) * float64(img.Bounds().Dy()) / float64(img.Bounds().Dx())
		height = int(math.Max(1.0, math.Floor(tmpH+0.5)))
	}

	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	dx := float64(img.Bounds().Dx()) / float64(width)
	dy := float64(img.Bounds().Dy()) / float64(height)

	if dx > 1 && dy > 1 {
		src := newScanner(img)
		parallel(0, height, func(ys <-chan int) {
			for y := range ys {
				srcY := int((float64(y) + 0.5) * dy)
				dstOff := y * dst.Stride
				for x := 0; x < width; x++ {
					srcX := int((float64(x) + 0.5) * dx)
					src.scan(srcX, srcY, srcX+1, srcY+1, dst.Pix[dstOff:dstOff+4])
					dstOff += 4
				}
			}
		})
	} else {
		src := toNRGBA(img)
		parallel(0, height, func(ys <-chan int) {
			for y := range ys {
				srcY := int((float64(y) + 0.5) * dy)
				srcOff0 := srcY * src.Stride
				dstOff := y * dst.Stride
				for x := 0; x < width; x++ {
					srcX := int((float64(x) + 0.5) * dx)
					srcOff := srcOff0 + srcX*4
					copy(dst.Pix[dstOff:dstOff+4], src.Pix[srcOff:srcOff+4])
					dstOff += 4
				}
			}
		})
	}

	return dst
}

func fixCropperJSConstrainsBug(imageBounds image.Rectangle, boundsToCrop *BoundsToCrop) {

	// readjust bounds according to given crop
	boundsToCrop.X1 = boundsToCrop.X + boundsToCrop.X1
	boundsToCrop.Y1 = boundsToCrop.Y + boundsToCrop.Y1

	if boundsToCrop.X < 0 {
		boundsToCrop.X = 0
	}
	if boundsToCrop.Y < 0 {
		boundsToCrop.Y = 0
	}
	if boundsToCrop.X1 > imageBounds.Max.X {
		boundsToCrop.X1 = imageBounds.Max.X
	}
	if boundsToCrop.Y1 > imageBounds.Max.Y {
		boundsToCrop.Y1 = imageBounds.Max.Y
	}

	// set the final width and height before
	// overriding the x1 and y1 values
	boundsToCrop.imageWidth = boundsToCrop.X1
	boundsToCrop.imageHeight = boundsToCrop.Y1

}

func validateCropBounds(imageBounds image.Rectangle, boundsToCrop *BoundsToCrop) error {

	if (boundsToCrop.X < imageBounds.Min.X || boundsToCrop.X >= imageBounds.Max.X) ||
		(boundsToCrop.Y < imageBounds.Min.Y || boundsToCrop.Y >= imageBounds.Max.Y) ||
		(boundsToCrop.X1 > imageBounds.Max.X || boundsToCrop.X1 <= boundsToCrop.X) ||
		(boundsToCrop.Y1 > imageBounds.Max.Y || boundsToCrop.Y1 <= boundsToCrop.Y) ||
		((boundsToCrop.X + boundsToCrop.X1) < boundsToCrop.MinWidth) ||
		((boundsToCrop.Y + boundsToCrop.Y1) < boundsToCrop.MinHeight) {
		return errors.New("Out of bounds crop not possible")
	}
	return nil
}

// parallel processes the data in separate goroutines.
func parallel(start, stop int, fn func(<-chan int)) {
	count := stop - start
	if count < 1 {
		return
	}

	procs := runtime.GOMAXPROCS(0)
	if procs > count {
		procs = count
	}

	c := make(chan int, count)
	for i := start; i < stop; i++ {
		c <- i
	}
	close(c)

	var wg sync.WaitGroup
	for i := 0; i < procs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(c)
		}()
	}
	wg.Wait()
}

func toNRGBA(img image.Image) *image.NRGBA {
	if img, ok := img.(*image.NRGBA); ok {
		return &image.NRGBA{
			Pix:    img.Pix,
			Stride: img.Stride,
			Rect:   img.Rect.Sub(img.Rect.Min),
		}
	}
	return Clone(img)
}

// New creates a new image with the specified width and height, and fills it with the specified color.
func New(width, height int, fillColor color.Color) *image.NRGBA {
	if width <= 0 || height <= 0 {
		return &image.NRGBA{}
	}

	c := color.NRGBAModel.Convert(fillColor).(color.NRGBA)
	if (c == color.NRGBA{0, 0, 0, 0}) {
		return image.NewNRGBA(image.Rect(0, 0, width, height))
	}

	return &image.NRGBA{
		Pix:    bytes.Repeat([]byte{c.R, c.G, c.B, c.A}, width*height),
		Stride: 4 * width,
		Rect:   image.Rect(0, 0, width, height),
	}
}

// Paste pastes the img image to the background image at the specified position and returns the combined image.
func Paste(background, img image.Image, pos image.Point) *image.NRGBA {
	dst := Clone(background)
	pos = pos.Sub(background.Bounds().Min)
	pasteRect := image.Rectangle{Min: pos, Max: pos.Add(img.Bounds().Size())}
	interRect := pasteRect.Intersect(dst.Bounds())
	if interRect.Empty() {
		return dst
	}
	src := newScanner(img)
	parallel(interRect.Min.Y, interRect.Max.Y, func(ys <-chan int) {
		for y := range ys {
			x1 := interRect.Min.X - pasteRect.Min.X
			x2 := interRect.Max.X - pasteRect.Min.X
			y1 := y - pasteRect.Min.Y
			y2 := y1 + 1
			i1 := y*dst.Stride + interRect.Min.X*4
			i2 := i1 + interRect.Dx()*4
			src.scan(x1, y1, x2, y2, dst.Pix[i1:i2])
		}
	})
	return dst
}

// Clone returns a copy of the given image.
func Clone(img image.Image) *image.NRGBA {
	src := newScanner(img)
	dst := image.NewNRGBA(image.Rect(0, 0, src.w, src.h))
	size := src.w * 4
	parallel(0, src.h, func(ys <-chan int) {
		for y := range ys {
			i := y * dst.Stride
			src.scan(0, y, src.w, y+1, dst.Pix[i:i+size])
		}
	})
	return dst
}

// clamp rounds and clamps float64 value to fit into uint8.
func clamp(x float64) uint8 {
	v := int64(x + 0.5)
	if v > 255 {
		return 255
	}
	if v > 0 {
		return uint8(v)
	}
	return 0
}

func reverse(pix []uint8) {
	if len(pix) <= 4 {
		return
	}
	i := 0
	j := len(pix) - 4
	for i < j {
		pi := pix[i : i+4 : i+4]
		pj := pix[j : j+4 : j+4]
		pi[0], pj[0] = pj[0], pi[0]
		pi[1], pj[1] = pj[1], pi[1]
		pi[2], pj[2] = pj[2], pi[2]
		pi[3], pj[3] = pj[3], pi[3]
		i += 4
		j -= 4
	}
}

/*
// absint returns the absolute value of i.
func absint(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// rgbToHSL converts a color from RGB to HSL.
func rgbToHSL(r, g, b uint8) (float64, float64, float64) {
	rr := float64(r) / 255
	gg := float64(g) / 255
	bb := float64(b) / 255

	max := math.Max(rr, math.Max(gg, bb))
	min := math.Min(rr, math.Min(gg, bb))

	l := (max + min) / 2

	if max == min {
		return 0, 0, l
	}

	var h, s float64
	d := max - min
	if l > 0.5 {
		s = d / (2 - max - min)
	} else {
		s = d / (max + min)
	}

	switch max {
	case rr:
		h = (gg - bb) / d
		if g < b {
			h += 6
		}
	case gg:
		h = (bb-rr)/d + 2
	case bb:
		h = (rr-gg)/d + 4
	}
	h /= 6

	return h, s, l
}

// hslToRGB converts a color from HSL to RGB.
func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	var r, g, b float64
	if s == 0 {
		v := clamp(l * 255)
		return v, v, v
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	r = hueToRGB(p, q, h+1/3.0)
	g = hueToRGB(p, q, h)
	b = hueToRGB(p, q, h-1/3.0)

	return clamp(r * 255), clamp(g * 255), clamp(b * 255)
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	if t < 1/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1/2.0 {
		return q
	}
	if t < 2/3.0 {
		return p + (q-p)*(2/3.0-t)*6
	}
	return p
}
*/
