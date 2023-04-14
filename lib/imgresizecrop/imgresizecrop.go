package imgresizecrop

import (
	"fmt"
)

type Path struct {
	SourceImage             string
	RenameTo                string
	PathSourceImage         string
	DestinationResized      string
	DestinationCroppedImage string
	DeleteOriginalFile      bool
}

type BoundsToCrop struct {
	X           int
	Y           int
	X1          int
	Y1          int
	imageWidth  int
	imageHeight int
	MinWidth    int
	MinHeight   int
}

type ResizeTo struct {
	Width  int
	Height int
}

func Resize(path *Path, resizeTo *ResizeTo) error {

	//	var imageSrc, _, err = readImageFromFilePath(path.PathSourceImage, path.SourceImage)
	var imageSrc, err = Open(fmt.Sprintf("%s/%s", path.PathSourceImage, path.SourceImage))
	if err != nil {
		return err
	}

	saveImage(resizeNearest(imageSrc, resizeTo.Width, resizeTo.Height), resizeTo.Width, resizeTo.Height, path.DestinationResized, path.RenameTo)
	imageSrc = nil

	defer handleDeleteOriginal(path.PathSourceImage, path.SourceImage, path.DeleteOriginalFile)

	return nil
}

func Crop(path *Path, boundsToCrop *BoundsToCrop, rotateAngle float64) error {

	var _, err = handleCrop(path, boundsToCrop, rotateAngle)
	if err != nil {
		return err
	}

	defer handleDeleteOriginal(path.PathSourceImage, path.SourceImage, path.DeleteOriginalFile)

	return nil
}

func CropRotateResize(path *Path, boundsToCrop *BoundsToCrop, rotateAngle float64, resizeTo *ResizeTo) error {

	var croppedImage, err = handleCrop(path, boundsToCrop, rotateAngle)
	if err != nil {
		return err
	}

	saveImage(resizeNearest(croppedImage, resizeTo.Width, resizeTo.Height), resizeTo.Width, resizeTo.Height, path.DestinationResized, path.RenameTo)
	//	fmt.Println(path, resizeTo)
	//	Resize(path, resizeTo)

	croppedImage = nil

	defer handleDeleteOriginal(path.PathSourceImage, path.SourceImage, path.DeleteOriginalFile)

	return nil
}

func CalculateThumbnailDimensions(path *Path, thumbnailSize int) *ResizeTo {
	return _calculateThumbnailDimensions(path, thumbnailSize)
}
